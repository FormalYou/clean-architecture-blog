package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/formal-you/clean-architecture-blog/internal/application/contracts"
	"github.com/formal-you/clean-architecture-blog/internal/interfaces/http/dto"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthService provides JWT generation and validation services.
type JWTAuthService struct {
	secretKey []byte
}

// NewJWTAuthService creates a new JWTAuthService.
func NewJWTAuthService(secret string) contracts.AuthService {
	return &JWTAuthService{secretKey: []byte(secret)}
}

// GenerateToken generates a new JWT for a given user ID.
func (s *JWTAuthService) GenerateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken validates a JWT string and returns the user ID.
func (s *JWTAuthService) ValidateToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(float64); ok {
			return int64(sub), nil
		}
	}

	return 0, errors.New("invalid token")
}

// GetUserIDFromContext extracts user ID from a context.
func (s *JWTAuthService) GetUserIDFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(dto.UserIDKey).(int64)
	if !ok {
		return 0, errors.New("invalid user ID in context")
	}
	return userID, nil
}
