package middleware

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"test-workmate/internal/metrics"
	"time"
)

// Metrics middleware measures duration of http requests
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start)
		path := c.FullPath()

		metrics.HttpDuration.WithLabelValues(c.Request.Method, path, strconv.Itoa(c.Writer.Status())).Observe(duration.Seconds())
	}
}
