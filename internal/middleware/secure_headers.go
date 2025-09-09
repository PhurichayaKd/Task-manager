package middleware

import "github.com/gin-gonic/gin"

func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		// หมายเหตุ: HSTS ควรเปิดเฉพาะ production ที่มี HTTPS จริงเท่านั้น
		c.Next()
	}
}
