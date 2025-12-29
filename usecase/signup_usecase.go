package usecase

import (
	"context"
	"time"

	"github.com/horaoen/go-backend-clean-architecture/domain"
	"golang.org/x/crypto/bcrypt"
)

type signupUsecase struct {
	userRepository domain.UserRepository
	tokenService   domain.TokenService
	contextTimeout time.Duration
}

func NewSignupUsecase(userRepository domain.UserRepository, tokenService domain.TokenService, timeout time.Duration) domain.SignupUsecase {
	return &signupUsecase{
		userRepository: userRepository,
		tokenService:   tokenService,
		contextTimeout: timeout,
	}
}

func (su *signupUsecase) Signup(c context.Context, name, email, password string) (domain.TokenPair, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	_, err := su.userRepository.GetByEmail(ctx, email)
	if err == nil {
		return domain.TokenPair{}, domain.ErrUserAlreadyExists
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.TokenPair{}, domain.ErrInternalServer
	}

	user := domain.User{
		Name:     name,
		Email:    email,
		Password: string(encryptedPassword),
	}

	if err := su.userRepository.Create(ctx, &user); err != nil {
		return domain.TokenPair{}, domain.ErrInternalServer
	}

	return su.tokenService.GenerateTokenPair(&user)
}
