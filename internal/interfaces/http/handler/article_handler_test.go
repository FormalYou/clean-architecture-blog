package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/formal-you/clean-architecture-blog/domain"
	"github.com/formal-you/clean-architecture-blog/internal/application/repository"
	mock_usecase "github.com/formal-you/clean-architecture-blog/internal/application/usecase/mocks"
	"github.com/formal-you/clean-architecture-blog/internal/errorx"
	"github.com/formal-you/clean-architecture-blog/internal/interfaces/http/handler/middleware"
)

func setupArticleRouter(articleHandler *ArticleHandler) *gin.Engine {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		// Mock logger for middleware
		c.Set("logger", zap.NewNop())
		c.Next()
	})
	router.Use(middleware.ErrorHandler(zap.NewNop())) // Add error middleware
	router.POST("/articles", articleHandler.Create)
	router.GET("/articles", articleHandler.GetAll)
	router.GET("/articles/:id", articleHandler.GetByID)
	router.PUT("/articles/:id", articleHandler.Update)
	router.DELETE("/articles/:id", articleHandler.Delete)
	return router
}

func TestArticleHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArticleUsecase := mock_usecase.NewMockArticleUsecaseInterface(ctrl)
	logger := zap.NewNop()

	handler := NewArticleHandler(mockArticleUsecase, logger)
	router := setupArticleRouter(handler)

	testCases := []struct {
		name           string
		requestBody    gin.H
		setupMocks     func()
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name: "Success",
			requestBody: gin.H{
				"title":   "Test Article",
				"content": "This is a test content.",
				"tags":    []string{"go", "test"},
			},
			setupMocks: func() {
				mockArticleUsecase.EXPECT().CreateArticle(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   gin.H{"ID": float64(0), "Title": "Test Article", "Content": "This is a test content.", "AuthorID": float64(0), "Tags": []interface{}{map[string]interface{}{"ID": float64(0), "Name": "go"}, map[string]interface{}{"ID": float64(0), "Name": "test"}}},
		},
		{
			name:           "Invalid JSON",
			requestBody:    gin.H{"title": "Test Article"},
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"code": float64(errorx.CodeInvalidParams), "message": "Invalid Parameters"},
		},
		{
			name: "Usecase Error",
			requestBody: gin.H{
				"title":   "Test Article",
				"content": "This is a test content.",
				"tags":    []string{"go", "test"},
			},
			setupMocks: func() {
				mockArticleUsecase.EXPECT().CreateArticle(gomock.Any(), gomock.Any()).Return(errorx.New(errorx.CodeInternalServerError, errors.New("db error")))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"code": float64(errorx.CodeInternalServerError), "message": "Internal Server Error"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			jsonValue, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(jsonValue))
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

func TestArticleHandler_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArticleUsecase := mock_usecase.NewMockArticleUsecaseInterface(ctrl)
	logger := zap.NewNop()

	handler := NewArticleHandler(mockArticleUsecase, logger)
	router := setupArticleRouter(handler)

	expectedArticle := &domain.Article{ID: 1, Title: "Test Article", Content: "Test Content"}

	testCases := []struct {
		name           string
		articleID      string
		setupMocks     func()
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name:      "Success",
			articleID: "1",
			setupMocks: func() {
				mockArticleUsecase.EXPECT().GetArticleByID(gomock.Any(), int64(1)).Return(expectedArticle, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   gin.H{"ID": float64(1), "Title": "Test Article", "Content": "Test Content", "AuthorID": float64(0), "Tags": interface{}(nil)},
		},
		{
			name:           "Invalid ID",
			articleID:      "abc",
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"code": float64(errorx.CodeInvalidParams), "message": "Invalid Parameters"},
		},
		{
			name:      "Article Not Found",
			articleID: "2",
			setupMocks: func() {
				mockArticleUsecase.EXPECT().GetArticleByID(gomock.Any(), int64(2)).Return(nil, errorx.New(errorx.CodeArticleNotFound, repository.ErrNotFound))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   gin.H{"code": float64(errorx.CodeArticleNotFound), "message": "Article not found"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/articles/%s", tc.articleID), nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)

			var responseBody gin.H
			json.Unmarshal(rr.Body.Bytes(), &responseBody)
			assert.Equal(t, tc.expectedBody, responseBody)
		})
	}
}

func TestArticleHandler_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArticleUsecase := mock_usecase.NewMockArticleUsecaseInterface(ctrl)
	logger := zap.NewNop()

	handler := NewArticleHandler(mockArticleUsecase, logger)
	router := setupArticleRouter(handler)

	expectedArticles := []*domain.Article{
		{ID: 1, Title: "Article 1"},
		{ID: 2, Title: "Article 2"},
	}

	testCases := []struct {
		name           string
		setupMocks     func()
		expectedStatus int
		expectedBody   interface{} // Use interface{} to handle different body types
	}{
		{
			name: "Success",
			setupMocks: func() {
				mockArticleUsecase.EXPECT().GetAllArticles(gomock.Any()).Return(expectedArticles, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   []gin.H{{"ID": float64(1), "Title": "Article 1", "AuthorID": float64(0), "Content": "", "Tags": interface{}(nil)}, {"ID": float64(2), "Title": "Article 2", "AuthorID": float64(0), "Content": "", "Tags": interface{}(nil)}},
		},
		{
			name: "Usecase Error",
			setupMocks: func() {
				mockArticleUsecase.EXPECT().GetAllArticles(gomock.Any()).Return(nil, errorx.New(errorx.CodeInternalServerError, errors.New("db error")))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"code": float64(errorx.CodeInternalServerError), "message": "Internal Server Error"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			req, _ := http.NewRequest(http.MethodGet, "/articles", nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)

			switch expected := tc.expectedBody.(type) {
			case []gin.H:
				var responseBody []gin.H
				err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.Equal(t, expected, responseBody)
			case gin.H:
				var responseBody gin.H
				err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.Equal(t, expected, responseBody)
			}
		})
	}
}

func TestArticleHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArticleUsecase := mock_usecase.NewMockArticleUsecaseInterface(ctrl)
	logger := zap.NewNop()

	handler := NewArticleHandler(mockArticleUsecase, logger)
	router := setupArticleRouter(handler)

	testCases := []struct {
		name           string
		articleID      string
		requestBody    gin.H
		setupMocks     func()
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name:      "Success",
			articleID: "1",
			requestBody: gin.H{
				"title":   "Updated Title",
				"content": "Updated Content",
			},
			setupMocks: func() {
				mockArticleUsecase.EXPECT().UpdateArticle(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   gin.H{"status": "article updated"},
		},
		{
			name:      "Invalid ID",
			articleID: "abc",
			requestBody: gin.H{
				"title":   "Updated Title",
				"content": "Updated Content",
			},
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"code": float64(errorx.CodeInvalidParams), "message": "Invalid Parameters"},
		},
		{
			name:      "Usecase Error",
			articleID: "1",
			requestBody: gin.H{
				"title":   "Updated Title",
				"content": "Updated Content",
			},
			setupMocks: func() {
				mockArticleUsecase.EXPECT().UpdateArticle(gomock.Any(), gomock.Any()).Return(errorx.New(errorx.CodeInternalServerError, errors.New("db error")))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"code": float64(errorx.CodeInternalServerError), "message": "Internal Server Error"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			jsonValue, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/articles/%s", tc.articleID), bytes.NewBuffer(jsonValue))
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

func TestArticleHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArticleUsecase := mock_usecase.NewMockArticleUsecaseInterface(ctrl)
	logger := zap.NewNop()

	handler := NewArticleHandler(mockArticleUsecase, logger)
	router := setupArticleRouter(handler)

	testCases := []struct {
		name           string
		articleID      string
		setupMocks     func()
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name:      "Success",
			articleID: "1",
			setupMocks: func() {
				mockArticleUsecase.EXPECT().DeleteArticle(gomock.Any(), int64(1)).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   gin.H{"status": "article deleted"},
		},
		{
			name:           "Invalid ID",
			articleID:      "abc",
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"code": float64(errorx.CodeInvalidParams), "message": "Invalid Parameters"},
		},
		{
			name:      "Usecase Error",
			articleID: "1",
			setupMocks: func() {
				mockArticleUsecase.EXPECT().DeleteArticle(gomock.Any(), int64(1)).Return(errorx.New(errorx.CodeInternalServerError, errors.New("db error")))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"code": float64(errorx.CodeInternalServerError), "message": "Internal Server Error"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/articles/%s", tc.articleID), nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)

			var responseBody gin.H
			json.Unmarshal(rr.Body.Bytes(), &responseBody)
			assert.Equal(t, tc.expectedBody, responseBody)
		})
	}
}
