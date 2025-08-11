package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/formal-you/clean-architecture-blog/domain"
	mock_contracts "github.com/formal-you/clean-architecture-blog/internal/application/contracts/mocks"
	"github.com/formal-you/clean-architecture-blog/internal/application/repository"
	mock_repo "github.com/formal-you/clean-architecture-blog/internal/application/repository/mocks"
	"github.com/formal-you/clean-architecture-blog/internal/errorx"
)

func TestUserUsecase_Register(t *testing.T) {
	// Setup gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mocks
	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	mockAuthSvc := mock_contracts.NewMockAuthService(ctrl)
	mockLogger := mock_contracts.NewMockLogger(ctrl)

	// Create usecase with mocks
	userUsecase := NewUserUsecase(mockUserRepo, mockAuthSvc, 15*time.Minute, mockLogger)

	// Define test cases
	testCases := []struct {
		name          string
		inputUser     *domain.User
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			inputUser: &domain.User{
				Username:     "testuser",
				PasswordHash: "password123",
				Email:        "test@example.com",
			},
			setupMocks: func() {
				mockUserRepo.EXPECT().GetByUsername("testuser").Return(nil, repository.ErrNotFound)
				mockUserRepo.EXPECT().Create(gomock.Any()).Return(nil)
				mockLogger.EXPECT().Info(gomock.Any(), gomock.Any())
			},
			expectedError: nil,
		},
		{
			name: "User Already Exists",
			inputUser: &domain.User{
				Username:     "existinguser",
				PasswordHash: "password123",
				Email:        "existing@example.com",
			},
			setupMocks: func() {
				mockUserRepo.EXPECT().GetByUsername("existinguser").Return(&domain.User{}, nil)
			},
			expectedError: errorx.New(errorx.CodeUserAlreadyExists, nil),
		},
		{
			name: "Database error on GetByUsername",
			inputUser: &domain.User{
				Username:     "dbuser",
				PasswordHash: "password123",
				Email:        "db@example.com",
			},
			setupMocks: func() {
				dbError := errors.New("database connection failed")
				mockUserRepo.EXPECT().GetByUsername("dbuser").Return(nil, dbError)
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any())
			},
			expectedError: errorx.New(errorx.CodeInternalServerError, errors.New("database connection failed")),
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := userUsecase.Register(tc.inputUser)

			// Assertions
			if tc.expectedError != nil {
				assert.Error(t, err)
				var detailErr *errorx.DetailError
				if assert.ErrorAs(t, err, &detailErr) {
					expectedDetailErr := tc.expectedError.(*errorx.DetailError)
					assert.Equal(t, expectedDetailErr.Code, detailErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
