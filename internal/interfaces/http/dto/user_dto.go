package dto

// RegisterRequest 是用户注册请求的 DTO
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// LoginRequest 是用户登录请求的 DTO
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 是用户登录响应的 DTO
type LoginResponse struct {
	Token string `json:"token"`
}