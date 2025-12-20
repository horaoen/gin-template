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

func TestRefreshTokenUsecase_GetUserByID(t *testing.T) {
	id := "1"
	user := domain.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("GetByID", mock.Anything, id).Return(user, nil)

		u := usecase.NewRefreshTokenUsecase(mockRepo, time.Second*2)
		fetchedUser, err := u.GetUserByID(context.Background(), id)

		assert.NoError(t, err)
		assert.Equal(t, user, fetchedUser)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("GetByID", mock.Anything, id).Return(domain.User{}, errors.New("not found"))

		u := usecase.NewRefreshTokenUsecase(mockRepo, time.Second*2)
		_, err := u.GetUserByID(context.Background(), id)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
