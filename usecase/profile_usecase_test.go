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

func TestProfileUsecase_ChangePassword(t *testing.T) {
	userID := "1"
	oldPassword := "old_password"
	newPassword := "new_password"
	hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)
	user := domain.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hashedOldPassword),
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("GetByID", mock.Anything, userID).Return(user, nil)
		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(newPassword))
			return u.ID == user.ID && err == nil
		})).Return(nil)

		pu := usecase.NewProfileUsecase(mockRepo, time.Second*2)
		err := pu.ChangePassword(context.Background(), userID, oldPassword, newPassword)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid_old_password", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("GetByID", mock.Anything, userID).Return(user, nil)

		pu := usecase.NewProfileUsecase(mockRepo, time.Second*2)
		err := pu.ChangePassword(context.Background(), userID, "wrong_old_password", newPassword)

		assert.Error(t, err)
		assert.Equal(t, "invalid old password", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("user_not_found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("GetByID", mock.Anything, userID).Return(domain.User{}, errors.New("user not found"))

		pu := usecase.NewProfileUsecase(mockRepo, time.Second*2)
		err := pu.ChangePassword(context.Background(), userID, oldPassword, newPassword)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
