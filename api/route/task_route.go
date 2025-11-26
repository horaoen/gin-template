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

func NewTaskRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	tr := repository.NewTaskRepository(db)
	tc := &controller.TaskController{
		TaskUsecase: usecase.NewTaskUsecase(tr, timeout),
	}
	group.GET("/task", tc.Fetch)
	group.POST("/task", tc.Create)
}
