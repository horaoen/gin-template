package repository

import (
	"context"
	"strconv"

	"github.com/horaoen/go-backend-clean-architecture/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) Create(c context.Context, user *domain.User) error {
	return ur.db.WithContext(c).Create(user).Error
}

func (ur *userRepository) Fetch(c context.Context) ([]domain.User, error) {
	var users []domain.User
	err := ur.db.WithContext(c).Select("id", "name", "email", "created_at", "updated_at").Find(&users).Error
	if err != nil {
		return nil, err
	}

	if users == nil {
		return []domain.User{}, nil
	}

	return users, nil
}

func (ur *userRepository) GetByEmail(c context.Context, email string) (domain.User, error) {
	var user domain.User
	err := ur.db.WithContext(c).Where("email = ?", email).First(&user).Error
	return user, err
}

func (ur *userRepository) GetByID(c context.Context, id string) (domain.User, error) {
	var user domain.User

	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return user, err
	}

	err = ur.db.WithContext(c).First(&user, userID).Error
	return user, err
}
