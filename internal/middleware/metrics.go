package middleware

import (
	"time"

	"datadog-exercise/platform/metrics"

	"github.com/gin-gonic/gin"
)

func Metrics(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		metrics.IncInFlight(c.Request.Context())

		c.Next()

		metrics.DecInFlight(c.Request.Context())
		duration := time.Since(start).Seconds()
		metrics.RecordHTTPRequest(
			c.Request.Context(),
			c.Request.Method,
			c.FullPath(),
			c.Writer.Status(),
			duration,
		)
	}
}
