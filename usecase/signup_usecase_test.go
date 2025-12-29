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
)

func TestSignupUsecase_Signup(t *testing.T) {
	name := "Test User"
	email := "test@example.com"
	password := "password"

	expectedTokens := domain.TokenPair{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenService := new(MockTokenService)

		mockRepo.On("GetByEmail", mock.Anything, email).Return(domain.User{}, errors.New("not found"))
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
		mockTokenService.On("GenerateTokenPair", mock.Anything).Return(expectedTokens, nil)

		u := usecase.NewSignupUsecase(mockRepo, mockTokenService, time.Second*2)
		tokens, err := u.Signup(context.Background(), name, email, password)

		assert.NoError(t, err)
		assert.Equal(t, expectedTokens, tokens)
		mockRepo.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})

	t.Run("user_already_exists", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenService := new(MockTokenService)

		existingUser := domain.User{Email: email}
		mockRepo.On("GetByEmail", mock.Anything, email).Return(existingUser, nil)

		u := usecase.NewSignupUsecase(mockRepo, mockTokenService, time.Second*2)
		_, err := u.Signup(context.Background(), name, email, password)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenService := new(MockTokenService)

		mockRepo.On("GetByEmail", mock.Anything, email).Return(domain.User{}, errors.New("not found"))
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("database error"))

		u := usecase.NewSignupUsecase(mockRepo, mockTokenService, time.Second*2)
		_, err := u.Signup(context.Background(), name, email, password)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInternalServer)
		mockRepo.AssertExpectations(t)
	})
}
