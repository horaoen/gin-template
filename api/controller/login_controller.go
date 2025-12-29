package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/horaoen/go-backend-clean-architecture/api/dto"
	"github.com/horaoen/go-backend-clean-architecture/domain"
)

type LoginController struct {
	LoginUsecase domain.LoginUsecase
}

func (lc *LoginController) Login(c *gin.Context) {
	var request dto.LoginRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	tokens, err := lc.LoginUsecase.Login(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "invalid email or password"})
		default:
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
