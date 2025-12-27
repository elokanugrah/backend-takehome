package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/elokanugrah/backend-takehome/internal/domain"
)

var _ domain.PostRepository = (*postRepository)(nil)

type postRepository struct {
	DB *sql.DB
}

func NewPostRepository(db *sql.DB) *postRepository {
	return &postRepository{DB: db}
}

func (r *postRepository) getExecutor(ctx context.Context) dbExecutor {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return r.DB
}

func (r *postRepository) Save(ctx context.Context, post *domain.Post) error {
	query := `INSERT INTO posts (title, content, author_id, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?)`

	res, err := r.getExecutor(ctx).ExecContext(ctx, query,
		post.Title,
		post.Content,
		post.AuthorID,
		post.CreatedAt,
		post.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error saving post: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}

	post.ID = int(id)
	return nil
}

func (r *postRepository) FindAll(ctx context.Context) ([]domain.Post, error) {
	query := `SELECT p.id, p.title, p.content, p.author_id, u.name, u.email, p.created_at, p.updated_at 
			  FROM posts p
			  JOIN users u ON p.author_id = u.id`

	rows, err := r.getExecutor(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %w", err)
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var p domain.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.AuthorName, &p.AuthorEmail, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning post: %w", err)
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts rows: %w", err)
	}

	return posts, nil
}

func (r *postRepository) FindByID(ctx context.Context, id int) (*domain.Post, error) {
	query := `SELECT p.id, p.title, p.content, p.author_id, u.name, u.email, p.created_at, p.updated_at 
			  FROM posts p
			  JOIN users u ON p.author_id = u.id WHERE p.id = ?`

	var p domain.Post
	err := r.getExecutor(ctx).QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.AuthorName, &p.AuthorEmail, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPostNotFound
		}
		return nil, fmt.Errorf("error scanning post: %w", err)
	}

	return &p, nil
}

func (r *postRepository) Update(ctx context.Context, post *domain.Post) error {
	query := `UPDATE posts SET title = ?, content = ?, updated_at = ? WHERE id = ?`

	res, err := r.getExecutor(ctx).ExecContext(ctx, query, post.Title, post.Content, post.UpdatedAt, post.ID)
	if err != nil {
		return fmt.Errorf("error updating post: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrPostNotFound
	}

	return nil
}

func (r *postRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM posts WHERE id = ?`

	res, err := r.getExecutor(ctx).ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting post: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrPostNotFound
	}

	return nil
}
