// Package bootstrap
package bootstrap

import (
	"github.com/horaoen/go-backend-clean-architecture/domain"
	"gorm.io/gorm"
)

type Application struct {
	Env *Env
	DB  *gorm.DB
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	InitLog(app.Env)

	app.DB = NewPostgres(app.Env)

	// 自动迁移数据库表
	err := app.DB.AutoMigrate(&domain.User{}, &domain.Task{})
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}

	return *app
}

func (app *Application) CloseDBConnection() {
	if app.DB != nil {
		sqlDB, err := app.DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}
}
