package api

import (
	"net/http"
	"strconv"

	"task-manager/internal/auth"
	"task-manager/internal/middleware"
	"task-manager/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	Svc service.TaskService
}

func RegisterTaskRoutes(r *gin.Engine, svc service.TaskService, jwt auth.JWT) {
	h := &TaskHandler{Svc: svc}

	g := r.Group("/api/tasks")
	g.Use(middleware.JWTMiddleware(&jwt))
	{
		g.GET("", h.getTasks)
		g.POST("", h.createTask)
		g.PUT("/:id", h.updateTask)
		g.DELETE("/:id", h.deleteTask)
	}

	// Also register /tasks for backward compatibility
	r.GET("/tasks", middleware.JWTMiddleware(&jwt), h.getTasks)
}

func (h *TaskHandler) getTasks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit := 200 // default limit
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	tasks, err := h.Svc.GetUserTasks(c.Request.Context(), userID.(int), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (h *TaskHandler) createTask(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (h *TaskHandler) updateTask(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (h *TaskHandler) deleteTask(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
