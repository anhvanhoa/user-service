package main

import (
	"auth-service/bootstrap"
	"auth-service/constants"
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
	discoveryClient, err := discovery.NewDiscovery(log)
	if err != nil {
		log.Fatal("Failed to create discovery client: " + err.Error())
	}
	discoveryClient.Register(constants.AuthService, env.PORT_GRPC)
	defer discoveryClient.Close(constants.AuthService)
	clientConfig := []*grpc_client.Config{}
	for name, client := range env.GRPC_CLIENTS {
		clientConfig = append(clientConfig, &grpc_client.Config{
			Name:          name,
			ServerAddress: client.ServerAddress,
			Timeout:       client.Timeout,
			MaxRetries:    client.MaxRetries,
			KeepAlive:     client.KeepAlive,
		})
	}
	clientFactory := grpc_client.NewClientFactory(log, clientConfig...)
	client := clientFactory.GetClient(constants.MailService)
	mailService := grpc_client.NewMailService(client)
	authService := grpcservice.NewAuthService(db, env, log, mailService, queueClient)
	grpcSrv := grpcservice.NewGRPCServer(env, authService, log)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := grpcSrv.Start(ctx); err != nil {
		log.Fatal("gRPC server error: " + err.Error())
	}
}
