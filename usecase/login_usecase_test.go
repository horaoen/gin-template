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

func TestLoginUsecase_GetUserByEmail(t *testing.T) {
	email := "test@example.com"
	user := domain.User{
		Name:     "Test User",
		Email:    email,
		Password: "password",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)

		u := usecase.NewLoginUsecase(mockRepo, time.Second*2)
		fetchedUser, err := u.GetUserByEmail(context.Background(), email)

		assert.NoError(t, err)
		assert.Equal(t, user, fetchedUser)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("GetByEmail", mock.Anything, email).Return(domain.User{}, errors.New("not found"))

		u := usecase.NewLoginUsecase(mockRepo, time.Second*2)
		_, err := u.GetUserByEmail(context.Background(), email)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
