package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name          string
		user          *User
		expectedError string
	}{
		{
			name: "Valid User",
			user: &User{
				Username: "testuser",
				Email:    "test@example.com",
			},
			expectedError: "",
		},
		{
			name: "Missing Username",
			user: &User{
				Username: "",
				Email:    "test@example.com",
			},
			expectedError: "username is required",
		},
		{
			name: "Invalid Email Format",
			user: &User{
				Username: "testuser",
				Email:    "invalid-email",
			},
			expectedError: "invalid email format",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.user.Validate()
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
