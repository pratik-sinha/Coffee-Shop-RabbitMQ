package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MaxBodyLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var w http.ResponseWriter = c.Writer
		c.Request.Body = http.MaxBytesReader(w, c.Request.Body, 5*1024*1024)
		c.Next()
	}
}
