package repository

import (
	"context"
	"time"

	"github.com/formal-you/clean-architecture-blog/domain"
)

type ArticleCacheRepository interface {
	GetArticle(ctx context.Context, id uint) (*domain.Article, error)
	SetArticle(ctx context.Context, article *domain.Article, expiration time.Duration) error
	GetArticles(ctx context.Context, key string) ([]*domain.Article, error)
	SetArticles(ctx context.Context, key string, articles []*domain.Article, expiration time.Duration) error
	DeleteArticle(ctx context.Context, id uint) error
}