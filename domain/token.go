package domain

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenService interface {
	GenerateTokenPair(user *User) (TokenPair, error)
	ExtractIDFromToken(token string) (string, error)
}
