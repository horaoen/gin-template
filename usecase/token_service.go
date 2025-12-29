package usecase

import (
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/horaoen/go-backend-clean-architecture/domain"
)

type tokenService struct {
	accessTokenSecret      string
	refreshTokenSecret     string
	accessTokenExpiryHour  int
	refreshTokenExpiryHour int
}

func NewTokenService(
	accessTokenSecret string,
	refreshTokenSecret string,
	accessTokenExpiryHour int,
	refreshTokenExpiryHour int,
) domain.TokenService {
	return &tokenService{
		accessTokenSecret:      accessTokenSecret,
		refreshTokenSecret:     refreshTokenSecret,
		accessTokenExpiryHour:  accessTokenExpiryHour,
		refreshTokenExpiryHour: refreshTokenExpiryHour,
	}
}

func (ts *tokenService) GenerateTokenPair(user *domain.User) (domain.TokenPair, error) {
	accessToken, err := ts.createAccessToken(user)
	if err != nil {
		return domain.TokenPair{}, err
	}

	refreshToken, err := ts.createRefreshToken(user)
	if err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (ts *tokenService) ExtractIDFromToken(requestToken string) (string, error) {
	claims := &domain.JwtCustomRefreshClaims{}
	token, err := jwt.ParseWithClaims(requestToken, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return []byte(ts.refreshTokenSecret), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", domain.ErrInvalidToken
	}
	return claims.ID, nil
}

func (ts *tokenService) createAccessToken(user *domain.User) (string, error) {
	exp := time.Now().Add(time.Hour * time.Duration(ts.accessTokenExpiryHour))
	claims := &domain.JwtCustomClaims{
		Name: user.Name,
		ID:   strconv.FormatUint(uint64(user.ID), 10),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(ts.accessTokenSecret))
}

func (ts *tokenService) createRefreshToken(user *domain.User) (string, error) {
	exp := time.Now().Add(time.Hour * time.Duration(ts.refreshTokenExpiryHour))
	claims := &domain.JwtCustomRefreshClaims{
		ID: strconv.FormatUint(uint64(user.ID), 10),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(ts.refreshTokenSecret))
}
