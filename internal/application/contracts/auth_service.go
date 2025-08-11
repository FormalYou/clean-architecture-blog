package contracts

import "context"

// AuthService defines the interface for authentication services.
type AuthService interface {
	GenerateToken(userID int64) (string, error)
	ValidateToken(tokenString string) (int64, error)
	GetUserIDFromContext(ctx context.Context) (int64, error)
}