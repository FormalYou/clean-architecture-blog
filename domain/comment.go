package domain

import "time"

// Comment 是文章评论的领域实体
type Comment struct {
	ID        int64
	ArticleID int64
	UserID    int64
	Content   string
	CreatedAt time.Time
}