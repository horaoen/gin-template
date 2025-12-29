package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/horaoen/go-backend-clean-architecture/api/dto"
	"github.com/horaoen/go-backend-clean-architecture/domain"
)

type ProfileController struct {
	ProfileUsecase domain.ProfileUsecase
}

func (pc *ProfileController) Fetch(c *gin.Context) {
	userID := c.GetString("x-user-id")

	profile, err := pc.ProfileUsecase.GetProfileByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ProfileResponse{
		Name:  profile.Name,
		Email: profile.Email,
	})
}

func (pc *ProfileController) ChangePassword(c *gin.Context) {
	var request dto.ChangePasswordRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	userID := c.GetString("x-user-id")

	err := pc.ProfileUsecase.ChangePassword(c, userID, request.OldPassword, request.NewPassword)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) || err.Error() == "invalid old password" {
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "invalid old password"})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Password changed successfully"})
}
