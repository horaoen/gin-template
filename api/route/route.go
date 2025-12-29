package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/horaoen/go-backend-clean-architecture/api/middleware"
	"github.com/horaoen/go-backend-clean-architecture/bootstrap"
	_ "github.com/horaoen/go-backend-clean-architecture/docs"
	"github.com/horaoen/go-backend-clean-architecture/repository"
	"github.com/horaoen/go-backend-clean-architecture/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, gin *gin.Engine) {
	userRepo := repository.NewUserRepository(db)
	tokenService := usecase.NewTokenService(
		env.AccessTokenSecret,
		env.RefreshTokenSecret,
		env.AccessTokenExpiryHour,
		env.RefreshTokenExpiryHour,
	)

	publicRouter := gin.Group("")
	publicRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	NewSignupRouter(userRepo, tokenService, timeout, publicRouter)
	NewLoginRouter(userRepo, tokenService, timeout, publicRouter)
	NewRefreshTokenRouter(userRepo, tokenService, timeout, publicRouter)

	protectedRouter := gin.Group("")
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	NewProfileRouter(userRepo, timeout, protectedRouter)
}
