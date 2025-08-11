package gorm

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/formal-you/clean-architecture-blog/domain"
	"github.com/formal-you/clean-architecture-blog/internal/application/repository"
)

// GormArticleRepository 是 ArticleRepository 的 GORM 实现
type GormArticleRepository struct {
	db *gorm.DB
}

// NewGormArticleRepository 创建一个新的 GormArticleRepository
func NewGormArticleRepository(db *gorm.DB) repository.ArticleRepository {
	return &GormArticleRepository{db: db}
}

func (r *GormArticleRepository) Create(ctx context.Context, article *domain.Article) error {
	articleModel := FromDomain(article)
	return r.db.WithContext(ctx).Create(articleModel).Error
}

func (r *GormArticleRepository) GetByID(ctx context.Context, id int64) (*domain.Article, error) {
	var articleModel ArticleModel
	if err := r.db.WithContext(ctx).First(&articleModel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return articleModel.ToDomain(), nil
}

func (r *GormArticleRepository) GetAll(ctx context.Context) ([]*domain.Article, error) {
	var articleModels []ArticleModel
	if err := r.db.WithContext(ctx).Find(&articleModels).Error; err != nil {
		return nil, err
	}
	var articles []*domain.Article
	for _, model := range articleModels {
		articles = append(articles, model.ToDomain())
	}
	return articles, nil
}

func (r *GormArticleRepository) Update(ctx context.Context, article *domain.Article) error {
	articleModel := FromDomain(article)
	return r.db.WithContext(ctx).Save(articleModel).Error
}

func (r *GormArticleRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&ArticleModel{}, id).Error
}