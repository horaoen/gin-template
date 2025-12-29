package controller_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/horaoen/go-backend-clean-architecture/api/controller"
	"github.com/horaoen/go-backend-clean-architecture/api/dto"
	"github.com/horaoen/go-backend-clean-architecture/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRefreshTokenUsecase struct {
	mock.Mock
}

func (m *MockRefreshTokenUsecase) Refresh(c context.Context, refreshToken string) (domain.TokenPair, error) {
	args := m.Called(c, refreshToken)
	return args.Get(0).(domain.TokenPair), args.Error(1)
}

func TestRefreshTokenController_RefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	expectedTokens := domain.TokenPair{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
	}

	t.Run("success", func(t *testing.T) {
		mockUsecase := new(MockRefreshTokenUsecase)
		rtc := controller.RefreshTokenController{
			RefreshTokenUsecase: mockUsecase,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := url.Values{}
		data.Set("refreshToken", "valid_refresh_token")

		req, _ := http.NewRequest(http.MethodPost, "/refresh", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		mockUsecase.On("Refresh", mock.Anything, "valid_refresh_token").Return(expectedTokens, nil)

		rtc.RefreshToken(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.RefreshTokenResponse
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
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := url.Values{}
		data.Set("refreshToken", "invalid_token")

		req, _ := http.NewRequest(http.MethodPost, "/refresh", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		mockUsecase.On("Refresh", mock.Anything, "invalid_token").Return(domain.TokenPair{}, domain.ErrInvalidToken)

		rtc.RefreshToken(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid or expired token", response.Message)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("user_not_found", func(t *testing.T) {
		mockUsecase := new(MockRefreshTokenUsecase)
		rtc := controller.RefreshTokenController{
			RefreshTokenUsecase: mockUsecase,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := url.Values{}
		data.Set("refreshToken", "valid_token_for_missing_user")

		req, _ := http.NewRequest(http.MethodPost, "/refresh", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		mockUsecase.On("Refresh", mock.Anything, "valid_token_for_missing_user").Return(domain.TokenPair{}, domain.ErrUserNotFound)

		rtc.RefreshToken(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "user not found", response.Message)

		mockUsecase.AssertExpectations(t)
	})
}
