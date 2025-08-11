package gorm

import (
	"time"

	"github.com/formal-you/clean-architecture-blog/domain"
)

// ArticleModel 是文章在GORM中的持久化模型
type ArticleModel struct {
	ID        int64      `gorm:"primaryKey"`
	Title     string     `gorm:"not null"`
	Content   string     `gorm:"type:text"`
	AuthorID  int64      `gorm:"not null"`
	Tags      []TagModel `gorm:"many2many:article_tags;"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
}

// ToDomain 将持久化模型转换为领域模型
func (m *ArticleModel) ToDomain() *domain.Article {
	tags := make([]domain.Tag, len(m.Tags))
	for i, tagModel := range m.Tags {
		tags[i] = *tagModel.ToDomain()
	}
	return &domain.Article{
		ID:       m.ID,
		Title:    m.Title,
		Content:  m.Content,
		AuthorID: m.AuthorID,
		Tags:     tags,
	}
}

// FromDomain 将领域模型转换为持久化模型
func FromDomain(a *domain.Article) *ArticleModel {
	tags := make([]TagModel, len(a.Tags))
	for i, tag := range a.Tags {
		tags[i] = *FromDomainTag(&tag)
	}
	return &ArticleModel{
		ID:       a.ID,
		Title:    a.Title,
		Content:  a.Content,
		AuthorID: a.AuthorID,
		Tags:     tags,
		// CreatedAt: a.CreatedAt,
		// UpdatedAt: a.UpdatedAt,
	}
}
