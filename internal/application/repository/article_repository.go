package repository

import (
	"context"

	"github.com/formal-you/clean-architecture-blog/domain"
)

// ArticleRepository 定义了文章实体的持久化操作接口
type ArticleRepository interface {
	Create(ctx context.Context, article *domain.Article) error
	GetByID(ctx context.Context, id int64) (*domain.Article, error)
	GetAll(ctx context.Context) ([]*domain.Article, error)
	Update(ctx context.Context, article *domain.Article) error
	Delete(ctx context.Context, id int64) error
}