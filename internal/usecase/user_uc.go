package usecase

import (
	"context"
	"errors"

	"github.com/elokanugrah/backend-takehome/internal/domain"
)

type UserUseCase struct {
	userRepo    domain.UserRepository
	authService AuthService
}

func NewUserUseCase(ur domain.UserRepository, as AuthService) *UserUseCase {
	return &UserUseCase{
		userRepo:    ur,
		authService: as,
	}
}

func (uc *UserUseCase) Register(ctx context.Context, username, email, password string) (*domain.User, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	// Create new user domain object
	newUser, err := domain.NewUser(username, email, password)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := uc.userRepo.Save(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (uc *UserUseCase) Login(ctx context.Context, username, password string) (string, error) {
	// Find user by username
	user, err := uc.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", domain.ErrInvalidCredentials
	}

	// Check password
	if err := user.CheckPassword(password); err != nil {
		// Invalid password
		return "", domain.ErrInvalidCredentials
	}

	// Generate JWT
	token, err := uc.authService.GenerateToken(ctx, user)
	if err != nil {
		return "", err
	}

	return token, nil
}
