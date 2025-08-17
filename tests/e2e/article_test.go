package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/FormalYou/clean-architecture-blog/internal/interfaces/http/dto"
	"github.com/stretchr/testify/assert"
)

func TestArticleWorkflow(t *testing.T) {
	var token string
	var articleID uint

	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	password := "password123"
	email := fmt.Sprintf("%s@example.com", username)

	t.Run("User Registration", func(t *testing.T) {
		registerBody := []byte(fmt.Sprintf(`{"username": "%s", "password": "%s", "email": "%s"}`, username, password, email))
		resp, err := newRequest("POST", "/api/v1/register", registerBody)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("User Login", func(t *testing.T) {
		loginBody := []byte(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password))
		resp, err := newRequest("POST", "/api/v1/login", loginBody)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var loginResponse dto.LoginResponse
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(body, &loginResponse)
		assert.NoError(t, err)
		token = loginResponse.Token
		assert.NotEmpty(t, token)
	})

	t.Run("Create Article", func(t *testing.T) {
		articleTitle := "Test Article"
		articleContent := "This is a test article."
		createArticleBody := []byte(fmt.Sprintf(`{"title": "%s", "content": "%s"}`, articleTitle, articleContent))

		resp, err := newRequestWithAuth("POST", "/api/v1/articles", createArticleBody, token)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var createArticleResponse dto.ArticleResponse
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(body, &createArticleResponse)
		assert.NoError(t, err)
		articleID = createArticleResponse.ID
		assert.NotZero(t, articleID)
	})

	t.Run("Get Article", func(t *testing.T) {
		resp, err := newRequest("GET", fmt.Sprintf("/api/v1/articles/%d", articleID), nil)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var getArticleResponse dto.ArticleResponse
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(body, &getArticleResponse)
		assert.NoError(t, err)
		assert.Equal(t, "Test Article", getArticleResponse.Title)
		assert.Equal(t, "This is a test article.", getArticleResponse.Content)
	})
}
