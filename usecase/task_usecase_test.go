package usecase_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/horaoen/go-backend-clean-architecture/domain"
	"github.com/horaoen/go-backend-clean-architecture/domain/mocks"
	"github.com/horaoen/go-backend-clean-architecture/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetchByUserID(t *testing.T) {
	mockTaskRepository := new(mocks.TaskRepository)
	var userIDUint uint = 1
	userID := strconv.FormatUint(uint64(userIDUint), 10)

	t.Run("success", func(t *testing.T) {

		mockTask := domain.Task{
			ID:     1,
			Title:  "Test Title",
			UserID: userIDUint,
		}

		mockListTask := make([]domain.Task, 0)
		mockListTask = append(mockListTask, mockTask)

		mockTaskRepository.On("FetchByUserID", mock.Anything, userID).Return(mockListTask, nil).Once()

		u := usecase.NewTaskUsecase(mockTaskRepository, time.Second*2)

		list, err := u.FetchByUserID(context.Background(), userID)

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.Len(t, list, len(mockListTask))

		mockTaskRepository.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockTaskRepository.On("FetchByUserID", mock.Anything, userID).Return(nil, errors.New("Unexpected")).Once()

		u := usecase.NewTaskUsecase(mockTaskRepository, time.Second*2)

		list, err := u.FetchByUserID(context.Background(), userID)

		assert.Error(t, err)
		assert.Nil(t, list)

		mockTaskRepository.AssertExpectations(t)
	})

}
