package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodyLimit returns middleware that limits request body size.
func BodyLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		}
		c.Next()
	}
}
