package gorm

import (
	"time"

	"github.com/FormalYou/clean-architecture-blog/domain"
)

// ArticleModel 是文章在GORM中的持久化模型
type ArticleModel struct {
	ID        int64     `gorm:"primaryKey"`
	Title     string    `gorm:"not null"`
	Content   string    `gorm:"type:text"`
	AuthorID  int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// ToDomain 将持久化模型转换为领域模型
func (m *ArticleModel) ToDomain() *domain.Article {
	return &domain.Article{
		ID:       m.ID,
		Title:    m.Title,
		Content:  m.Content,
		AuthorID: m.AuthorID,
	}
}

// FromDomain 将领域模型转换为持久化模型
func FromDomain(a *domain.Article) *ArticleModel {
	return &ArticleModel{
		ID:       a.ID,
		Title:    a.Title,
		Content:  a.Content,
		AuthorID: a.AuthorID,
	}
}
