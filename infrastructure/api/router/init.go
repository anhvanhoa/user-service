package router

import (
	"cms-server/bootstrap"
	"cms-server/domain/service/cache"
	pkglog "cms-server/infrastructure/service/logger"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	db    *pg.DB
	app   *fiber.App
	log   pkglog.Logger
	env   *bootstrap.Env
	cache cache.RedisConfigImpl
	valid bootstrap.IValidator
}

func InitRouter(
	app *fiber.App,
	db *pg.DB,
	log pkglog.Logger,
	env *bootstrap.Env,
	cache cache.RedisConfigImpl,
	valid bootstrap.IValidator,
) {
	router := &Router{
		db:    db,
		app:   app,
		log:   log,
		env:   env,
		cache: cache,
		valid: valid,
	}
	router.initAuthRouter()
}
