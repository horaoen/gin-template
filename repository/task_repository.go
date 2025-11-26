package repository

import (
	"context"
	"strconv"

	"github.com/horaoen/go-backend-clean-architecture/domain"
	"gorm.io/gorm"
)

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) domain.TaskRepository {
	return &taskRepository{
		db: db,
	}
}

func (tr *taskRepository) Create(c context.Context, task *domain.Task) error {
	return tr.db.WithContext(c).Create(task).Error
}

func (tr *taskRepository) FetchByUserID(c context.Context, userID string) ([]domain.Task, error) {
	var tasks []domain.Task

	uid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return tasks, err
	}

	err = tr.db.WithContext(c).Where("user_id = ?", uid).Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	if tasks == nil {
		return []domain.Task{}, nil
	}

	return tasks, nil
}
