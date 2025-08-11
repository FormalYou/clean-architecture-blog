package gorm

import (
	"time"

	"github.com/formal-you/clean-architecture-blog/domain"
)

// UserModel 是用户在GORM中的持久化模型
type UserModel struct {
	ID           int64  `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Email        string `gorm:"unique;not null"`
	Nickname     string
	Avatar       string
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

// ToDomain 将持久化模型转换为领域模型
func (m *UserModel) ToDomain() *domain.User {
	return &domain.User{
		ID:           m.ID,
		Username:     m.Username,
		PasswordHash: m.PasswordHash,
		Email:        m.Email,
		Profile: domain.UserProfile{
			Nickname: m.Nickname,
			Avatar:   m.Avatar,
		},
		// CreatedAt: m.CreatedAt,
		// UpdatedAt: m.UpdatedAt,
	}
}

// FromDomainUser 将领域模型转换为持久化模型
func FromDomainUser(u *domain.User) *UserModel {
	return &UserModel{
		ID:           u.ID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		Email:        u.Email,
		Nickname:     u.Profile.Nickname,
		Avatar:       u.Profile.Avatar,
		// CreatedAt:    u.CreatedAt,
		// UpdatedAt:    u.UpdatedAt,
	}
}
