package usecase_test

import (
	"context"
	"testing"

	"github.com/elokanugrah/backend-takehome/internal/domain"
	domainmocks "github.com/elokanugrah/backend-takehome/internal/domain/mocks"
	"github.com/elokanugrah/backend-takehome/internal/usecase"
	usecasemocks "github.com/elokanugrah/backend-takehome/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePost(t *testing.T) {
	ctx := context.Background()
	title := "My First Post"
	content := "Hello World"
	authorID := 1

	t.Run("success", func(t *testing.T) {
		mockPostRepo := new(domainmocks.PostRepository)
		mockTxManager := new(usecasemocks.TransactionManager)
		uc := usecase.NewPostUseCase(mockPostRepo, mockTxManager)

		// Mock Transaction
		mockTxManager.On("WithTransaction", ctx, mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		// Mock Save
		mockPostRepo.On("Save", ctx, mock.MatchedBy(func(p *domain.Post) bool {
			return p.Title == title && p.Content == content && p.AuthorID == authorID
		})).Return(nil).Run(func(args mock.Arguments) {
			p := args.Get(1).(*domain.Post)
			p.ID = 100 // Simulate DB ID generation
		})

		// Mock GetByID (called inside Create to fetch full details)
		expectedPost := &domain.Post{
			ID: 100, Title: title, Content: content, AuthorID: authorID,
			AuthorName: "John", AuthorEmail: "john@example.com",
		}
		mockPostRepo.On("FindByID", ctx, 100).Return(expectedPost, nil)

		post, err := uc.Create(ctx, title, content, authorID)

		assert.NoError(t, err)
		assert.NotNil(t, post)
		assert.Equal(t, 100, post.ID)
		assert.Equal(t, "John", post.AuthorName)

		mockPostRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		mockPostRepo := new(domainmocks.PostRepository)
		mockTxManager := new(usecasemocks.TransactionManager)
		uc := usecase.NewPostUseCase(mockPostRepo, mockTxManager)

		post, err := uc.Create(ctx, "", content, authorID)

		assert.Error(t, err)
		assert.Nil(t, post)
		assert.Equal(t, "title and content are required", err.Error())
	})
}

func TestUpdatePost(t *testing.T) {
	ctx := context.Background()
	id := 100
	authorID := 1
	title := "Updated Title"
	content := "Updated Content"

	t.Run("success", func(t *testing.T) {
		mockPostRepo := new(domainmocks.PostRepository)
		mockTxManager := new(usecasemocks.TransactionManager)
		uc := usecase.NewPostUseCase(mockPostRepo, mockTxManager)

		existingPost := &domain.Post{
			ID: id, Title: "Old Title", Content: "Old Content", AuthorID: authorID,
		}

		mockTxManager.On("WithTransaction", ctx, mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		mockPostRepo.On("FindByID", ctx, id).Return(existingPost, nil)

		mockPostRepo.On("Update", ctx, mock.MatchedBy(func(p *domain.Post) bool {
			return p.ID == id && p.Title == title && p.Content == content
		})).Return(nil)

		post, err := uc.Update(ctx, id, title, content, authorID)

		assert.NoError(t, err)
		assert.NotNil(t, post)
		assert.Equal(t, title, post.Title)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockPostRepo := new(domainmocks.PostRepository)
		mockTxManager := new(usecasemocks.TransactionManager)
		uc := usecase.NewPostUseCase(mockPostRepo, mockTxManager)

		existingPost := &domain.Post{
			ID: id, Title: "Old Title", Content: "Old Content", AuthorID: 999, // Different author
		}

		mockTxManager.On("WithTransaction", ctx, mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		mockPostRepo.On("FindByID", ctx, id).Return(existingPost, nil)

		post, err := uc.Update(ctx, id, title, content, authorID)

		assert.Error(t, err)
		assert.Nil(t, post)
		assert.Contains(t, err.Error(), "forbidden")
	})

	t.Run("post not found", func(t *testing.T) {
		mockPostRepo := new(domainmocks.PostRepository)
		mockTxManager := new(usecasemocks.TransactionManager)
		uc := usecase.NewPostUseCase(mockPostRepo, mockTxManager)

		mockTxManager.On("WithTransaction", ctx, mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		mockPostRepo.On("FindByID", ctx, id).Return(nil, domain.ErrPostNotFound)

		post, err := uc.Update(ctx, id, title, content, authorID)

		assert.Error(t, err)
		assert.Nil(t, post)
		assert.Equal(t, domain.ErrPostNotFound, err)
	})
}

func TestDeletePost(t *testing.T) {
	ctx := context.Background()
	id := 100
	authorID := 1

	t.Run("success", func(t *testing.T) {
		mockPostRepo := new(domainmocks.PostRepository)
		mockTxManager := new(usecasemocks.TransactionManager)
		uc := usecase.NewPostUseCase(mockPostRepo, mockTxManager)

		existingPost := &domain.Post{ID: id, AuthorID: authorID}

		mockTxManager.On("WithTransaction", ctx, mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		mockPostRepo.On("FindByID", ctx, id).Return(existingPost, nil)
		mockPostRepo.On("Delete", ctx, id).Return(nil)

		err := uc.Delete(ctx, id, authorID)

		assert.NoError(t, err)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockPostRepo := new(domainmocks.PostRepository)
		mockTxManager := new(usecasemocks.TransactionManager)
		uc := usecase.NewPostUseCase(mockPostRepo, mockTxManager)

		existingPost := &domain.Post{ID: id, AuthorID: 999} // Different author

		mockTxManager.On("WithTransaction", ctx, mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		mockPostRepo.On("FindByID", ctx, id).Return(existingPost, nil)

		err := uc.Delete(ctx, id, authorID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "forbidden")
	})
}
