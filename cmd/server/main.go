package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"task-manager/internal/api"
	"task-manager/internal/auth"
	"task-manager/internal/config"
	"task-manager/internal/middleware"
	"task-manager/internal/repo"
	"task-manager/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.MustLoad()

	// ตั้งโหมดจาก config
	if cfg.GinMode != "" {
		gin.SetMode(cfg.GinMode)
	}

	// DB
	db := repo.MustOpen(cfg.DBDSN)
	defer db.Close()

	// DI
	j := auth.NewJWT(cfg)
	pw := auth.NewPasswordHasher()

	userRepo := repo.NewUserRepo(db)
	taskRepo := repo.NewTaskRepo(db)

	authSvc := service.NewAuthService(userRepo, pw, j)
	userSvc := service.NewUserService(userRepo, pw)
	taskSvc := service.NewTaskService(taskRepo)

	// Google OAuth
	googleCfg := auth.NewGoogleOAuthFromEnv()

	// Router
	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	r.Use(
		gin.Recovery(),
		middleware.CORSMiddleware(),
		middleware.SecureHeaders(),
		func(c *gin.Context) {
			start := time.Now()
			c.Next()
			log.Printf("%s %s %d (%s)",
				c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start),
			)
		},
	)

	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })



	// Static file serving สำหรับ frontend
	r.Static("/assets", "./frontend/vanilla/assets")
	r.Static("/auth", "./frontend/vanilla/auth")
	r.Static("/profile", "./frontend/vanilla/profile")
	r.Static("/settings", "./frontend/vanilla/settings")
	r.Static("/utils", "./frontend/vanilla/utils")
	r.StaticFile("/styles.css", "./frontend/vanilla/styles.css")
	r.StaticFile("/create_account.html", "./frontend/vanilla/create_account.html")
	r.StaticFile("/account-type.html", "./frontend/vanilla/account-type.html")
	r.StaticFile("/team-invitation.html", "./frontend/vanilla/team-invitation.html")
	r.StaticFile("/home.html", "./frontend/vanilla/home.html")
	r.StaticFile("/dashboard-reporting.html", "./frontend/vanilla/dashboard-reporting.html")
	r.StaticFile("/index.html", "./frontend/vanilla/index.html")

	// ส่ง arg ให้ครบ (เพิ่ม frontendURL เข้าไปเป็นตัวสุดท้าย)
	api.RegisterAuthRoutes(r, authSvc, userRepo, googleCfg, cfg.FrontendURL)
	api.RegisterUserRoutes(r, userSvc, middleware.JWTMiddleware(&j))
	
	// Root route - ต้องอยู่ท้ายสุดเพื่อไม่ให้ override routes อื่น
	r.StaticFile("/", "./frontend/vanilla/index.html")
	api.RegisterTaskRoutes(r, taskSvc, j)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("listening on :%s (mode=%s)", cfg.Port, gin.Mode())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exited gracefully")
}
