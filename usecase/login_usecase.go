package usecase

import (
	"context"
	"time"

	"github.com/horaoen/go-backend-clean-architecture/domain"
	"golang.org/x/crypto/bcrypt"
)

type loginUsecase struct {
	userRepository domain.UserRepository
	tokenService   domain.TokenService
	contextTimeout time.Duration
}

func NewLoginUsecase(userRepository domain.UserRepository, tokenService domain.TokenService, timeout time.Duration) domain.LoginUsecase {
	return &loginUsecase{
		userRepository: userRepository,
		tokenService:   tokenService,
		contextTimeout: timeout,
	}
}

func (lu *loginUsecase) Login(c context.Context, email, password string) (domain.TokenPair, error) {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()

	user, err := lu.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return domain.TokenPair{}, domain.ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return domain.TokenPair{}, domain.ErrInvalidCredentials
	}

	return lu.tokenService.GenerateTokenPair(&user)
}
