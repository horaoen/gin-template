package bootstrap

import (
	"fmt"

	zlog "github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(env *Env) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		env.PostgresHost, env.PostgresUser, env.PostgresPassword, env.PostgresDB, env.PostgresPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		zlog.Fatal().AnErr("postgres connect fail", err)
	}
	zlog.Info().Msg("PostgreSQL 连接成功")
	return db
}
