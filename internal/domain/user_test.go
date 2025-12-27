package domain_test

import (
	"testing"

	"github.com/elokanugrah/backend-takehome/internal/domain"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	t.Run("should create user successfully with valid data", func(t *testing.T) {
		name := "John Doe"
		email := "john@example.com"
		password := "secret123"

		user, err := domain.NewUser(name, email, password)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, name, user.Name)
		assert.Equal(t, email, user.Email)
		assert.NotEmpty(t, user.PasswordHash)
		assert.NotEqual(t, password, user.PasswordHash) // Password should be hashed
	})

	t.Run("should return error for invalid email", func(t *testing.T) {
		name := "John Doe"
		email := "invalid-email"
		password := "secret123"

		user, err := domain.NewUser(name, email, password)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, domain.ErrInvalidEmail, err)
	})
}

func TestCheckPassword(t *testing.T) {
	password := "secret123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &domain.User{
		PasswordHash: string(hashed),
	}

	t.Run("should return nil for correct password", func(t *testing.T) {
		err := user.CheckPassword(password)
		assert.NoError(t, err)
	})

	t.Run("should return error for incorrect password", func(t *testing.T) {
		err := user.CheckPassword("wrongpassword")
		assert.Error(t, err)
	})
}
