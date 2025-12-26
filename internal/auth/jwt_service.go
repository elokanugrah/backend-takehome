package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/elokanugrah/backend-takehome/internal/domain"
	"github.com/elokanugrah/backend-takehome/internal/usecase"
	"github.com/golang-jwt/jwt/v5"
)

// JWTService provides an implementation of the AuthService using JWT.
type JWTService struct {
	secretKey []byte
	issuer    string
}

func NewJWTService(secret string) usecase.AuthService {
	return &JWTService{
		secretKey: []byte(secret),
		issuer:    "go-order-system",
	}
}

// GenerateToken creates a new JWT for a given user.
func (s *JWTService) GenerateToken(ctx context.Context, user *domain.User) (string, error) {
	claims := &usecase.AuthClaims{
		UserID: int64(user.ID),
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
			Subject:   user.Name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken parses and validates a JWT string.
func (s *JWTService) ValidateToken(ctx context.Context, tokenString string) (*usecase.AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &usecase.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*usecase.AuthClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
