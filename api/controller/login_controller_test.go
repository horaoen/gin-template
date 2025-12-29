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

type MockLoginUsecase struct {
	mock.Mock
}

func (m *MockLoginUsecase) Login(c context.Context, email, password string) (domain.TokenPair, error) {
	args := m.Called(c, email, password)
	return args.Get(0).(domain.TokenPair), args.Error(1)
}

func TestLoginController_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	expectedTokens := domain.TokenPair{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}

	t.Run("success", func(t *testing.T) {
		mockUsecase := new(MockLoginUsecase)
		lc := controller.LoginController{
			LoginUsecase: mockUsecase,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := url.Values{}
		data.Set("email", "test@example.com")
		data.Set("password", "password")

		req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		mockUsecase.On("Login", mock.Anything, "test@example.com", "password").Return(expectedTokens, nil)

		lc.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "access_token", response.AccessToken)
		assert.Equal(t, "refresh_token", response.RefreshToken)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("invalid_credentials", func(t *testing.T) {
		mockUsecase := new(MockLoginUsecase)
		lc := controller.LoginController{
			LoginUsecase: mockUsecase,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := url.Values{}
		data.Set("email", "test@example.com")
		data.Set("password", "wrong_password")

		req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		mockUsecase.On("Login", mock.Anything, "test@example.com", "wrong_password").Return(domain.TokenPair{}, domain.ErrInvalidCredentials)

		lc.Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid email or password", response.Message)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("bad_request", func(t *testing.T) {
		mockUsecase := new(MockLoginUsecase)
		lc := controller.LoginController{
			LoginUsecase: mockUsecase,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := url.Values{}
		data.Set("email", "invalid-email")
		data.Set("password", "password")

		req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		lc.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
