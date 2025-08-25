package bootstrap

import (
	"auth-service/domain/service/cache"
	loggerI "auth-service/domain/service/logger"
	"auth-service/infrastructure/service/logger"

	"github.com/go-pg/pg/v10"
	"go.uber.org/zap/zapcore"
)

type Application struct {
	Env   *Env
	DB    *pg.DB
	Log   loggerI.Log
	Cache cache.RedisConfigImpl
	Queue *queueClient
}

func App() *Application {
	env := Env{}
	NewEnv(&env)
	logConfig := logger.NewConfig()
	log := logger.InitLogger(logConfig, zapcore.DebugLevel, env.IsProduction())

	db := NewPostgresDB(&env, log)
	configRedis := NewRedisConfig(
		env.DB_CACHE.Addr,
		env.DB_CACHE.Password,
		env.DB_CACHE.DB,
		env.DB_CACHE.Network,
		env.DB_CACHE.MaxIdle,
		env.DB_CACHE.MaxActive,
		env.DB_CACHE.IdleTimeout,
	)
	cache := NewRedis(configRedis)
	queue := NewQueueClient(&env, log)
	return &Application{
		Env:   &env,
		DB:    db,
		Log:   log,
		Cache: cache,
		Queue: queue,
	}
}
