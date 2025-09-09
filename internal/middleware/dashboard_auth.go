package middleware

import (
	"log"
	"net/http"
	"strings"
	"task-manager/internal/auth"

	"github.com/gin-gonic/gin"
)

// DashboardAuthMiddleware protects dashboard routes
func DashboardAuthMiddleware(jwtAuth *auth.JWT, frontendURL string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Check for token in Authorization header first
		authHeader := c.GetHeader("Authorization")
		var token string
		
		if authHeader != "" {
			// Extract token from "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}
		
		// If no Authorization header, check cookie
		if token == "" {
			if cookie, err := c.Cookie("access_token"); err == nil {
				token = cookie
			}
		}
		
		// If still no token, redirect to login
		if token == "" {
			log.Printf("Dashboard auth: No token found, redirecting to login")
			c.Redirect(http.StatusFound, frontendURL+"/auth/login.html")
			c.Abort()
			return
		}
		
		// Validate token
		claims, err := (*jwtAuth).ParseAccess(token)
		if err != nil {
			// Invalid token, redirect to login
			log.Printf("Dashboard auth: Token validation failed: %v", err)
			c.Redirect(http.StatusFound, frontendURL+"/auth/login.html")
			c.Abort()
			return
		}
		
		// Set user ID in context for potential API calls
		c.Set("userID", claims.UserID)
		c.Next()
	})
}