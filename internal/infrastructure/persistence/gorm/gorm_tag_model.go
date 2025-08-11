package gorm

import "github.com/formal-you/clean-architecture-blog/domain"

// TagModel 是标签在GORM中的持久化模型
type TagModel struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

// ToDomain 将持久化模型转换为领域模型
func (m *TagModel) ToDomain() *domain.Tag {
	return &domain.Tag{
		ID:   m.ID,
		Name: m.Name,
	}
}

// FromDomainTag 将领域模型转换为持久化模型
func FromDomainTag(t *domain.Tag) *TagModel {
	return &TagModel{
		ID:   t.ID,
		Name: t.Name,
	}
}