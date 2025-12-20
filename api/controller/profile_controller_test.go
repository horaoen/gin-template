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
	"github.com/horaoen/go-backend-clean-architecture/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProfileUsecase
type MockProfileUsecase struct {
	mock.Mock
}

func (m *MockProfileUsecase) GetProfileByID(c context.Context, userID string) (*domain.Profile, error) {
	args := m.Called(c, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Profile), args.Error(1)
}

func (m *MockProfileUsecase) ChangePassword(c context.Context, userID string, oldPassword string, newPassword string) error {
	args := m.Called(c, userID, oldPassword, newPassword)
	return args.Error(0)
}

func TestProfileController_ChangePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := "1"
	oldPassword := "old_password"
	newPassword := "new_password123"

	t.Run("success", func(t *testing.T) {
		mockUsecase := new(MockProfileUsecase)
		pc := controller.ProfileController{
			ProfileUsecase: mockUsecase,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("x-user-id", userID)

		data := url.Values{}
		data.Set("oldPassword", oldPassword)
		data.Set("newPassword", newPassword)

		req, _ := http.NewRequest(http.MethodPost, "/profile/change-password", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		mockUsecase.On("ChangePassword", mock.Anything, userID, oldPassword, newPassword).Return(nil)

		pc.ChangePassword(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response domain.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Password changed successfully", response.Message)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("invalid_old_password", func(t *testing.T) {
		mockUsecase := new(MockProfileUsecase)
		pc := controller.ProfileController{
			ProfileUsecase: mockUsecase,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("x-user-id", userID)

		data := url.Values{}
		data.Set("oldPassword", "wrong")
		data.Set("newPassword", newPassword)

		req, _ := http.NewRequest(http.MethodPost, "/profile/change-password", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request = req

		mockUsecase.On("ChangePassword", mock.Anything, userID, "wrong", newPassword).Return(errors.New("invalid old password"))

		pc.ChangePassword(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockUsecase.AssertExpectations(t)
	})
}
