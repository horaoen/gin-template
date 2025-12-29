package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/horaoen/go-backend-clean-architecture/api/dto"
	"github.com/horaoen/go-backend-clean-architecture/domain"
)

type RefreshTokenController struct {
	RefreshTokenUsecase domain.RefreshTokenUsecase
}

func (rtc *RefreshTokenController) RefreshToken(c *gin.Context) {
	var request dto.RefreshTokenRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	tokens, err := rtc.RefreshTokenUsecase.Refresh(c.Request.Context(), request.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidToken):
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "invalid or expired token"})
		case errors.Is(err, domain.ErrUserNotFound):
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
