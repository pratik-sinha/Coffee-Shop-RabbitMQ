package middleware

import "github.com/gin-gonic/gin"

// CORS CORS
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:8999")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		c.Header("Access-Control-Request-Headers", "*")

		c.Header("Access-Control-Allow-Credentials", "true")

		c.Header("Content-Type", "application/json")

		if c.Request.Method != "OPTIONS" {
			c.Next()
		} else {
			c.AbortWithStatus(200)
			return
		}
	}
}
