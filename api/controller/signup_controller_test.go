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
)

// MockSignupUsecase
type MockSignupUsecase struct {
	mock.Mock
}

func (m *MockSignupUsecase) Create(c context.Context, user *domain.User) error {
	args := m.Called(c, user)
	return args.Error(0)
}

func (m *MockSignupUsecase) GetUserByEmail(c context.Context, email string) (domain.User, error) {
	args := m.Called(c, email)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockSignupUsecase) CreateAccessToken(user *domain.User, secret string, expiry int) (string, error) {
	args := m.Called(user, secret, expiry)
	return args.String(0), args.Error(1)
}

func (m *MockSignupUsecase) CreateRefreshToken(user *domain.User, secret string, expiry int) (string, error) {
	args := m.Called(user, secret, expiry)
	return args.String(0), args.Error(1)
}

func TestSignupController_Signup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	env := &bootstrap.Env{
		AccessTokenSecret:      "access_secret",
		AccessTokenExpiryHour:  2,
		RefreshTokenSecret:     "refresh_secret",
		RefreshTokenExpiryHour: 24,
	}

	t.Run("success", func(t *testing.T) {
		mockUsecase := new(MockSignupUsecase)
		sc := controller.SignupController{
			SignupUsecase: mockUsecase,
			Env:           env,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := url.Values{}
		data.Set("name", "Test User")
		data.Set("email", "test@example.com")
		data.Set("password", "password")

		req, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		// Mock expectations
		mockUsecase.On("GetUserByEmail", mock.Anything, "test@example.com").Return(domain.User{}, errors.New("not found"))
		mockUsecase.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
		mockUsecase.On("CreateAccessToken", mock.Anything, env.AccessTokenSecret, env.AccessTokenExpiryHour).Return("access_token", nil)
		mockUsecase.On("CreateRefreshToken", mock.Anything, env.RefreshTokenSecret, env.RefreshTokenExpiryHour).Return("refresh_token", nil)

		sc.Signup(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.SignupResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "access_token", response.AccessToken)
		assert.Equal(t, "refresh_token", response.RefreshToken)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("user_already_exists", func(t *testing.T) {
		mockUsecase := new(MockSignupUsecase)
		sc := controller.SignupController{
			SignupUsecase: mockUsecase,
			Env:           env,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := url.Values{}
		data.Set("name", "Test User")
		data.Set("email", "existing@example.com")
		data.Set("password", "password")

		req, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		// Mock expectations
		mockUsecase.On("GetUserByEmail", mock.Anything, "existing@example.com").Return(domain.User{Email: "existing@example.com"}, nil)

		sc.Signup(c)

		assert.Equal(t, http.StatusConflict, w.Code)

		var response domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User already exists with the given email", response.Message)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("bad_request", func(t *testing.T) {
		mockUsecase := new(MockSignupUsecase)
		sc := controller.SignupController{
			SignupUsecase: mockUsecase,
			Env:           env,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Missing email
		data := url.Values{}
		data.Set("name", "Test User")
		// data.Set("email", "test@example.com")
		data.Set("password", "password")

		req, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		sc.Signup(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		// No mock expectations needed as it fails at binding
	})
}
