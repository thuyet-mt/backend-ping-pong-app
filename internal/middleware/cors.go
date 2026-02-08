package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {

		// ⚠️ Trong giai đoạn dev, cho phép tất cả
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// Nếu sau này cần auth bằng cookie / credentials
		// thì KHÔNG được dùng "*" nữa

		c.Writer.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, PUT, PATCH, DELETE, OPTIONS",
		)

		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Origin, Content-Type, Authorization",
		)

		c.Writer.Header().Set(
			"Access-Control-Expose-Headers",
			"Content-Length",
		)

		// Preflight request
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
