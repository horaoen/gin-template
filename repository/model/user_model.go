package model

import (
	"time"

	"github.com/horaoen/go-backend-clean-architecture/domain"
)

type UserModel struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null"`
	Email     string `gorm:"size:255;uniqueIndex;not null"`
	Password  string `gorm:"column:password;size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserModel) TableName() string {
	return "users"
}

func (m *UserModel) ToDomain() domain.User {
	return domain.User{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToUserModel(u *domain.User) UserModel {
	return UserModel{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
