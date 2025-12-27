package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/elokanugrah/backend-takehome/internal/domain"
)

var _ domain.CommentRepository = (*commentRepository)(nil)

type commentRepository struct {
	DB *sql.DB
}

func NewCommentRepository(db *sql.DB) *commentRepository {
	return &commentRepository{DB: db}
}

func (r *commentRepository) getExecutor(ctx context.Context) dbExecutor {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return r.DB
}

func (r *commentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	query := `INSERT INTO comments (post_id, author_name, content, created_at) 
			  VALUES (?, ?, ?, ?)`

	res, err := r.getExecutor(ctx).ExecContext(ctx, query,
		comment.PostID,
		comment.AuthorName,
		comment.Content,
		comment.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("error saving comment: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}

	comment.ID = int(id)
	return nil
}

func (r *commentRepository) FindByPostID(ctx context.Context, postID int) ([]domain.Comment, error) {
	query := `SELECT id, post_id, author_name, content, created_at FROM comments WHERE post_id = ?`

	rows, err := r.getExecutor(ctx).QueryContext(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("error querying comments: %w", err)
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var c domain.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.AuthorName, &c.Content, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning comment: %w", err)
		}
		comments = append(comments, c)
	}

	return comments, nil
}
