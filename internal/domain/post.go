package domain

import (
	"context"
	"errors"
	"time"
)

var ErrPostNotFound = errors.New("post not found")

type Post struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	AuthorID    int       `json:"author_id"`
	AuthorName  string    `json:"author_name"`
	AuthorEmail string    `json:"author_email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewPost(title, content string, authorID int) *Post {
	return &Post{
		Title:     title,
		Content:   content,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// PostRepository defines the contract for post data access.
//
//go:generate mockery --name PostRepository --output ./mocks --case=snake
type PostRepository interface {
	Save(ctx context.Context, post *Post) error
	FindAll(ctx context.Context) ([]Post, error)
	FindByID(ctx context.Context, id int) (*Post, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id int) error
}
