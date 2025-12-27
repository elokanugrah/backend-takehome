package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/elokanugrah/backend-takehome/internal/domain"
)

type PostUseCase struct {
	postRepo domain.PostRepository
}

func NewPostUseCase(pr domain.PostRepository) *PostUseCase {
	return &PostUseCase{
		postRepo: pr,
	}
}

func (uc *PostUseCase) Create(ctx context.Context, title, content string, authorID int) (*domain.Post, error) {
	if title == "" || content == "" {
		return nil, errors.New("title and content are required")
	}

	post := domain.NewPost(title, content, authorID)
	if err := uc.postRepo.Save(ctx, post); err != nil {
		return nil, err
	}

	// Fetch the full post details (including author info) to return
	return uc.GetByID(ctx, post.ID)
}

func (uc *PostUseCase) GetAll(ctx context.Context) ([]domain.Post, error) {
	return uc.postRepo.FindAll(ctx)
}

func (uc *PostUseCase) GetByID(ctx context.Context, id int) (*domain.Post, error) {
	return uc.postRepo.FindByID(ctx, id)
}

func (uc *PostUseCase) Update(ctx context.Context, id int, title, content string, authorID int) (*domain.Post, error) {
	post, err := uc.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if post.AuthorID != authorID {
		return nil, errors.New("forbidden: only the author can update this post")
	}

	if title != "" {
		post.Title = title
	}
	if content != "" {
		post.Content = content
	}
	post.UpdatedAt = time.Now()

	if err := uc.postRepo.Update(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (uc *PostUseCase) Delete(ctx context.Context, id int, authorID int) error {
	post, err := uc.postRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if post.AuthorID != authorID {
		return errors.New("forbidden: only the author can delete this post")
	}

	return uc.postRepo.Delete(ctx, id)
}
