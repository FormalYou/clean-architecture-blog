package repository

import "github.com/formal-you/clean-architecture-blog/domain"

// TagRepository defines the interface for tag persistence.
type TagRepository interface {
	FindAll() ([]*domain.Tag, error)
	FindByName(name string) (*domain.Tag, error)
	Save(tag *domain.Tag) error
}