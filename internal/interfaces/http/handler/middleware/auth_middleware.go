package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/FormalYou/clean-architecture-blog/internal/application/contracts"
	"github.com/FormalYou/clean-architecture-blog/internal/interfaces/http/dto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthMiddleware creates a Gin middleware for JWT authentication.
func AuthMiddleware(authSvc contracts.AuthService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("authorization header is missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("invalid token format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}

		userID, err := authSvc.ValidateToken(parts[1])
		if err != nil {
			logger.Warn("invalid token", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		logger.Info("user authenticated", zap.Int64("user_id", userID))
		ctx := context.WithValue(c.Request.Context(), dto.UserIDKey, userID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
