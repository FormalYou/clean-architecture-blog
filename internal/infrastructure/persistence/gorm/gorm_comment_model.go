package gorm

import (
	"time"

	"github.com/FormalYou/clean-architecture-blog/domain"
)

// CommentModel 是评论在GORM中的持久化模型
type CommentModel struct {
	ID        int64     `gorm:"primaryKey"`
	ArticleID int64     `gorm:"not null"`
	UserID    int64     `gorm:"not null"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// ToDomain 将持久化模型转换为领域模型
func (m *CommentModel) ToDomain() *domain.Comment {
	return &domain.Comment{
		ID:        m.ID,
		ArticleID: m.ArticleID,
		UserID:    m.UserID,
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
	}
}

// FromDomainComment 将领域模型转换为持久化模型
func FromDomainComment(c *domain.Comment) *CommentModel {
	return &CommentModel{
		ID:        c.ID,
		ArticleID: c.ArticleID,
		UserID:    c.UserID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
	}
}
