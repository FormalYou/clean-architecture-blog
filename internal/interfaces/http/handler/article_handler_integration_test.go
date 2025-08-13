package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/formal-you/clean-architecture-blog/domain"
	"github.com/formal-you/clean-architecture-blog/internal/interfaces/http/dto"
)

func TestArticleHandler_Integration_Create(t *testing.T) {
	// ClearAllData()

	// 准备：注册和登录以获取 token
	registerReq := dto.RegisterRequest{Username: "testuser", Password: "password", Email: "test@example.com"}
	PerformRequest(TestRouter, "POST", "/users/register", ToJSON(registerReq))
	loginReq := dto.LoginRequest{Email: "test@example.com", Password: "password"}
	w := PerformRequest(TestRouter, "POST", "/users/login", ToJSON(loginReq))
	var loginResp dto.LoginResponse
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	token := loginResp.Token

	// 测试：创建文章
	createReq := dto.CreateArticleRequest{Title: "Integration Test", Content: "This is a test.", Tags: []string{"Go", "Test"}}
	w = PerformRequest(TestRouter, "POST", "/articles", ToJSON(createReq), token)

	assert.Equal(t, http.StatusCreated, w.Code)
	var articleResp domain.Article
	json.Unmarshal(w.Body.Bytes(), &articleResp)
	assert.Equal(t, "Integration Test", articleResp.Title)
}
