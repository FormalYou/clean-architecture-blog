package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/formal-you/clean-architecture-blog/domain"
	"github.com/formal-you/clean-architecture-blog/internal/application/repository"
	"github.com/redis/go-redis/v9"
)

type articleCacheRepository struct {
	redisClient *redis.Client
}

func NewArticleCacheRepository(redisClient *redis.Client) repository.ArticleCacheRepository {
	return &articleCacheRepository{redisClient: redisClient}
}

func (r *articleCacheRepository) GetArticle(ctx context.Context, id uint) (*domain.Article, error) {
	key := fmt.Sprintf("article:%d", id)
	val, err := r.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, err
	}

	var article domain.Article
	err = json.Unmarshal([]byte(val), &article)
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleCacheRepository) SetArticle(ctx context.Context, article *domain.Article, expiration time.Duration) error {
	key := fmt.Sprintf("article:%d", article.ID)
	val, err := json.Marshal(article)
	if err != nil {
		return err
	}
	return r.redisClient.Set(ctx, key, val, expiration).Err()
}

func (r *articleCacheRepository) GetArticles(ctx context.Context, key string) ([]*domain.Article, error) {
	val, err := r.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, err
	}

	var articles []*domain.Article
	err = json.Unmarshal([]byte(val), &articles)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *articleCacheRepository) SetArticles(ctx context.Context, key string, articles []*domain.Article, expiration time.Duration) error {
	val, err := json.Marshal(articles)
	if err != nil {
		return err
	}
	return r.redisClient.Set(ctx, key, val, expiration).Err()
}

func (r *articleCacheRepository) DeleteArticle(ctx context.Context, id uint) error {
	key := fmt.Sprintf("article:%d", id)
	return r.redisClient.Del(ctx, key).Err()
}
