package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/formal-you/clean-architecture-blog/internal/interfaces/http/dto"
	"github.com/stretchr/testify/assert"
)

func TestUserRegistrationAndLogin(t *testing.T) {
	// 1. Register a new user
	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	password := "password123"
	email := fmt.Sprintf("%s@example.com", username)

	registerBody := []byte(fmt.Sprintf(`{"username": "%s", "password": "%s", "email": "%s"}`, username, password, email))
	resp, err := newRequest("POST", "/api/v1/register", registerBody)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// 2. Log in with the new user
	loginBody := []byte(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password))
	resp, err = newRequest("POST", "/api/v1/login", loginBody)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResponse dto.LoginResponse
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(body, &loginResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, loginResponse.Token)
}