package main

import (
	"context"
	"user-service/bootstrap"
	grpcservice "user-service/infrastructure/grpc_service"
	role_server "user-service/infrastructure/grpc_service/role"
	session_server "user-service/infrastructure/grpc_service/session"
	user_server "user-service/infrastructure/grpc_service/user"

	"github.com/anhvanhoa/service-core/domain/discovery"
)

func main() {
	app := bootstrap.App()
	env := app.Env
	log := app.Log
	db := app.DB
	cache := app.Cache
	discoveryConfig := discovery.DiscoveryConfig{
		ServiceName:   env.NameService,
		ServicePort:   env.PortGrpc,
		ServiceHost:   env.HostGrpc,
		IntervalCheck: env.IntervalCheck,
		TimeoutCheck:  env.TimeoutCheck,
	}
	discoveryClient, err := discovery.NewDiscovery(&discoveryConfig)
	if err != nil {
		log.Fatal("Failed to create discovery client: " + err.Error())
	}
	discoveryClient.Register()
	defer discoveryClient.Close(env.NameService)
	userService := user_server.NewUserServer(db)
	sessionService := session_server.NewSessionServer(db, cache)
	roleService := role_server.NewRoleServer(db)
	grpcSrv := grpcservice.NewGRPCServer(env, log, userService, sessionService, roleService)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := grpcSrv.Start(ctx); err != nil {
		log.Fatal("gRPC server error: " + err.Error())
	}
}
