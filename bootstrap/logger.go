package bootstrap

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLog(env *Env) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.Level(env.LogLevel))
	log.Info().Msgf("logging level: %d", env.LogLevel)
}
