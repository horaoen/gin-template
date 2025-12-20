package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/horaoen/go-backend-clean-architecture/api/route"
	"github.com/horaoen/go-backend-clean-architecture/bootstrap"
	"github.com/rs/zerolog/log"
)

// @title           Go Backend Clean Architecture API
// @version         1.0
// @description     This is a sample server for Go Backend Clean Architecture.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /
// @schemes         http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	app := bootstrap.App()

	env := app.Env

	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	engine := gin.Default()

	route.Setup(env, timeout, app.DB, engine)

	err := engine.Run(env.ServerAddress)
	if err != nil {
		log.Err(err).Msg("应用启动失败")
	}
}
