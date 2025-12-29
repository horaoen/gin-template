package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/horaoen/go-backend-clean-architecture/api/dto"
	"github.com/horaoen/go-backend-clean-architecture/domain"
)

type SignupController struct {
	SignupUsecase domain.SignupUsecase
}

func (sc *SignupController) Signup(c *gin.Context) {
	var request dto.SignupRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	tokens, err := sc.SignupUsecase.Signup(c.Request.Context(), request.Name, request.Email, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, domain.ErrorResponse{Message: "user already exists with the given email"})
		default:
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.SignupResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
