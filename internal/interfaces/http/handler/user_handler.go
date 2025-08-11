package handler

import (
	"net/http"

	"github.com/formal-you/clean-architecture-blog/domain"
	"github.com/formal-you/clean-architecture-blog/internal/application/usecase"
	"github.com/formal-you/clean-architecture-blog/internal/errorx"
	"github.com/formal-you/clean-architecture-blog/internal/interfaces/http/dto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler 处理用户相关的 HTTP 请求
type UserHandler struct {
	userUsecase usecase.UserUsecaseInterface
	logger      *zap.Logger
}

// NewUserHandler 创建一个新的 UserHandler
func NewUserHandler(userUsecase usecase.UserUsecaseInterface, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		logger:      logger.Named("UserHandler"),
	}
}

// Register 处理用户注册请求
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errorx.New(errorx.CodeInvalidParams, err))
		return
	}

	user := &domain.User{
		Username:     req.Username,
		PasswordHash: req.Password,
		Email:        req.Email,
	}

	if err := h.userUsecase.Register(user); err != nil {
		c.Error(err) // Pass the error to the middleware
		return
	}

	h.logger.Info("user registered successfully", zap.String("username", req.Username))
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login 处理用户登录请求
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errorx.New(errorx.CodeInvalidParams, err))
		return
	}

	token, err := h.userUsecase.Login(req.Username, req.Password)
	if err != nil {
		c.Error(err) // Pass the error to the middleware
		return
	}

	h.logger.Info("user logged in successfully", zap.String("username", req.Username))
	c.JSON(http.StatusOK, dto.LoginResponse{Token: token})
}