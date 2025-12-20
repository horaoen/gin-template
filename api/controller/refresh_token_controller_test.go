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

// MockRefreshTokenUsecase
type MockRefreshTokenUsecase struct {
	mock.Mock
}

func (m *MockRefreshTokenUsecase) GetUserByID(c context.Context, id string) (domain.User, error) {
	args := m.Called(c, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockRefreshTokenUsecase) CreateAccessToken(user *domain.User, secret string, expiry int) (string, error) {
	args := m.Called(user, secret, expiry)
	return args.String(0), args.Error(1)
}

func (m *MockRefreshTokenUsecase) CreateRefreshToken(user *domain.User, secret string, expiry int) (string, error) {
	args := m.Called(user, secret, expiry)
	return args.String(0), args.Error(1)
}

func (m *MockRefreshTokenUsecase) ExtractIDFromToken(requestToken string, secret string) (string, error) {
	args := m.Called(requestToken, secret)
	return args.String(0), args.Error(1)
}

func TestRefreshTokenController_RefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	env := &bootstrap.Env{
		AccessTokenSecret:      "access_secret",
		AccessTokenExpiryHour:  2,
		RefreshTokenSecret:     "refresh_secret",
		RefreshTokenExpiryHour: 24,
	}

	user := domain.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	}

	t.Run("success", func(t *testing.T) {
		mockUsecase := new(MockRefreshTokenUsecase)
		rtc := controller.RefreshTokenController{
			RefreshTokenUsecase: mockUsecase,
			Env:                 env,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := url.Values{}
		data.Set("refreshToken", "valid_refresh_token")

		req, _ := http.NewRequest(http.MethodPost, "/refresh", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		// Mock expectations
		mockUsecase.On("ExtractIDFromToken", "valid_refresh_token", env.RefreshTokenSecret).Return("1", nil)
		mockUsecase.On("GetUserByID", mock.Anything, "1").Return(user, nil)
		mockUsecase.On("CreateAccessToken", mock.Anything, env.AccessTokenSecret, env.AccessTokenExpiryHour).Return("new_access_token", nil)
		mockUsecase.On("CreateRefreshToken", mock.Anything, env.RefreshTokenSecret, env.RefreshTokenExpiryHour).Return("new_refresh_token", nil)

		rtc.RefreshToken(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.RefreshTokenResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "new_access_token", response.AccessToken)
		assert.Equal(t, "new_refresh_token", response.RefreshToken)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("invalid_token", func(t *testing.T) {
		mockUsecase := new(MockRefreshTokenUsecase)
		rtc := controller.RefreshTokenController{
			RefreshTokenUsecase: mockUsecase,
			Env:                 env,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := url.Values{}
		data.Set("refreshToken", "invalid_token")

		req, _ := http.NewRequest(http.MethodPost, "/refresh", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		// Mock expectations
		mockUsecase.On("ExtractIDFromToken", "invalid_token", env.RefreshTokenSecret).Return("", errors.New("invalid token"))

		rtc.RefreshToken(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User not found", response.Message) // The controller returns "User not found" on extraction error

		mockUsecase.AssertExpectations(t)
	})

	t.Run("user_not_found", func(t *testing.T) {
		mockUsecase := new(MockRefreshTokenUsecase)
		rtc := controller.RefreshTokenController{
			RefreshTokenUsecase: mockUsecase,
			Env:                 env,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := url.Values{}
		data.Set("refreshToken", "valid_refresh_token")

		req, _ := http.NewRequest(http.MethodPost, "/refresh", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		// Mock expectations
		mockUsecase.On("ExtractIDFromToken", "valid_refresh_token", env.RefreshTokenSecret).Return("1", nil)
		mockUsecase.On("GetUserByID", mock.Anything, "1").Return(domain.User{}, errors.New("not found"))

		rtc.RefreshToken(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User not found", response.Message)

		mockUsecase.AssertExpectations(t)
	})
}
