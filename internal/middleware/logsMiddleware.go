package middleware

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type loggerKey int

const (
	LoggerContextKey loggerKey = iota
)

func StructuredLogger(serviceName string) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	logger.SetOutput(os.Stdout) // лучше для Docker/Loki
	logger.SetLevel(logrus.InfoLevel)
	logger.WithFields(logrus.Fields{
		"service": serviceName,
		"env":     os.Getenv("ENV"),
		"version": os.Getenv("VERSION"),
	})

	return logger
}

func LogMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ctx := c.Request.Context()
		// ctx = context.WithValue(ctx, LoggerContextKey, logger)
		// c.Request = c.Request.WithContext(ctx)
		// c.Next()

		start := time.Now()
		requestLogger := logger.WithFields(logrus.Fields{
			"request_id": c.GetHeader("X-Request-ID"),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		ctx := context.WithValue(c.Request.Context(), LoggerContextKey, requestLogger)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		fields := logrus.Fields{
			"status":         status,
			"duration":       duration.Milliseconds(),
			"duration_human": duration.String(),
			"bytes":          c.Writer.Size(),
		}

		if len(c.Errors) > 0 {
			fields["error"] = c.Errors.String()
		}

		switch {
		case status >= 500:
			requestLogger.WithFields(fields).Error("Server error")
		case status >= 400:
			requestLogger.WithFields(fields).Warn("Client error")
		default:
			requestLogger.WithFields(fields).Info("Request completed")
		}
	}
}

func GetLoggerFromContext(ctx context.Context) *logrus.Entry {
	if logger, ok := ctx.Value(LoggerContextKey).(*logrus.Entry); ok {
		return logger
	}
	return logrus.NewEntry(logrus.StandardLogger()).WithFields(logrus.Fields{
		"warning": "logger_not_found_in_context",
	})
}
