package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		attrs := []slog.Attr{
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", status),
			slog.String("latency", latency.String()),
			slog.String("client_ip", c.ClientIP()),
		}

		level := slog.LevelInfo
		msg := "request completed"
		if status >= 500 {
			level = slog.LevelError
			msg = "request failed"
		} else if status >= 400 {
			level = slog.LevelWarn
			msg = "request error"
		}

		slog.LogAttrs(c.Request.Context(), level, msg, attrs...)
	}
}
