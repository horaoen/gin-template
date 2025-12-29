package repository

import (
	"context"
	"strconv"

	"github.com/horaoen/go-backend-clean-architecture/domain"
	"github.com/horaoen/go-backend-clean-architecture/repository/model"
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
	userModel := model.ToUserModel(user)
	if err := ur.db.WithContext(c).Create(&userModel).Error; err != nil {
		return err
	}
	user.ID = userModel.ID
	user.CreatedAt = userModel.CreatedAt
	user.UpdatedAt = userModel.UpdatedAt
	return nil
}

func (ur *userRepository) Fetch(c context.Context) ([]domain.User, error) {
	var userModels []model.UserModel
	err := ur.db.WithContext(c).Select("id", "name", "email", "created_at", "updated_at").Find(&userModels).Error
	if err != nil {
		return nil, err
	}

	if userModels == nil {
		return []domain.User{}, nil
	}

	users := make([]domain.User, len(userModels))
	for i, m := range userModels {
		users[i] = m.ToDomain()
	}

	return users, nil
}

func (ur *userRepository) GetByEmail(c context.Context, email string) (domain.User, error) {
	var userModel model.UserModel
	err := ur.db.WithContext(c).Where("email = ?", email).First(&userModel).Error
	if err != nil {
		return domain.User{}, err
	}
	return userModel.ToDomain(), nil
}

func (ur *userRepository) GetByID(c context.Context, id string) (domain.User, error) {
	var userModel model.UserModel

	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return domain.User{}, err
	}

	err = ur.db.WithContext(c).First(&userModel, userID).Error
	if err != nil {
		return domain.User{}, err
	}
	return userModel.ToDomain(), nil
}

func (ur *userRepository) Update(c context.Context, user *domain.User) error {
	userModel := model.ToUserModel(user)
	if err := ur.db.WithContext(c).Save(&userModel).Error; err != nil {
		return err
	}
	user.UpdatedAt = userModel.UpdatedAt
	return nil
}
