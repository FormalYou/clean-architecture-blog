package domain

import (
	"errors"
	"regexp"
)

// User 是用户的领域实体
type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Email        string
	Profile      UserProfile
}

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

// Validate 检查用户实体的业务规则
func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("username is required")
	}
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}
	return nil
}

// UserProfile 存放用户的个人资料信息
type UserProfile struct {
	Nickname string
	Avatar   string
}

// UserRepository 定义了用户数据的存储库接口
// type UserRepository interface {
// 	Create(user *User) error
// 	GetByUsername(username string) (*User, error)
// }
