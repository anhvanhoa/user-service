package grpcservice

import (
	"user-service/bootstrap"

	grpc_server "github.com/anhvanhoa/service-core/bootstrap/grpc"
	"github.com/anhvanhoa/service-core/domain/log"
	proto_session "github.com/anhvanhoa/sf-proto/gen/session/v1"
	proto_user "github.com/anhvanhoa/sf-proto/gen/user/v1"
	"google.golang.org/grpc"
)

func NewGRPCServer(env *bootstrap.Env, log *log.LogGRPCImpl, userService proto_user.UserServiceServer, sessionService proto_session.SessionServiceServer) *grpc_server.GRPCServer {
	config := &grpc_server.GRPCServerConfig{
		IsProduction: env.IsProduction(),
		PortGRPC:     env.PortGrpc,
		NameService:  env.NameService,
	}
	return grpc_server.NewGRPCServer(
		config,
		log,
		func(server *grpc.Server) {
			proto_user.RegisterUserServiceServer(server, userService)
			proto_session.RegisterSessionServiceServer(server, sessionService)
		},
	)
}
