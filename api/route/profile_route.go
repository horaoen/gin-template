package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/horaoen/go-backend-clean-architecture/api/controller"
	"github.com/horaoen/go-backend-clean-architecture/domain"
	"github.com/horaoen/go-backend-clean-architecture/usecase"
)

func NewProfileRouter(userRepo domain.UserRepository, timeout time.Duration, group *gin.RouterGroup) {
	pc := &controller.ProfileController{
		ProfileUsecase: usecase.NewProfileUsecase(userRepo, timeout),
	}
	group.GET("/profile", pc.Fetch)
	group.POST("/profile/change-password", pc.ChangePassword)
}
