package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	uuid := uuid.New()
	return func(c *gin.Context) {
		c.Header("X-Request-ID", uuid.String())
		c.Next()
	}
}
