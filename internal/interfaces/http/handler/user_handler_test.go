package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	mock_usecase "github.com/formal-you/clean-architecture-blog/internal/application/usecase/mocks"
	"github.com/formal-you/clean-architecture-blog/internal/errorx"
	"github.com/formal-you/clean-architecture-blog/internal/interfaces/http/handler/middleware"
)

func setupRouter(userHandler *UserHandler) *gin.Engine {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		// Mock logger for middleware
		c.Set("logger", zap.NewNop())
		c.Next()
	})
	router.Use(middleware.ErrorHandler(zap.NewNop())) // Add error middleware
	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)
	return router
}

func TestUserHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUsecase := mock_usecase.NewMockUserUsecaseInterface(ctrl)
	logger := zap.NewNop() // No-op logger for tests

	handler := NewUserHandler(mockUserUsecase, logger)
	router := setupRouter(handler)

	testCases := []struct {
		name           string
		requestBody    gin.H
		setupMocks     func()
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name: "Successful Registration",
			requestBody: gin.H{
				"username": "testuser",
				"password": "password123",
				"email":    "test@example.com",
			},
			setupMocks: func() {
				mockUserUsecase.EXPECT().Register(gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   gin.H{"message": "User registered successfully"},
		},
		{
			name: "Invalid JSON",
			requestBody:    gin.H{"username": "testuser"},
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"code": float64(errorx.CodeInvalidParams), "message": "Invalid Parameters"},
		},
		{
			name: "User Already Exists",
			requestBody: gin.H{
				"username": "existinguser",
				"password": "password123",
				"email":    "existing@example.com",
			},
			setupMocks: func() {
				mockUserUsecase.EXPECT().Register(gomock.Any()).Return(errorx.New(errorx.CodeUserAlreadyExists, nil))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"code": float64(errorx.CodeUserAlreadyExists), "message": "User already exists"},
		},
		{
			name: "Internal Server Error",
			requestBody: gin.H{
				"username": "servererror",
				"password": "password123",
				"email":    "server@example.com",
			},
			setupMocks: func() {
				mockUserUsecase.EXPECT().Register(gomock.Any()).Return(errorx.New(errorx.CodeInternalServerError, errors.New("db error")))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"code": float64(errorx.CodeInternalServerError), "message": "Internal Server Error"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			// Create request
			jsonValue, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			// Record response
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			// Assertions
			assert.Equal(t, tc.expectedStatus, rr.Code)

			var responseBody gin.H
			json.Unmarshal(rr.Body.Bytes(), &responseBody)
			assert.Equal(t, tc.expectedBody, responseBody)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUsecase := mock_usecase.NewMockUserUsecaseInterface(ctrl)
	logger := zap.NewNop()

	handler := NewUserHandler(mockUserUsecase, logger)
	router := setupRouter(handler)

	testCases := []struct {
		name           string
		requestBody    gin.H
		setupMocks     func()
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name: "Successful Login",
			requestBody: gin.H{
				"username": "testuser",
				"password": "password123",
			},
			setupMocks: func() {
				mockUserUsecase.EXPECT().Login("testuser", "password123").Return("mock_jwt_token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   gin.H{"token": "mock_jwt_token"},
		},
		{
			name: "Invalid Credentials",
			requestBody: gin.H{
				"username": "wronguser",
				"password": "wrongpass",
			},
			setupMocks: func() {
				mockUserUsecase.EXPECT().Login("wronguser", "wrongpass").Return("", errorx.New(errorx.CodeInvalidCredentials, errors.New("invalid credentials")))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   gin.H{"code": float64(errorx.CodeInvalidCredentials), "message": "Invalid username or password"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			jsonValue, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)

			var responseBody gin.H
			json.Unmarshal(rr.Body.Bytes(), &responseBody)
			assert.Equal(t, tc.expectedBody, responseBody)
		})
	}
}
