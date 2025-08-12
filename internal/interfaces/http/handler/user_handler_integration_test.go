package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/formal-you/clean-architecture-blog/internal/interfaces/http/dto"
)

func TestUserHandler_Integration_Register(t *testing.T) {
	ClearAllData()

	registerReq := dto.RegisterRequest{Username: "newuser", Password: "password", Email: "new@example.com"}
	w := PerformRequest(TestRouter, "POST", "/users/register", ToJSON(registerReq))

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestUserHandler_Integration_Login(t *testing.T) {
	ClearAllData()

	// 准备: 先注册一个用户
	registerReq := dto.RegisterRequest{Username: "loginuser", Password: "password", Email: "login@example.com"}
	PerformRequest(TestRouter, "POST", "/users/register", ToJSON(registerReq))

	// 测试: 登录
	loginReq := dto.LoginRequest{Email: "login@example.com", Password: "password"}
	w := PerformRequest(TestRouter, "POST", "/users/login", ToJSON(loginReq))

	assert.Equal(t, http.StatusOK, w.Code)
	var loginResp dto.LoginResponse
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NotEmpty(t, loginResp.Token)
}