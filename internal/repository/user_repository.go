package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/elokanugrah/backend-takehome/internal/domain"
)

var _ domain.UserRepository = (*userRepository)(nil)

type userRepository struct {
	DB *sql.DB
}

// FindByEmail implements domain.UserRepository.
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, name, email, password_hash, created_at, updated_at 
			   FROM users WHERE email = ?`
	var u domain.User

	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil, nil to indicate not found
		}
		return nil, fmt.Errorf("error scanning user: %w", err)
	}

	return &u, nil
}

// Save implements domain.UserRepository.
func (r *userRepository) Save(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (name, email, password_hash, created_at, updated_at) 
			   VALUES (?, ?, ?, ?, ?)`

	res, err := r.DB.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}

	user.ID = int(id)
	return nil
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{DB: db}
}
