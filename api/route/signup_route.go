package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/horaoen/go-backend-clean-architecture/api/controller"
	"github.com/horaoen/go-backend-clean-architecture/domain"
	"github.com/horaoen/go-backend-clean-architecture/usecase"
)

func NewSignupRouter(userRepo domain.UserRepository, tokenService domain.TokenService, timeout time.Duration, group *gin.RouterGroup) {
	sc := controller.SignupController{
		SignupUsecase: usecase.NewSignupUsecase(userRepo, tokenService, timeout),
	}
	group.POST("/signup", sc.Signup)
}
