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

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(c context.Context, user *domain.User) error {
	args := m.Called(c, user)
	return args.Error(0)
}

func (m *MockUserRepository) Fetch(c context.Context) ([]domain.User, error) {
	args := m.Called(c)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(c context.Context, email string) (domain.User, error) {
	args := m.Called(c, email)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(c context.Context, id string) (domain.User, error) {
	args := m.Called(c, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func TestSignupUsecase_Create(t *testing.T) {
	user := &domain.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("Create", mock.Anything, user).Return(nil)

		u := usecase.NewSignupUsecase(mockRepo, time.Second*2)
		err := u.Create(context.Background(), user)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("Create", mock.Anything, user).Return(errors.New("database error"))

		u := usecase.NewSignupUsecase(mockRepo, time.Second*2)
		err := u.Create(context.Background(), user)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSignupUsecase_GetUserByEmail(t *testing.T) {
	email := "test@example.com"
	user := domain.User{
		Name:     "Test User",
		Email:    email,
		Password: "password",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)

		u := usecase.NewSignupUsecase(mockRepo, time.Second*2)
		fetchedUser, err := u.GetUserByEmail(context.Background(), email)

		assert.NoError(t, err)
		assert.Equal(t, user, fetchedUser)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("GetByEmail", mock.Anything, email).Return(domain.User{}, errors.New("not found"))

		u := usecase.NewSignupUsecase(mockRepo, time.Second*2)
		_, err := u.GetUserByEmail(context.Background(), email)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
