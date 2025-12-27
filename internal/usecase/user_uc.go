package usecase

import (
	"context"
	"errors"

	"github.com/elokanugrah/backend-takehome/internal/domain"
)

type UserUseCase struct {
	userRepo    domain.UserRepository
	authService AuthService
	txManager   TransactionManager
}

func NewUserUseCase(ur domain.UserRepository, as AuthService, tm TransactionManager) *UserUseCase {
	return &UserUseCase{
		userRepo:    ur,
		authService: as,
		txManager:   tm,
	}
}

func (uc *UserUseCase) Register(ctx context.Context, name, email, password string) (*domain.User, error) {
	var registeredUser *domain.User

	err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		existingUser, err := uc.userRepo.FindByEmail(txCtx, email)
		if err != nil {
			return err
		}
		if existingUser != nil {
			return errors.New("email already taken")
		}

		newUser, err := domain.NewUser(name, email, password)
		if err != nil {
			return err
		}

		if err := uc.userRepo.Save(txCtx, newUser); err != nil {
			return err
		}

		registeredUser = newUser
		return nil
	})

	if err != nil {
		return nil, err
	}

	return registeredUser, nil
}

func (uc *UserUseCase) Login(ctx context.Context, email, password string) (string, error) {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, email)
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
