package gorm

import (
	"github.com/FormalYou/clean-architecture-blog/domain"
	"github.com/FormalYou/clean-architecture-blog/internal/application/repository"
	"gorm.io/gorm"
)

// GormCommentRepository 是 CommentRepository 的 GORM 实现
type GormCommentRepository struct {
	db *gorm.DB
}

// NewGormCommentRepository 创建一个新的 GormCommentRepository
func NewGormCommentRepository(db *gorm.DB) repository.CommentRepository {
	return &GormCommentRepository{db: db}
}

// FindByArticleID 通过文章 ID 从数据库中获取评论
func (r *GormCommentRepository) FindByArticleID(articleID uint) ([]*domain.Comment, error) {
	var commentModels []CommentModel
	err := r.db.Where("article_id = ?", articleID).Find(&commentModels).Error
	if err != nil {
		return nil, err
	}

	var comments []*domain.Comment
	for _, model := range commentModels {
		comments = append(comments, model.ToDomain())
	}
	return comments, nil
}

// Save 在数据库中保存一条评论
func (r *GormCommentRepository) Save(comment *domain.Comment) error {
	commentModel := FromDomainComment(comment)
	return r.db.Save(commentModel).Error
}
