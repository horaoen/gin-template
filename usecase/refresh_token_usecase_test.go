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

func TestRefreshTokenUsecase_Refresh(t *testing.T) {
	userID := "1"
	refreshToken := "valid_refresh_token"
	user := domain.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	}

	expectedTokens := domain.TokenPair{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenService := new(MockTokenService)

		mockTokenService.On("ExtractIDFromToken", refreshToken).Return(userID, nil)
		mockRepo.On("GetByID", mock.Anything, userID).Return(user, nil)
		mockTokenService.On("GenerateTokenPair", mock.Anything).Return(expectedTokens, nil)

		u := usecase.NewRefreshTokenUsecase(mockRepo, mockTokenService, time.Second*2)
		tokens, err := u.Refresh(context.Background(), refreshToken)

		assert.NoError(t, err)
		assert.Equal(t, expectedTokens, tokens)
		mockRepo.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})

	t.Run("invalid_token", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenService := new(MockTokenService)

		mockTokenService.On("ExtractIDFromToken", "invalid_token").Return("", errors.New("invalid token"))

		u := usecase.NewRefreshTokenUsecase(mockRepo, mockTokenService, time.Second*2)
		_, err := u.Refresh(context.Background(), "invalid_token")

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidToken)
		mockTokenService.AssertExpectations(t)
	})

	t.Run("user_not_found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenService := new(MockTokenService)

		mockTokenService.On("ExtractIDFromToken", refreshToken).Return(userID, nil)
		mockRepo.On("GetByID", mock.Anything, userID).Return(domain.User{}, errors.New("not found"))

		u := usecase.NewRefreshTokenUsecase(mockRepo, mockTokenService, time.Second*2)
		_, err := u.Refresh(context.Background(), refreshToken)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		mockRepo.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})
}
