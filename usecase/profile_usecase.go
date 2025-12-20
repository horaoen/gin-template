package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/horaoen/go-backend-clean-architecture/domain"
	"golang.org/x/crypto/bcrypt"
)

type profileUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewProfileUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.ProfileUsecase {
	return &profileUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (pu *profileUsecase) GetProfileByID(c context.Context, userID string) (*domain.Profile, error) {
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	user, err := pu.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.Profile{Name: user.Name, Email: user.Email}, nil
}

func (pu *profileUsecase) ChangePassword(c context.Context, userID string, oldPassword string, newPassword string) error {
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	user, err := pu.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)) != nil {
		return errors.New("invalid old password")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(encryptedPassword)
	return pu.userRepository.Update(ctx, &user)
}
