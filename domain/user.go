package domain

import (
	"context"

	"gorm.io/gorm"
)

const (
	CollectionUser = "users"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"size:255;not null" json:"name"`
	Email    string `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	gorm.Model
}

type UserRepository interface {
	Create(c context.Context, user *User) error
	Fetch(c context.Context) ([]User, error)
	GetByEmail(c context.Context, email string) (User, error)
	GetByID(c context.Context, id string) (User, error)
	Update(c context.Context, user *User) error
}
