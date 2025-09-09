package api

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"task-manager/internal/auth"
	"task-manager/internal/repo"
	"task-manager/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Svc         service.AuthService
	UserRepo    repo.UserRepo
	GoogleCfg   *auth.GoogleOAuth
	FrontendURL string
}

func RegisterAuthRoutes(r *gin.Engine, svc service.AuthService, userRepo repo.UserRepo, gcfg *auth.GoogleOAuth, frontendURL string) {
	h := &AuthHandler{Svc: svc, UserRepo: userRepo, GoogleCfg: gcfg, FrontendURL: frontendURL}

	api := r.Group("/api/auth")
	{
		api.POST("/check-email", h.checkEmail)
		api.POST("/login", h.login)
		api.POST("/register", h.register)
		api.POST("/complete-google-registration", h.completeGoogleRegistration)
		api.GET("/google/login", h.googleLogin)
		api.GET("/google/callback", h.googleCallback)
	}
}

func (h *AuthHandler) checkEmail(c *gin.Context) {
	var in struct {
		Email string `json:"email"`
	}
	if err := c.BindJSON(&in); err != nil || in.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}
	exists, err := h.UserRepo.EmailExists(c.Request.Context(), in.Email)
	if err != nil {
		// Log the actual error for debugging
		println("EmailExists error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"exists": exists})
}

func (h *AuthHandler) login(c *gin.Context) {
	var in struct {
		Email    string `json:"email"` // Actually usernameOrEmail
		Password string `json:"password"`
	}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	if in.Email == "" || in.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username/email and password are required"})
		return
	}

	// Authenticate user
	user, err := h.Svc.Login(c.Request.Context(), in.Email, in.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate token
	token, err := h.Svc.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	// Set cookie
	c.SetCookie("access_token", token, int((24 * time.Hour).Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "login successful",
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func (h *AuthHandler) register(c *gin.Context) {
	var in struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	if in.Email == "" || in.Username == "" || in.Password == "" || in.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email, username, password, and name are required"})
		return
	}

	// Check if email already exists
	exists, err := h.UserRepo.EmailExists(c.Request.Context(), in.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	// Check if username already exists
	usernameExists, err := h.UserRepo.UsernameExists(c.Request.Context(), in.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	if usernameExists {
		c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
		return
	}

	// Register user
	user, err := h.Svc.Register(c.Request.Context(), in.Email, in.Username, in.Password, in.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	// Generate token for the new user
	token, err := h.Svc.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	// Set cookie
	c.SetCookie("access_token", token, int((24 * time.Hour).Seconds()), "/", "", false, true)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "registration successful",
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func (h *AuthHandler) completeGoogleRegistration(c *gin.Context) {
	var in struct {
		Email    string `json:"email" binding:"required,email"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	// Check if username already exists
	usernameExists, err := h.UserRepo.UsernameExists(c.Request.Context(), in.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	if usernameExists {
		c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
		return
	}

	// Complete Google registration
	user, err := h.Svc.CompleteGoogleRegistration(c.Request.Context(), in.Email, in.Username, in.Password, in.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	// Generate token for the updated user
	token, err := h.Svc.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	// Set cookie
	c.SetCookie("access_token", token, int((24 * time.Hour).Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "registration completed successfully",
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name.String,
		},
	})
}

func (h *AuthHandler) googleLogin(c *gin.Context) {
	next := c.Query("next")
	onboardIfNew := c.Query("onboardIfNew")
	
	// Debug logging
	log.Printf("Google login - Next: %s, OnboardIfNew: %s", next, onboardIfNew)

	st := struct {
		Next         string `json:"next"`
		OnboardIfNew string `json:"onboardIfNew"`
	}{Next: next, OnboardIfNew: onboardIfNew}

	b, _ := json.Marshal(st)
	state := base64.URLEncoding.EncodeToString(b)
	
	// Debug state
	log.Printf("State encoded: %s", state)

	url := h.GoogleCfg.LoginURL(state)
	c.Redirect(http.StatusFound, url)
}

func (h *AuthHandler) googleCallback(c *gin.Context) {
	code := c.Query("code")
	stateB64 := c.Query("state")

	var st struct {
		Next         string `json:"next"`
		OnboardIfNew string `json:"onboardIfNew"`
	}
	if sb, err := base64.URLEncoding.DecodeString(stateB64); err == nil {
		_ = json.Unmarshal(sb, &st)
	}
	
	// Debug logging
	log.Printf("OAuth callback - Next: %s, OnboardIfNew: %s", st.Next, st.OnboardIfNew)

	gu, err := h.GoogleCfg.FetchUser(c.Request.Context(), code)
	if err != nil {
		c.String(http.StatusBadRequest, "oauth error: %v", err)
		return
	}

	u, token, created, err := h.Svc.LoginOrSignupGoogle( // ← รับ 4 ค่า
		c.Request.Context(), gu.Email, gu.Name, gu.Sub, gu.Picture,
	)
	if err != nil {
		c.String(http.StatusInternalServerError, "auth error: %v", err)
		return
	}

	c.SetCookie("access_token", token, int((24 * time.Hour).Seconds()), "/", "", false, true)

	// flow ไปหน้าต่อ - ตรวจสอบผู้ใช้ใหม่ก่อน
  if created {
    // ผู้ใช้ใหม่ไปหน้า create_account.html พร้อมกับอีเมล
    c.Redirect(http.StatusFound, h.FrontendURL+"/create_account.html?email="+gu.Email)
    return
  }

  // ตรวจสอบว่าผู้ใช้เก่ามี username และ password หรือยัง
  if !u.Username.Valid || u.Username.String == "" || !u.PasswordHash.Valid || u.PasswordHash.String == "" {
    // ผู้ใช้เก่าที่ยังไม่มี username หรือ password ให้ไปตั้งค่า
    c.Redirect(http.StatusFound, h.FrontendURL+"/create_account.html?email="+gu.Email)
    return
  }

	if st.Next != "" {
		c.Redirect(http.StatusFound, st.Next)
		return
	}
	// ผู้ใช้เก่าที่มีข้อมูลครบแล้ว ไปหน้า dashboard
	c.Redirect(http.StatusFound, h.FrontendURL+"/dashboard/index.html")
}
