package bootstrap

import (
	"auth-service/infrastructure/service/database"
	pkglog "auth-service/infrastructure/service/logger"

	"github.com/go-pg/pg/extra/pgdebug/v10"
	"github.com/go-pg/pg/v10"
)

func NewPostgresDB(env *Env, log pkglog.Logger) *pg.DB {
	opt, err := pg.ParseURL(env.URL_DB)
	if err != nil {
		log.Fatal("Error parsing the database URL: " + err.Error())
	}

	db := pg.Connect(opt)
	if err := db.Ping(db.Context()); err != nil {
		log.Fatal("Error connecting to the database: " + err.Error())
	}

	if !env.IsProduction() {
		db.AddQueryHook(pgdebug.NewDebugHook())
		db.AddQueryHook(database.NewQueryHook())
	}
	return db
}
