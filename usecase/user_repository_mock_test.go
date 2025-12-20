package usecase_test

import (
	"context"

	"github.com/horaoen/go-backend-clean-architecture/domain"
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

func (m *MockUserRepository) Update(c context.Context, user *domain.User) error {
	args := m.Called(c, user)
	return args.Error(0)
}
