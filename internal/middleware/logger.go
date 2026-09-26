package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"tmp/internal/logger"
)

func GinLogger(baseLog *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		reqLog := baseLog.With(
			slog.String("request_id", reqID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
		)
		ctx := logger.WithLogger(c.Request.Context(), reqLog)
		c.Request = c.Request.WithContext(ctx)
		c.Header("X-Request-ID", reqID)

		c.Next()
		reqLog.Info("http request served",
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", time.Since(start)),
			slog.String("ip", c.ClientIP()),
		)
	}
}