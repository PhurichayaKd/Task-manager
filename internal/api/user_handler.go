package api

import (
	"net/http"
	"task-manager/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserSvc service.UserService
}

func RegisterUserRoutes(r *gin.Engine, userSvc service.UserService, authMw gin.HandlerFunc) {
	h := &UserHandler{UserSvc: userSvc}

	g := r.Group("/api/users")
	g.Use(authMw) // Require authentication
	{
		g.GET("/me", h.getMe)
		g.PUT("/profile", h.updateProfile)
	}
}

func (h *UserHandler) getMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.UserSvc.GetByID(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
}

func (h *UserHandler) updateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Name     string `json:"name" binding:"required,min=2"`
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Update name
	err := h.UserSvc.UpdateName(c.Request.Context(), userID.(int64), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	// Update username if provided
	if req.Username != "" {
		err = h.UserSvc.UpdateUsername(c.Request.Context(), userID.(int64), req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update username"})
			return
		}
	}

	// Update password if provided
	if req.Password != "" {
		err = h.UserSvc.UpdatePassword(c.Request.Context(), userID.(int64), req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "profile updated successfully"})
}
