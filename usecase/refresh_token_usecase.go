package usecase

import (
	"context"
	"time"

	"github.com/horaoen/go-backend-clean-architecture/domain"
)

type refreshTokenUsecase struct {
	userRepository domain.UserRepository
	tokenService   domain.TokenService
	contextTimeout time.Duration
}

func NewRefreshTokenUsecase(userRepository domain.UserRepository, tokenService domain.TokenService, timeout time.Duration) domain.RefreshTokenUsecase {
	return &refreshTokenUsecase{
		userRepository: userRepository,
		tokenService:   tokenService,
		contextTimeout: timeout,
	}
}

func (rtu *refreshTokenUsecase) Refresh(c context.Context, refreshToken string) (domain.TokenPair, error) {
	ctx, cancel := context.WithTimeout(c, rtu.contextTimeout)
	defer cancel()

	userID, err := rtu.tokenService.ExtractIDFromToken(refreshToken)
	if err != nil {
		return domain.TokenPair{}, domain.ErrInvalidToken
	}

	user, err := rtu.userRepository.GetByID(ctx, userID)
	if err != nil {
		return domain.TokenPair{}, domain.ErrUserNotFound
	}

	return rtu.tokenService.GenerateTokenPair(&user)
}
