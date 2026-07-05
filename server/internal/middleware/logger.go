package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"starledger/internal/pkg"
)

// Logger is the request logging middleware.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		if pkg.Logger != nil {
			pkg.Logger.Info("request",
				zap.Int("status", status),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.Duration("latency", latency),
				zap.String("client_ip", c.ClientIP()),
				zap.Int("body_size", c.Writer.Size()),
			)
		}
	}
}
