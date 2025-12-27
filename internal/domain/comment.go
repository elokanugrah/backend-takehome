package domain

import (
	"context"
	"time"
)

type Comment struct {
	ID         int       `json:"id"`
	PostID     int       `json:"post_id"`
	AuthorName string    `json:"author_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

func NewComment(postID int, authorName, content string) *Comment {
	return &Comment{
		PostID:     postID,
		AuthorName: authorName,
		Content:    content,
		CreatedAt:  time.Now(),
	}
}

// CommentRepository defines the contract for comment data access.
//
//go:generate mockery --name CommentRepository --output ./mocks --case=snake
type CommentRepository interface {
	Create(ctx context.Context, comment *Comment) error
	FindByPostID(ctx context.Context, postID int) ([]Comment, error)
}
