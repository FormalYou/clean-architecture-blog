package domain

import (
	"errors"
)

// Article 是文章的领域实体
type Article struct {
	ID        int64
	Title     string
	Content   string
	AuthorID  int64
	Tags      []Tag 
	
}

// Validate 检查文章实体的业务规则
func (a *Article) Validate() error {
	if a.Title == "" {
		return errors.New("title is required")
	}
	if a.Content == "" {
		return errors.New("content is required")
	}
	if a.AuthorID == 0 {
		return errors.New("author is required")
	}
	return nil
}