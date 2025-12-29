package domain

import "context"

type LoginUsecase interface {
	Login(c context.Context, email, password string) (TokenPair, error)
}
