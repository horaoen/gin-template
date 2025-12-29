package domain

import "context"

type SignupUsecase interface {
	Signup(c context.Context, name, email, password string) (TokenPair, error)
}
