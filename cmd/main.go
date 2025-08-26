package main

import (
	"auth-service/bootstrap"
	"auth-service/infrastructure/discovery"
	"auth-service/infrastructure/grpc_client"
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
	queueClient := app.Queue
	discoveryClient, err := discovery.NewDiscovery(log, env)
	if err != nil {
		log.Fatal("Failed to create discovery client: " + err.Error())
	}
	discoveryClient.Register(env.NAME_SERVICE)
	defer discoveryClient.Close(env.NAME_SERVICE)
	clientFactory := grpc_client.NewClientFactory(log, env.GRPC_CLIENTS...)
	client := clientFactory.GetClient(env.MAIL_SERVICE_ADDR)
	mailService := grpc_client.NewMailService(client)

	authService := grpcservice.NewAuthService(db, env, log, mailService, queueClient)
	grpcSrv := grpcservice.NewGRPCServer(env, authService, log)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := grpcSrv.Start(ctx); err != nil {
		log.Fatal("gRPC server error: " + err.Error())
	}
}
