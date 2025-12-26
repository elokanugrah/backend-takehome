package usecase

import (
	"context"

	"github.com/elokanugrah/backend-takehome/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

// AuthClaims contains the data stored inside the JWT.
type AuthClaims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthService defines the contract for authentication and token management.
//
//go:generate mockery --name AuthService --output ./mocks --case=snake
type AuthService interface {
	GenerateToken(ctx context.Context, user *domain.User) (string, error)
	ValidateToken(ctx context.Context, tokenString string) (*AuthClaims, error)
}

// TransactionManager defines the contract for database transaction management.
// This allows use cases to run operations within a single transaction
// without being coupled to a specific database implementation.
//
//go:generate mockery --name TransactionManager --output ./mocks --case=snake
type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
