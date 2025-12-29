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

type MockSignupUsecase struct {
	mock.Mock
}

func (m *MockSignupUsecase) Signup(c context.Context, name, email, password string) (domain.TokenPair, error) {
	args := m.Called(c, name, email, password)
	return args.Get(0).(domain.TokenPair), args.Error(1)
}

func TestSignupController_Signup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	expectedTokens := domain.TokenPair{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}

	t.Run("success", func(t *testing.T) {
		mockUsecase := new(MockSignupUsecase)
		sc := controller.SignupController{
			SignupUsecase: mockUsecase,
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

		mockUsecase.On("Signup", mock.Anything, "Test User", "test@example.com", "password").Return(expectedTokens, nil)

		sc.Signup(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.SignupResponse
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

		mockUsecase.On("Signup", mock.Anything, "Test User", "existing@example.com", "password").Return(domain.TokenPair{}, domain.ErrUserAlreadyExists)

		sc.Signup(c)

		assert.Equal(t, http.StatusConflict, w.Code)

		var response domain.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "user already exists with the given email", response.Message)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("bad_request", func(t *testing.T) {
		mockUsecase := new(MockSignupUsecase)
		sc := controller.SignupController{
			SignupUsecase: mockUsecase,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := url.Values{}
		data.Set("name", "Test User")
		data.Set("password", "password")

		req, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		sc.Signup(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
