package gorm

import (
	"github.com/formal-you/clean-architecture-blog/domain"
	"github.com/formal-you/clean-architecture-blog/internal/application/repository"
	"gorm.io/gorm"
)

// GormTagRepository 是 TagRepository 的 GORM 实现
type GormTagRepository struct {
	db *gorm.DB
}

// NewGormTagRepository 创建一个新的 GormTagRepository
func NewGormTagRepository(db *gorm.DB) repository.TagRepository {
	return &GormTagRepository{db: db}
}

// FindAll 从数据库中获取所有标签
func (r *GormTagRepository) FindAll() ([]*domain.Tag, error) {
	var tagModels []TagModel
	err := r.db.Find(&tagModels).Error
	if err != nil {
		return nil, err
	}

	var tags []*domain.Tag
	for _, model := range tagModels {
		tags = append(tags, model.ToDomain())
	}
	return tags, nil
}

// FindByName 通过名称从数据库中获取标签
func (r *GormTagRepository) FindByName(name string) (*domain.Tag, error) {
	var tagModel TagModel
	err := r.db.Where("name = ?", name).First(&tagModel).Error
	if err != nil {
		return nil, err
	}
	return tagModel.ToDomain(), nil
}

// Save 在数据库中保存一个标签
func (r *GormTagRepository) Save(tag *domain.Tag) error {
	tagModel := FromDomainTag(tag)
	return r.db.Save(tagModel).Error
}