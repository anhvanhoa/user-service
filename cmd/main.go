package main

import (
	"auth-service/bootstrap"
	grpcservice "auth-service/infrastructure/grpc_service"
	"context"
)

func main() {
	StartGRPCServer()
}

func StartGRPCServer() {
	app := bootstrap.App()
	env := app.Env
	log := app.Log
	db := app.DB
	authService := grpcservice.NewAuthService(db, env, log)
	grpcSrv := grpcservice.NewGRPCServer(env.PORT_GRPC, authService, log)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := grpcSrv.Start(ctx); err != nil {
		log.Fatal("gRPC server error: " + err.Error())
	}
}
