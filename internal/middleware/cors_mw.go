package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var allowed = map[string]bool{
	"http://localhost:5173": true, // Vite
	"http://localhost:5500": true, // Live Server
	"http://127.0.0.1:5500": true, // Live Server (127)
	"https://task-manager-production-6c61.up.railway.app": true, // Railway Production
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		h := c.Writer.Header()

		// อนุญาตเฉพาะ origin ที่ใช้จริง (เพิ่มเติมได้ตามต้องการ)
		if allowed[origin] {
			h.Set("Access-Control-Allow-Origin", origin)
		}
		h.Set("Vary", "Origin")
		h.Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		h.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// ไม่ใช้ cookies จึงไม่ต้อง Allow-Credentials

		// ตอบ preflight ทันที
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
