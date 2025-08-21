package main

import (
	"cms-server/bootstrap"
	"cms-server/infrastructure/api/router"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

func main() {
	StartGRPCServer()
}

func StartFiberServer() {
	app := bootstrap.App()
	env := app.Env
	db := app.DB
	log := app.Log
	cacheApp := app.Cache
	valid := app.Valid
	defer db.Close()
	fiberApp := fiber.New(fiber.Config{
		AppName:       env.NAME_APP,
		CaseSensitive: true,
		Prefork:       false,
		StrictRouting: true,
	})

	fiberApp.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	fiberApp.Use(cache.New((cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Get("No-Cache") == "true"
		},
		Expiration:   10 * time.Minute,
		CacheControl: true,
		Storage:      cacheApp,
	})))

	// Registering the route
	router.InitRouter(fiberApp, db, log, env, cacheApp, valid)

	if err := fiberApp.Listen(":" + env.PORT_APP); err != nil {
		log.Fatal("Error starting the server: " + err.Error())
	}
}

func StartGRPCServer() {
	app := bootstrap.App()
	env := app.Env
	log := app.Log
	db := app.DB
	grpcSrv := newGRPCServer(db, env, log)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := grpcSrv.Start(ctx); err != nil {
		log.Fatal("gRPC server error: " + err.Error())
	}
}
