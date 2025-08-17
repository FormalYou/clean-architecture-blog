package repository

import "github.com/FormalYou/clean-architecture-blog/domain"

// CommentRepository defines the interface for comment persistence.
type CommentRepository interface {
	FindByArticleID(articleID uint) ([]*domain.Comment, error)
	Save(comment *domain.Comment) error
}
