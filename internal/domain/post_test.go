package domain_test

import (
	"testing"

	"github.com/elokanugrah/backend-takehome/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewPost(t *testing.T) {
	title := "My Post"
	content := "Content"
	authorID := 1

	post := domain.NewPost(title, content, authorID)

	assert.NotNil(t, post)
	assert.Equal(t, title, post.Title)
	assert.Equal(t, content, post.Content)
	assert.Equal(t, authorID, post.AuthorID)
	assert.False(t, post.CreatedAt.IsZero())
	assert.False(t, post.UpdatedAt.IsZero())
}
