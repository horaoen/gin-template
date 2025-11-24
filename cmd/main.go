package main

import (
	"time"

	"github.com/gin-gonic/gin"
	route "github.com/horaoen/go-backend-clean-architecture/api/route"
	"github.com/horaoen/go-backend-clean-architecture/bootstrap"
	"github.com/rs/zerolog/log"
)

func main() {

	app := bootstrap.App()

	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	gin := gin.Default()

	route.Setup(env, timeout, db, gin)

	err := gin.Run(env.ServerAddress)
	if err != nil {
		log.Err(err).Msg("应用启动失败")
	}
}
