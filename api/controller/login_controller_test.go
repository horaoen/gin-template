package controller_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/horaoen/go-backend-clean-architecture/api/controller"
	"github.com/horaoen/go-backend-clean-architecture/bootstrap"
	"github.com/horaoen/go-backend-clean-architecture/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockLoginUsecase
type MockLoginUsecase struct {
	mock.Mock
}

func (m *MockLoginUsecase) GetUserByEmail(c context.Context, email string) (domain.User, error) {
	args := m.Called(c, email)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockLoginUsecase) CreateAccessToken(user *domain.User, secret string, expiry int) (string, error) {
	args := m.Called(user, secret, expiry)
	return args.String(0), args.Error(1)
}

func (m *MockLoginUsecase) CreateRefreshToken(user *domain.User, secret string, expiry int) (string, error) {
	args := m.Called(user, secret, expiry)
	return args.String(0), args.Error(1)
}

func TestLoginController_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	env := &bootstrap.Env{
		AccessTokenSecret:      "access_secret",
		AccessTokenExpiryHour:  2,
		RefreshTokenSecret:     "refresh_secret",
		RefreshTokenExpiryHour: 24,
	}

	password := "password"
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := domain.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: string(encryptedPassword),
	}

	t.Run("success", func(t *testing.T) {
		mockUsecase := new(MockLoginUsecase)
		lc := controller.LoginController{
			LoginUsecase: mockUsecase,
			Env:          env,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := url.Values{}
		data.Set("email", "test@example.com")
		data.Set("password", password)

		req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		// Mock expectations
		mockUsecase.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)
		mockUsecase.On("CreateAccessToken", mock.Anything, env.AccessTokenSecret, env.AccessTokenExpiryHour).Return("access_token", nil)
		mockUsecase.On("CreateRefreshToken", mock.Anything, env.RefreshTokenSecret, env.RefreshTokenExpiryHour).Return("refresh_token", nil)

		lc.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "access_token", response.AccessToken)
		assert.Equal(t, "refresh_token", response.RefreshToken)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("user_not_found", func(t *testing.T) {
		mockUsecase := new(MockLoginUsecase)
		lc := controller.LoginController{
			LoginUsecase: mockUsecase,
			Env:          env,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := url.Values{}
		data.Set("email", "notfound@example.com")
		data.Set("password", "password")

		req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		// Mock expectations
		mockUsecase.On("GetUserByEmail", mock.Anything, "notfound@example.com").Return(domain.User{}, errors.New("not found"))

		lc.Login(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User not found with the given email", response.Message)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("invalid_credentials", func(t *testing.T) {
		mockUsecase := new(MockLoginUsecase)
		lc := controller.LoginController{
			LoginUsecase: mockUsecase,
			Env:          env,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := url.Values{}
		data.Set("email", "test@example.com")
		data.Set("password", "wrong_password")

		req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		// Mock expectations
		mockUsecase.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)

		lc.Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid credentials", response.Message)

		mockUsecase.AssertExpectations(t)
	})
}
