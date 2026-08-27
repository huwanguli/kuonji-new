package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuonji/blog/pkg/logger"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()

		logger.Info("HTTP Request",
			zap.String("path", path),
			zap.String("query", query),
			zap.String("method", method),
			zap.Int("status", status),
			zap.String("ip", clientIP),
			zap.String("user-agent", userAgent),
			zap.Duration("latency", latency),
		)
	}
}
