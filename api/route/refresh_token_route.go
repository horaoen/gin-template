package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/horaoen/go-backend-clean-architecture/api/controller"
	"github.com/horaoen/go-backend-clean-architecture/domain"
	"github.com/horaoen/go-backend-clean-architecture/usecase"
)

func NewRefreshTokenRouter(userRepo domain.UserRepository, tokenService domain.TokenService, timeout time.Duration, group *gin.RouterGroup) {
	rtc := &controller.RefreshTokenController{
		RefreshTokenUsecase: usecase.NewRefreshTokenUsecase(userRepo, tokenService, timeout),
	}
	group.POST("/refresh", rtc.RefreshToken)
}
