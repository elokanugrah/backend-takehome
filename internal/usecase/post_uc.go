package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/elokanugrah/backend-takehome/internal/domain"
)

type PostUseCase struct {
	postRepo  domain.PostRepository
	txManager TransactionManager
}

func NewPostUseCase(pr domain.PostRepository, tm TransactionManager) *PostUseCase {
	return &PostUseCase{
		postRepo:  pr,
		txManager: tm,
	}
}

func (uc *PostUseCase) Create(ctx context.Context, title, content string, authorID int) (*domain.Post, error) {
	if title == "" || content == "" {
		return nil, errors.New("title and content are required")
	}

	var createdPost *domain.Post

	err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		post := domain.NewPost(title, content, authorID)
		if err := uc.postRepo.Save(txCtx, post); err != nil {
			return err
		}

		// Fetch the full post details (including author info) to return
		p, err := uc.GetByID(txCtx, post.ID)
		if err != nil {
			return err
		}
		createdPost = p
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdPost, nil
}

func (uc *PostUseCase) GetAll(ctx context.Context) ([]domain.Post, error) {
	return uc.postRepo.FindAll(ctx)
}

func (uc *PostUseCase) GetByID(ctx context.Context, id int) (*domain.Post, error) {
	return uc.postRepo.FindByID(ctx, id)
}

func (uc *PostUseCase) Update(ctx context.Context, id int, title, content string, authorID int) (*domain.Post, error) {
	var updatedPost *domain.Post

	err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		post, err := uc.postRepo.FindByID(txCtx, id)
		if err != nil {
			return err
		}

		if post.AuthorID != authorID {
			return errors.New("forbidden: only the author can update this post")
		}

		if title != "" {
			post.Title = title
		}
		if content != "" {
			post.Content = content
		}
		post.UpdatedAt = time.Now()

		if err := uc.postRepo.Update(txCtx, post); err != nil {
			return err
		}

		updatedPost = post
		return nil
	})

	if err != nil {
		return nil, err
	}

	return updatedPost, nil
}

func (uc *PostUseCase) Delete(ctx context.Context, id int, authorID int) error {
	return uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		post, err := uc.postRepo.FindByID(txCtx, id)
		if err != nil {
			return err
		}

		if post.AuthorID != authorID {
			return errors.New("forbidden: only the author can delete this post")
		}

		return uc.postRepo.Delete(txCtx, id)
	})
}
