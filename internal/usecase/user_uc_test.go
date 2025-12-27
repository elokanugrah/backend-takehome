package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/elokanugrah/backend-takehome/internal/domain"
	domainmocks "github.com/elokanugrah/backend-takehome/internal/domain/mocks"
	"github.com/elokanugrah/backend-takehome/internal/usecase"
	usecasemocks "github.com/elokanugrah/backend-takehome/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	ctx := context.Background()
	name := "John Doe"
	email := "john@example.com"
	password := "password123"

	t.Run("success", func(t *testing.T) {
		// Setup Mocks
		mockUserRepo := new(domainmocks.UserRepository)
		mockAuthService := new(usecasemocks.AuthService)
		mockTxManager := new(usecasemocks.TransactionManager)

		uc := usecase.NewUserUseCase(mockUserRepo, mockAuthService, mockTxManager)

		// Mock TransactionManager to execute the callback function
		mockTxManager.On("WithTransaction", ctx, mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		// Mock Repository calls inside transaction
		mockUserRepo.On("FindByEmail", ctx, email).Return(nil, nil)
		mockUserRepo.On("Save", ctx, mock.MatchedBy(func(u *domain.User) bool {
			return u.Name == name && u.Email == email
		})).Return(nil)

		user, err := uc.Register(ctx, name, email, password)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		mockUserRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("email already taken", func(t *testing.T) {
		// Setup Mocks
		mockUserRepo := new(domainmocks.UserRepository)
		mockAuthService := new(usecasemocks.AuthService)
		mockTxManager := new(usecasemocks.TransactionManager)

		uc := usecase.NewUserUseCase(mockUserRepo, mockAuthService, mockTxManager)

		existingUser := &domain.User{ID: 1, Email: email}

		mockTxManager.On("WithTransaction", ctx, mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		mockUserRepo.On("FindByEmail", ctx, email).Return(existingUser, nil)

		user, err := uc.Register(ctx, name, email, password)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "email already taken", err.Error())
		mockUserRepo.AssertNotCalled(t, "Save", ctx, mock.Anything)
	})

	t.Run("error saving user", func(t *testing.T) {
		// Setup Mocks
		mockUserRepo := new(domainmocks.UserRepository)
		mockAuthService := new(usecasemocks.AuthService)
		mockTxManager := new(usecasemocks.TransactionManager)

		uc := usecase.NewUserUseCase(mockUserRepo, mockAuthService, mockTxManager)

		mockTxManager.On("WithTransaction", ctx, mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		mockUserRepo.On("FindByEmail", ctx, email).Return(nil, nil)
		mockUserRepo.On("Save", ctx, mock.Anything).Return(errors.New("db error"))

		user, err := uc.Register(ctx, name, email, password)

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	email := "john@example.com"
	password := "password123"

	// Create a user with hashed password for testing
	user, _ := domain.NewUser("John Doe", email, password)

	t.Run("success", func(t *testing.T) {
		// Setup Mocks
		mockUserRepo := new(domainmocks.UserRepository)
		mockAuthService := new(usecasemocks.AuthService)
		mockTxManager := new(usecasemocks.TransactionManager)

		uc := usecase.NewUserUseCase(mockUserRepo, mockAuthService, mockTxManager)

		mockUserRepo.On("FindByEmail", ctx, email).Return(user, nil)
		mockAuthService.On("GenerateToken", ctx, user).Return("mock-token", nil)

		token, err := uc.Login(ctx, email, password)

		assert.NoError(t, err)
		assert.Equal(t, "mock-token", token)
	})

	t.Run("user not found", func(t *testing.T) {
		// Setup Mocks
		mockUserRepo := new(domainmocks.UserRepository)
		mockAuthService := new(usecasemocks.AuthService)
		mockTxManager := new(usecasemocks.TransactionManager)

		uc := usecase.NewUserUseCase(mockUserRepo, mockAuthService, mockTxManager)

		mockUserRepo.On("FindByEmail", ctx, email).Return(nil, nil)

		token, err := uc.Login(ctx, email, password)

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, domain.ErrInvalidCredentials, err)
	})

	t.Run("wrong password", func(t *testing.T) {
		// Setup Mocks
		mockUserRepo := new(domainmocks.UserRepository)
		mockAuthService := new(usecasemocks.AuthService)
		mockTxManager := new(usecasemocks.TransactionManager)

		uc := usecase.NewUserUseCase(mockUserRepo, mockAuthService, mockTxManager)

		mockUserRepo.On("FindByEmail", ctx, email).Return(user, nil)

		token, err := uc.Login(ctx, email, "wrongpass")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, domain.ErrInvalidCredentials, err)
	})

	t.Run("token generation error", func(t *testing.T) {
		// Setup Mocks
		mockUserRepo := new(domainmocks.UserRepository)
		mockAuthService := new(usecasemocks.AuthService)
		mockTxManager := new(usecasemocks.TransactionManager)

		uc := usecase.NewUserUseCase(mockUserRepo, mockAuthService, mockTxManager)

		mockUserRepo.On("FindByEmail", ctx, email).Return(user, nil)
		mockAuthService.On("GenerateToken", ctx, user).Return("", errors.New("token error"))

		token, err := uc.Login(ctx, email, password)

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}
