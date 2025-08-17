package usecase

import (
	"errors"
	"time"

	"github.com/FormalYou/clean-architecture-blog/domain"
	"github.com/FormalYou/clean-architecture-blog/internal/application/contracts"
	"github.com/FormalYou/clean-architecture-blog/internal/application/repository"
	"github.com/FormalYou/clean-architecture-blog/internal/errorx"
	"golang.org/x/crypto/bcrypt"
)

// UserUsecaseInterface 定义了用户相关的业务逻辑接口
type UserUsecaseInterface interface {
	Register(user *domain.User) error
	Login(email, password string) (string, error)
}

// UserUsecase 提供了用户相关的业务逻辑
type UserUsecase struct {
	userRepo   repository.UserRepository
	authSvc    contracts.AuthService
	jwtExpires time.Duration
	logger     contracts.Logger
}

// NewUserUsecase 创建一个新的 UserUsecase
func NewUserUsecase(userRepo repository.UserRepository, authSvc contracts.AuthService, jwtExpires time.Duration, logger contracts.Logger) UserUsecaseInterface {
	return &UserUsecase{
		userRepo:   userRepo,
		authSvc:    authSvc,
		jwtExpires: jwtExpires,
		logger:     logger,
	}
}

// Register 处理用户注册
func (uc *UserUsecase) Register(user *domain.User) error {
	// Check if user already exists
	_, err := uc.userRepo.FindByEmail(user.Email)
	if err == nil {
		return errorx.New(errorx.CodeUserAlreadyExists, nil)
	} else if !errors.Is(err, repository.ErrNotFound) {
		// A real database error occurred
		uc.logger.Error("failed to get user by email during registration", "error", err)
		return errorx.New(errorx.CodeInternalServerError, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		uc.logger.Error("failed to hash password", "error", err)
		return errorx.New(errorx.CodeInternalServerError, err)
	}
	user.PasswordHash = string(hashedPassword)

	err = uc.userRepo.Create(user)
	if err != nil {
		uc.logger.Error("failed to create user", "error", err)
		return errorx.New(errorx.CodeInternalServerError, err)
	}

	uc.logger.Info("user registered successfully", "username", user.Username)
	return nil
}

// Login 处理用户登录
func (uc *UserUsecase) Login(email, password string) (string, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", errorx.New(errorx.CodeInvalidCredentials, err)
		}
		uc.logger.Warn("failed to get user by email", "email", email, "error", err)
		return "", errorx.New(errorx.CodeInternalServerError, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		uc.logger.Warn("invalid password", "email", email, "error", err)
		return "", errorx.New(errorx.CodeInvalidCredentials, err)
	}

	token, err := uc.authSvc.GenerateToken(user.ID)
	if err != nil {
		uc.logger.Error("failed to generate token", "email", email, "error", err)
		return "", errorx.New(errorx.CodeInternalServerError, err)
	}

	uc.logger.Info("user logged in successfully", "email", email)
	return token, nil
}
