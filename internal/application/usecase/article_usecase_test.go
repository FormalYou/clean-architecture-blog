package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/formal-you/clean-architecture-blog/domain"
	"github.com/formal-you/clean-architecture-blog/internal/application/repository"
	mock_contracts "github.com/formal-you/clean-architecture-blog/internal/application/contracts/mocks"
	mock_repo "github.com/formal-you/clean-architecture-blog/internal/application/repository/mocks"
	"github.com/formal-you/clean-architecture-blog/internal/errorx"
)

func TestArticleUsecase_CreateArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArticleRepo := mock_repo.NewMockArticleRepository(ctrl)
	mockArticleCacheRepo := mock_repo.NewMockArticleCacheRepository(ctrl)
	mockAuthSvc := mock_contracts.NewMockAuthService(ctrl)
	mockLogger := mock_contracts.NewMockLogger(ctrl)

	usecase := NewArticleUsecase(mockArticleRepo, mockArticleCacheRepo, mockAuthSvc, mockLogger)

	testCases := []struct {
		name          string
		inputArticle  *domain.Article
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			inputArticle: &domain.Article{
				Title:   "Test Article",
				Content: "This is a test content.",
			},
			setupMocks: func() {
				mockAuthSvc.EXPECT().GetUserIDFromContext(gomock.Any()).Return(int64(1), nil)
				mockArticleRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				mockLogger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedError: nil,
		},
		{
			name: "Validation Error - Missing Title",
			inputArticle: &domain.Article{
				Title:   "",
				Content: "This is a test content.",
			},
			setupMocks: func() {
				mockAuthSvc.EXPECT().GetUserIDFromContext(gomock.Any()).Return(int64(1), nil)
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any())
			},
			expectedError: errorx.New(errorx.CodeInvalidParams, errors.New("title is required")),
		},
		{
			name: "Auth Error",
			inputArticle: &domain.Article{
				Title:   "Test Article",
				Content: "This is a test content.",
			},
			setupMocks: func() {
				mockAuthSvc.EXPECT().GetUserIDFromContext(gomock.Any()).Return(int64(0), errors.New("auth error"))
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any())
			},
			expectedError: errorx.New(errorx.CodeUnauthorized, errors.New("auth error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := usecase.CreateArticle(context.Background(), tc.inputArticle)
			if tc.expectedError != nil {
				assert.Error(t, err)
				var detailErr *errorx.DetailError
				if assert.ErrorAs(t, err, &detailErr) {
					assert.Equal(t, tc.expectedError.(*errorx.DetailError).Code, detailErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestArticleUsecase_GetArticleByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArticleRepo := mock_repo.NewMockArticleRepository(ctrl)
	mockArticleCacheRepo := mock_repo.NewMockArticleCacheRepository(ctrl)
	mockAuthSvc := mock_contracts.NewMockAuthService(ctrl)
	mockLogger := mock_contracts.NewMockLogger(ctrl)

	usecase := NewArticleUsecase(mockArticleRepo, mockArticleCacheRepo, mockAuthSvc, mockLogger)

	expectedArticle := &domain.Article{ID: 1, Title: "Cached Article"}

	testCases := []struct {
		name          string
		articleID     int64
		setupMocks    func()
		expectedArticle *domain.Article
		expectedError error
	}{
		{
			name:      "Cache Hit",
			articleID: 1,
			setupMocks: func() {
				mockArticleCacheRepo.EXPECT().GetArticle(gomock.Any(), uint(1)).Return(expectedArticle, nil)
				mockLogger.EXPECT().Info(gomock.Any(), gomock.Any())
			},
			expectedArticle: expectedArticle,
			expectedError: nil,
		},
		{
			name:      "Cache Miss, DB Success",
			articleID: 2,
			setupMocks: func() {
				mockArticleCacheRepo.EXPECT().GetArticle(gomock.Any(), uint(2)).Return(nil, repository.ErrNotFound)
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any())
				mockLogger.EXPECT().Info(gomock.Any(), gomock.Any())
				mockArticleRepo.EXPECT().GetByID(gomock.Any(), int64(2)).Return(expectedArticle, nil)
				mockArticleCacheRepo.EXPECT().SetArticle(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedArticle: expectedArticle,
			expectedError: nil,
		},
		{
			name:      "Cache Miss, DB Not Found",
			articleID: 3,
			setupMocks: func() {
				mockArticleCacheRepo.EXPECT().GetArticle(gomock.Any(), uint(3)).Return(nil, repository.ErrNotFound)
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any())
				mockLogger.EXPECT().Info(gomock.Any(), gomock.Any())
				mockArticleRepo.EXPECT().GetByID(gomock.Any(), int64(3)).Return(nil, repository.ErrNotFound)
			},
			expectedArticle: nil,
			expectedError: errorx.New(errorx.CodeArticleNotFound, repository.ErrNotFound),
		},
		{
			name:      "Cache Miss, DB Error",
			articleID: 4,
			setupMocks: func() {
				mockArticleCacheRepo.EXPECT().GetArticle(gomock.Any(), uint(4)).Return(nil, repository.ErrNotFound)
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any())
				mockLogger.EXPECT().Info(gomock.Any(), gomock.Any())
				dbError := errors.New("db connection failed")
				mockArticleRepo.EXPECT().GetByID(gomock.Any(), int64(4)).Return(nil, dbError)
			},
			expectedArticle: nil,
			expectedError: errorx.New(errorx.CodeInternalServerError, errors.New("db connection failed")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			article, err := usecase.GetArticleByID(context.Background(), tc.articleID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				var detailErr *errorx.DetailError
				if assert.ErrorAs(t, err, &detailErr) {
					assert.Equal(t, tc.expectedError.(*errorx.DetailError).Code, detailErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedArticle, article)
			}
		})
	}
}

func TestArticleUsecase_UpdateArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArticleRepo := mock_repo.NewMockArticleRepository(ctrl)
	mockArticleCacheRepo := mock_repo.NewMockArticleCacheRepository(ctrl)
	mockAuthSvc := mock_contracts.NewMockAuthService(ctrl)
	mockLogger := mock_contracts.NewMockLogger(ctrl)

	usecase := NewArticleUsecase(mockArticleRepo, mockArticleCacheRepo, mockAuthSvc, mockLogger)

	existingArticle := &domain.Article{ID: 1, AuthorID: 100, Title: "Old Title", Content: "Old Content"}
	updatedArticle := &domain.Article{ID: 1, AuthorID: 100, Title: "New Title", Content: "New Content"}

	testCases := []struct {
		name          string
		inputArticle  *domain.Article
		ctxUserID     int64
		setupMocks    func()
		expectedError error
	}{
		{
			name:         "Success",
			inputArticle: updatedArticle,
			ctxUserID:    100,
			setupMocks: func() {
				mockAuthSvc.EXPECT().GetUserIDFromContext(gomock.Any()).Return(int64(100), nil)
				mockArticleRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(existingArticle, nil)
				mockArticleRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				mockArticleCacheRepo.EXPECT().DeleteArticle(gomock.Any(), uint(1)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:         "Unauthorized User",
			inputArticle: updatedArticle,
			ctxUserID:    999,
			setupMocks: func() {
				mockAuthSvc.EXPECT().GetUserIDFromContext(gomock.Any()).Return(int64(999), nil)
				mockArticleRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(existingArticle, nil)
			},
			expectedError: errorx.New(errorx.CodeUnauthorized, errors.New("user not authorized to update this article")),
		},
		{
			name:         "Article Not Found",
			inputArticle: updatedArticle,
			ctxUserID:    100,
			setupMocks: func() {
				mockAuthSvc.EXPECT().GetUserIDFromContext(gomock.Any()).Return(int64(100), nil)
				mockArticleRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(nil, repository.ErrNotFound)
			},
			expectedError: errorx.New(errorx.CodeArticleNotFound, repository.ErrNotFound),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock GetUserIDFromContext to return the desired ctxUserID
			// This is handled by the mockAuthSvc.EXPECT().GetUserIDFromContext in setupMocks

			tc.setupMocks()
			err := usecase.UpdateArticle(context.Background(), tc.inputArticle)

			if tc.expectedError != nil {
				assert.Error(t, err)
				var detailErr *errorx.DetailError
				if assert.ErrorAs(t, err, &detailErr) {
					assert.Equal(t, tc.expectedError.(*errorx.DetailError).Code, detailErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestArticleUsecase_DeleteArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArticleRepo := mock_repo.NewMockArticleRepository(ctrl)
	mockArticleCacheRepo := mock_repo.NewMockArticleCacheRepository(ctrl)
	mockAuthSvc := mock_contracts.NewMockAuthService(ctrl)
	mockLogger := mock_contracts.NewMockLogger(ctrl)

	usecase := NewArticleUsecase(mockArticleRepo, mockArticleCacheRepo, mockAuthSvc, mockLogger)

	existingArticle := &domain.Article{ID: 1, AuthorID: 100}

	testCases := []struct {
		name          string
		articleID     int64
		ctxUserID     int64
		setupMocks    func()
		expectedError error
	}{
		{
			name:      "Success",
			articleID: 1,
			ctxUserID: 100,
			setupMocks: func() {
				mockAuthSvc.EXPECT().GetUserIDFromContext(gomock.Any()).Return(int64(100), nil)
				mockArticleRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(existingArticle, nil)
				mockArticleRepo.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil)
				mockArticleCacheRepo.EXPECT().DeleteArticle(gomock.Any(), uint(1)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "Unauthorized User",
			articleID: 1,
			ctxUserID: 999,
			setupMocks: func() {
				mockAuthSvc.EXPECT().GetUserIDFromContext(gomock.Any()).Return(int64(999), nil)
				mockArticleRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(existingArticle, nil)
			},
			expectedError: errorx.New(errorx.CodeUnauthorized, errors.New("user not authorized to delete this article")),
		},
		{
			name:      "Article Not Found",
			articleID: 1,
			ctxUserID: 100,
			setupMocks: func() {
				mockAuthSvc.EXPECT().GetUserIDFromContext(gomock.Any()).Return(int64(100), nil)
				mockArticleRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(nil, repository.ErrNotFound)
			},
			expectedError: errorx.New(errorx.CodeArticleNotFound, repository.ErrNotFound),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := usecase.DeleteArticle(context.Background(), tc.articleID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				var detailErr *errorx.DetailError
				if assert.ErrorAs(t, err, &detailErr) {
					assert.Equal(t, tc.expectedError.(*errorx.DetailError).Code, detailErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
