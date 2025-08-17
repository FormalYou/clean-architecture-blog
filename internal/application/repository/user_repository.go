package repository

import "github.com/FormalYou/clean-architecture-blog/domain"

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	FindByID(id uint) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	Save(user *domain.User) error
	Create(user *domain.User) error
	GetByUsername(username string) (*domain.User, error)
}
