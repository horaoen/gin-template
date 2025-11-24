package bootstrap

import (
	"github.com/horaoen/go-backend-clean-architecture/mongo"
	"gorm.io/gorm"
)

type Application struct {
	Env   *Env
	DB    *gorm.DB
	Mongo mongo.Client
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	InitLog(app.Env)

	app.DB = NewPostgres(app.Env)
	app.Mongo = NewMongoDatabase(app.Env)
	return *app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.Mongo)
}
