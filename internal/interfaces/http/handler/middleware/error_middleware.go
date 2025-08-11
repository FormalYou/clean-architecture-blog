package middleware

import (
	"errors"

	"github.com/formal-you/clean-architecture-blog/internal/errorx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ErrorHandler is a middleware to handle errors gracefully.
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only handle errors if they exist
		if len(c.Errors) == 0 {
			return
		}

		// We only handle the first error
		err := c.Errors[0].Err

		var detailErr *errorx.DetailError
		// Check if it's our custom error
		if errors.As(err, &detailErr) {
			// Log the detailed, internal error with the appropriate level
			logFields := []zap.Field{
				zap.Int("code", detailErr.Code),
				zap.String("message", detailErr.Message),
				zap.String("request_uri", c.Request.RequestURI),
			}

			switch detailErr.LogLevel {
			case zapcore.WarnLevel:
				logger.Warn(detailErr.Error(), logFields...)
			case zapcore.ErrorLevel:
				logger.Error(detailErr.Error(), logFields...)
			default:
				logger.Error(detailErr.Error(), logFields...)
			}

			// Respond with the user-facing business error
			c.JSON(detailErr.HTTPStatus, detailErr.BusinessError)
			return
		}

		// For all other unexpected errors, log them and return a generic 500 error
		logger.Error(
			err.Error(),
			zap.String("request_uri", c.Request.RequestURI),
		)

		// Respond with a generic internal server error
		internalErr := errorx.New(errorx.CodeInternalServerError, nil)
		c.JSON(internalErr.HTTPStatus, internalErr.BusinessError)
	}
}
