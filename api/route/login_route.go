package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/horaoen/go-backend-clean-architecture/api/controller"
	"github.com/horaoen/go-backend-clean-architecture/domain"
	"github.com/horaoen/go-backend-clean-architecture/usecase"
)

func NewLoginRouter(userRepo domain.UserRepository, tokenService domain.TokenService, timeout time.Duration, group *gin.RouterGroup) {
	lc := &controller.LoginController{
		LoginUsecase: usecase.NewLoginUsecase(userRepo, tokenService, timeout),
	}
	group.POST("/login", lc.Login)
}
