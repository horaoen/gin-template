package domain

import "context"

type RefreshTokenUsecase interface {
	Refresh(c context.Context, refreshToken string) (TokenPair, error)
}
