package main

import (
	"context"
	"user-service/bootstrap"
	"user-service/infrastructure/grpc_client"
	grpcservice "user-service/infrastructure/grpc_service"
	session_server "user-service/infrastructure/grpc_service/session"
	user_server "user-service/infrastructure/grpc_service/user"

	gc "github.com/anhvanhoa/service-core/domain/grpc_client"
)

func main() {
	app := bootstrap.App()
	env := app.Env
	log := app.Log
	db := app.DB
	cache := app.Cache

	clientFactory := gc.NewClientFactory(env.GrpcClients...)
	client := clientFactory.GetClient(env.PermissionServiceAddr)
	permissionClient := grpc_client.NewPermissionClient(client)

	userService := user_server.NewUserServer(db, app.Helper)
	sessionService := session_server.NewSessionServer(db, cache)
	grpcSrv := grpcservice.NewGRPCServer(env, log, cache, userService, sessionService)
	ctx, cancel := context.WithCancel(context.Background())
	permissions := app.Helper.ConvertResourcesToPermissions(grpcSrv.GetResources())
	if _, err := permissionClient.PermissionServiceClient.RegisterPermission(ctx, permissions); err != nil {
		log.Fatal("Failed to register permission: " + err.Error())
	}
	defer cancel()
	if err := grpcSrv.Start(ctx); err != nil {
		log.Fatal("gRPC server error: " + err.Error())
	}
}
