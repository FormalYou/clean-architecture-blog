package repository

import "github.com/FormalYou/clean-architecture-blog/domain"

// TagRepository defines the interface for tag persistence.
type TagRepository interface {
	FindAll() ([]*domain.Tag, error)
	FindByName(name string) (*domain.Tag, error)
	Save(tag *domain.Tag) error
}
