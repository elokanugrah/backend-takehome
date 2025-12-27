package domain_test

import (
	"testing"

	"github.com/elokanugrah/backend-takehome/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewComment(t *testing.T) {
	postID := 1
	authorName := "John Doe"
	content := "Great post!"

	comment := domain.NewComment(postID, authorName, content)

	assert.NotNil(t, comment)
	assert.Equal(t, postID, comment.PostID)
	assert.Equal(t, authorName, comment.AuthorName)
	assert.Equal(t, content, comment.Content)
	assert.False(t, comment.CreatedAt.IsZero())
}
