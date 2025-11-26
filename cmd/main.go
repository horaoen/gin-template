package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/horaoen/go-backend-clean-architecture/api/route"
	"github.com/horaoen/go-backend-clean-architecture/bootstrap"
	"github.com/rs/zerolog/log"
)

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
