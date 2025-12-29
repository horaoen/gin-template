package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/horaoen/go-backend-clean-architecture/domain"
	"github.com/horaoen/go-backend-clean-architecture/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateTokenPair(user *domain.User) (domain.TokenPair, error) {
	args := m.Called(user)
	return args.Get(0).(domain.TokenPair), args.Error(1)
}

func (m *MockTokenService) ExtractIDFromToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func TestLoginUsecase_Login(t *testing.T) {
	email := "test@example.com"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := domain.User{
		ID:       1,
		Name:     "Test User",
		Email:    email,
		Password: string(hashedPassword),
	}

	expectedTokens := domain.TokenPair{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenService := new(MockTokenService)

		mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)
		mockTokenService.On("GenerateTokenPair", mock.Anything).Return(expectedTokens, nil)

		u := usecase.NewLoginUsecase(mockRepo, mockTokenService, time.Second*2)
		tokens, err := u.Login(context.Background(), email, password)

		assert.NoError(t, err)
		assert.Equal(t, expectedTokens, tokens)
		mockRepo.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})

	t.Run("user_not_found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenService := new(MockTokenService)

		mockRepo.On("GetByEmail", mock.Anything, email).Return(domain.User{}, errors.New("not found"))

		u := usecase.NewLoginUsecase(mockRepo, mockTokenService, time.Second*2)
		_, err := u.Login(context.Background(), email, password)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid_password", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenService := new(MockTokenService)

		mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)

		u := usecase.NewLoginUsecase(mockRepo, mockTokenService, time.Second*2)
		_, err := u.Login(context.Background(), email, "wrong_password")

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
		mockRepo.AssertExpectations(t)
	})
}
