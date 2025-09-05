package bootstrap

import (
	"github.com/anhvanhoa/service-core/bootstrap/db"
	"github.com/anhvanhoa/service-core/domain/cache"
	"github.com/anhvanhoa/service-core/domain/log"
	"github.com/go-pg/pg/v10"
	"go.uber.org/zap/zapcore"
)

type Application struct {
	Env   *Env
	DB    *pg.DB
	Log   *log.LogGRPCImpl
	Cache cache.CacheI
}

func App() *Application {
	env := Env{}
	NewEnv(&env)
	logConfig := log.NewConfig()
	log := log.InitLogGRPC(logConfig, zapcore.DebugLevel, env.IsProduction())
	db := db.NewPostgresDB(db.ConfigDB{
		URL:  env.UrlDb,
		Mode: env.NodeEnv,
	})
	configRedis := cache.NewConfigCache(
		env.DbCache.Addr,
		env.DbCache.Password,
		env.DbCache.Db,
		env.DbCache.Network,
		env.DbCache.MaxIdle,
		env.DbCache.MaxActive,
		env.DbCache.IdleTimeout,
	)
	cache := cache.NewCache(configRedis)
	return &Application{
		Env:   &env,
		DB:    db,
		Log:   log,
		Cache: cache,
	}
}
