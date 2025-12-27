package usecase

import (
	"context"
	"errors"

	"github.com/elokanugrah/backend-takehome/internal/domain"
)

type CommentUseCase struct {
	commentRepo domain.CommentRepository
	postRepo    domain.PostRepository
	userRepo    domain.UserRepository
	txManager   TransactionManager
}

func NewCommentUseCase(cr domain.CommentRepository, pr domain.PostRepository, ur domain.UserRepository, tm TransactionManager) *CommentUseCase {
	return &CommentUseCase{
		commentRepo: cr,
		postRepo:    pr,
		userRepo:    ur,
		txManager:   tm,
	}
}

func (uc *CommentUseCase) Create(ctx context.Context, postID int, userID int, content string) (*domain.Comment, error) {
	if content == "" {
		return nil, errors.New("content is required")
	}

	var createdComment *domain.Comment

	err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// Get User to set AuthorName
		user, err := uc.userRepo.FindByID(txCtx, userID)
		if err != nil {
			return err
		}
		if user == nil {
			return errors.New("user not found")
		}

		// Ensure post exists
		if _, err := uc.postRepo.FindByID(txCtx, postID); err != nil {
			return err
		}

		comment := domain.NewComment(postID, user.Name, content)
		if err := uc.commentRepo.Create(txCtx, comment); err != nil {
			return err
		}

		createdComment = comment
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdComment, nil
}

func (uc *CommentUseCase) GetByPostID(ctx context.Context, postID int) ([]domain.Comment, error) {
	// Verify post exists first
	if _, err := uc.postRepo.FindByID(ctx, postID); err != nil {
		return nil, err
	}

	return uc.commentRepo.FindByPostID(ctx, postID)
}
