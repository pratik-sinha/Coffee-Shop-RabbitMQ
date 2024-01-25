package middleware

import (
	"coffee-shop/pkg/metrics"
	"time"

	"github.com/gin-gonic/gin"
)

func RecordMetrics(m metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		timeDur := time.Since(start).Milliseconds()
		statusCode := c.Writer.Status()

		m.IncHits(c.Request.Context(), statusCode, c.Request.Method, c.Request.URL.Path)
		m.ObserveResponseTime(c.Request.Context(), statusCode, c.Request.Method, c.Request.URL.Path, float64(timeDur))

	}
}
