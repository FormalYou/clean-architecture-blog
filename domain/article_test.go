package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArticle_Validate(t *testing.T) {
	testCases := []struct {
		name          string
		article       *Article
		expectedError string
	}{
		{
			name: "Valid Article",
			article: &Article{
				Title:    "Test Title",
				Content:  "Test Content",
				AuthorID: 1,
			},
			expectedError: "",
		},
		{
			name: "Missing Title",
			article: &Article{
				Title:    "",
				Content:  "Test Content",
				AuthorID: 1,
			},
			expectedError: "title is required",
		},
		{
			name: "Missing Content",
			article: &Article{
				Title:    "Test Title",
				Content:  "",
				AuthorID: 1,
			},
			expectedError: "content is required",
		},
		{
			name: "Missing AuthorID",
			article: &Article{
				Title:    "Test Title",
				Content:  "Test Content",
				AuthorID: 0,
			},
			expectedError: "author is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.article.Validate()
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
