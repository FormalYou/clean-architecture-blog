package gorm

import (
	"errors"

	"github.com/FormalYou/clean-architecture-blog/domain"
	"github.com/FormalYou/clean-architecture-blog/internal/application/repository"
	"gorm.io/gorm"
)

// GormUserRepository 是 UserRepository 的 GORM 实现
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository 创建一个新的 GormUserRepository
func NewGormUserRepository(db *gorm.DB) repository.UserRepository {
	return &GormUserRepository{db: db}
}

// FindByID 通过 ID 从数据库中获取用户
func (r *GormUserRepository) FindByID(id uint) (*domain.User, error) {
	var userModel UserModel
	err := r.db.First(&userModel, id).Error
	if err != nil {
		return nil, err
	}
	return userModel.ToDomain(), nil
}

// FindByEmail 通过 Email 从数据库中获取用户
func (r *GormUserRepository) FindByEmail(email string) (*domain.User, error) {
	var userModel UserModel
	err := r.db.Where("email = ?", email).First(&userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return userModel.ToDomain(), nil
}

// GetByUsername 通过用户名从数据库中获取用户
func (r *GormUserRepository) GetByUsername(username string) (*domain.User, error) {
	var userModel UserModel
	err := r.db.Where("username = ?", username).First(&userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return userModel.ToDomain(), nil
}

// Save 在数据库中保存一个用户
func (r *GormUserRepository) Save(user *domain.User) error {
	userModel := FromDomainUser(user)
	return r.db.Save(userModel).Error
}

// Create 在数据库中创建一个新用户
func (r *GormUserRepository) Create(user *domain.User) error {
	userModel := FromDomainUser(user)
	return r.db.Create(userModel).Error
}
