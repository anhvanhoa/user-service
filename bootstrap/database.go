package bootstrap

import (
	"auth-service/domain/service/logger"
	"auth-service/infrastructure/service/database"

	"github.com/go-pg/pg/extra/pgdebug/v10"
	"github.com/go-pg/pg/v10"
)

func NewPostgresDB(env *Env, log logger.Log) *pg.DB {
	// Connect to the database
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
