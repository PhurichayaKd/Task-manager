package middleware

import (
	"strings"
	"task-manager/internal/auth"

	"github.com/gin-gonic/gin"
)

// Authentication middleware
func Authn(j *auth.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ดึง Authorization: Bearer <token>, verify, set userID/role ใน context
	}
}

// JWTMiddleware validates JWT tokens
func JWTMiddleware(jwtAuth *auth.JWT) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := (*jwtAuth).ParseAccess(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("userID", claims.UserID)
		c.Next()
	})
}
