package domain

import (
	"context"

	"gorm.io/gorm"
)

const (
	CollectionTask = "tasks"
)

type Task struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Title  string `gorm:"size:255;not null" form:"title" binding:"required" json:"title"`
	UserID uint   `gorm:"not null;index" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID" json:"-"`
	gorm.Model
}

type TaskRepository interface {
	Create(c context.Context, task *Task) error
	FetchByUserID(c context.Context, userID string) ([]Task, error)
}

type TaskUsecase interface {
	Create(c context.Context, task *Task) error
	FetchByUserID(c context.Context, userID string) ([]Task, error)
}
