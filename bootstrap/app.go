package bootstrap

import (
	"auth-service/domain/service/cache"
	pkglog "auth-service/infrastructure/service/logger"

	"github.com/go-pg/pg/v10"
	valid "github.com/go-playground/validator/v10"
	"go.uber.org/zap/zapcore"
)

type Application struct {
	Env   *Env
	DB    *pg.DB
	Log   pkglog.Logger
	Cache cache.RedisConfigImpl
	Valid IValidator
}

func App() *Application {
	env := Env{}
	NewEnv(&env)

	logConfig := pkglog.NewConfig()
	log := pkglog.InitLogger(logConfig, zapcore.DebugLevel, env.IsProduction())

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
	valid := RegisterCustomValidations(valid.New())
	return &Application{
		Env:   &env,
		DB:    db,
		Log:   log,
		Cache: cache,
		Valid: valid,
	}
}
