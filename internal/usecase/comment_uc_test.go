package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/elokanugrah/backend-takehome/internal/domain"
	domainmocks "github.com/elokanugrah/backend-takehome/internal/domain/mocks"
	"github.com/elokanugrah/backend-takehome/internal/usecase"
	usecasemocks "github.com/elokanugrah/backend-takehome/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateComment(t *testing.T) {
	ctx := context.Background()
	postID := 100
	userID := 1
	content := "Nice article!"
	authorName := "John Doe"

	t.Run("success", func(t *testing.T) {
		mockCommentRepo := new(domainmocks.CommentRepository)
		mockPostRepo := new(domainmocks.PostRepository)
		mockUserRepo := new(domainmocks.UserRepository)
		mockTxManager := new(usecasemocks.TransactionManager)

		uc := usecase.NewCommentUseCase(mockCommentRepo, mockPostRepo, mockUserRepo, mockTxManager)

		mockTxManager.On("WithTransaction", ctx, mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		// Mock User check
		user := &domain.User{ID: userID, Name: authorName}
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)

		// Mock Post check
		post := &domain.Post{ID: postID}
		mockPostRepo.On("FindByID", ctx, postID).Return(post, nil)

		// Mock Create
		mockCommentRepo.On("Create", ctx, mock.MatchedBy(func(c *domain.Comment) bool {
			return c.PostID == postID && c.AuthorName == authorName && c.Content == content
		})).Return(nil).Run(func(args mock.Arguments) {
			c := args.Get(1).(*domain.Comment)
			c.ID = 1
		})

		comment, err := uc.Create(ctx, postID, userID, content)

		assert.NoError(t, err)
		assert.NotNil(t, comment)
		assert.Equal(t, authorName, comment.AuthorName)
		assert.Equal(t, 1, comment.ID)

		mockUserRepo.AssertExpectations(t)
		mockPostRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		mockCommentRepo := new(domainmocks.CommentRepository)
		mockPostRepo := new(domainmocks.PostRepository)
		mockUserRepo := new(domainmocks.UserRepository)
		mockTxManager := new(usecasemocks.TransactionManager)

		uc := usecase.NewCommentUseCase(mockCommentRepo, mockPostRepo, mockUserRepo, mockTxManager)

		comment, err := uc.Create(ctx, postID, userID, "")

		assert.Error(t, err)
		assert.Nil(t, comment)
		assert.Equal(t, "content is required", err.Error())
	})

	t.Run("user not found", func(t *testing.T) {
		mockCommentRepo := new(domainmocks.CommentRepository)
		mockPostRepo := new(domainmocks.PostRepository)
		mockUserRepo := new(domainmocks.UserRepository)
		mockTxManager := new(usecasemocks.TransactionManager)

		uc := usecase.NewCommentUseCase(mockCommentRepo, mockPostRepo, mockUserRepo, mockTxManager)

		mockTxManager.On("WithTransaction", ctx, mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		mockUserRepo.On("FindByID", ctx, userID).Return(nil, nil) // User nil means not found in logic

		comment, err := uc.Create(ctx, postID, userID, content)

		assert.Error(t, err)
		assert.Nil(t, comment)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("post not found", func(t *testing.T) {
		mockCommentRepo := new(domainmocks.CommentRepository)
		mockPostRepo := new(domainmocks.PostRepository)
		mockUserRepo := new(domainmocks.UserRepository)
		mockTxManager := new(usecasemocks.TransactionManager)

		uc := usecase.NewCommentUseCase(mockCommentRepo, mockPostRepo, mockUserRepo, mockTxManager)

		mockTxManager.On("WithTransaction", ctx, mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		user := &domain.User{ID: userID, Name: authorName}
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)

		mockPostRepo.On("FindByID", ctx, postID).Return(nil, domain.ErrPostNotFound)

		comment, err := uc.Create(ctx, postID, userID, content)

		assert.Error(t, err)
		assert.Nil(t, comment)
		assert.Equal(t, domain.ErrPostNotFound, err)
	})
}

func TestGetCommentsByPostID(t *testing.T) {
	ctx := context.Background()
	postID := 100

	t.Run("success", func(t *testing.T) {
		mockCommentRepo := new(domainmocks.CommentRepository)
		mockPostRepo := new(domainmocks.PostRepository)
		mockUserRepo := new(domainmocks.UserRepository)
		mockTxManager := new(usecasemocks.TransactionManager)

		uc := usecase.NewCommentUseCase(mockCommentRepo, mockPostRepo, mockUserRepo, mockTxManager)

		// Mock Post check
		post := &domain.Post{ID: postID}
		mockPostRepo.On("FindByID", ctx, postID).Return(post, nil)

		// Mock FindByPostID
		expectedComments := []domain.Comment{
			{ID: 1, PostID: postID, AuthorName: "John", Content: "Comment 1", CreatedAt: time.Now()},
			{ID: 2, PostID: postID, AuthorName: "Jane", Content: "Comment 2", CreatedAt: time.Now()},
		}
		mockCommentRepo.On("FindByPostID", ctx, postID).Return(expectedComments, nil)

		comments, err := uc.GetByPostID(ctx, postID)

		assert.NoError(t, err)
		assert.Len(t, comments, 2)
		assert.Equal(t, "John", comments[0].AuthorName)
	})

	t.Run("post not found", func(t *testing.T) {
		mockCommentRepo := new(domainmocks.CommentRepository)
		mockPostRepo := new(domainmocks.PostRepository)
		mockUserRepo := new(domainmocks.UserRepository)
		mockTxManager := new(usecasemocks.TransactionManager)

		uc := usecase.NewCommentUseCase(mockCommentRepo, mockPostRepo, mockUserRepo, mockTxManager)

		mockPostRepo.On("FindByID", ctx, postID).Return(nil, domain.ErrPostNotFound)

		comments, err := uc.GetByPostID(ctx, postID)

		assert.Error(t, err)
		assert.Nil(t, comments)
		assert.Equal(t, domain.ErrPostNotFound, err)
	})
}
