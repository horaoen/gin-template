package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/horaoen/go-backend-clean-architecture/domain"
)

type ProfileController struct {
	ProfileUsecase domain.ProfileUsecase
}

// Fetch godoc
// @Summary Get Profile
// @Description Get user profile
// @Tags Profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.Profile
// @Failure 500 {object} domain.ErrorResponse
// @Router /profile [get]
func (pc *ProfileController) Fetch(c *gin.Context) {
	userID := c.GetString("x-user-id")

	profile, err := pc.ProfileUsecase.GetProfileByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (pc *ProfileController) ChangePassword(c *gin.Context) {
	var request domain.ChangePasswordRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	userID := c.GetString("x-user-id")

	err = pc.ProfileUsecase.ChangePassword(c, userID, request.OldPassword, request.NewPassword)
	if err != nil {
		if err.Error() == "invalid old password" {
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Password changed successfully"})
}
