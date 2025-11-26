package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/horaoen/go-backend-clean-architecture/api/controller"
	"github.com/horaoen/go-backend-clean-architecture/bootstrap"
	"github.com/horaoen/go-backend-clean-architecture/repository"
	"github.com/horaoen/go-backend-clean-architecture/usecase"
	"gorm.io/gorm"
)

func NewProfileRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	pc := &controller.ProfileController{
		ProfileUsecase: usecase.NewProfileUsecase(ur, timeout),
	}
	group.GET("/profile", pc.Fetch)
}
