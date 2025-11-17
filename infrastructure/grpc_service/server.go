package grpcservice

import (
	"user-service/bootstrap"

	grpc_service "github.com/anhvanhoa/service-core/bootstrap/grpc"
	"github.com/anhvanhoa/service-core/domain/cache"
	"github.com/anhvanhoa/service-core/domain/log"
	"github.com/anhvanhoa/service-core/domain/token"
	"github.com/anhvanhoa/service-core/domain/user_context"
	proto_session "github.com/anhvanhoa/sf-proto/gen/session/v1"
	proto_user "github.com/anhvanhoa/sf-proto/gen/user/v1"
	"google.golang.org/grpc"
)

func NewGRPCServer(
	env *bootstrap.Env,
	log *log.LogGRPCImpl,
	cache cache.CacheI,
	userService proto_user.UserServiceServer,
	sessionService proto_session.SessionServiceServer,
) *grpc_service.GRPCServer {
	config := &grpc_service.GRPCServerConfig{
		IsProduction: env.IsProduction(),
		PortGRPC:     env.PortGrpc,
		NameService:  env.NameService,
	}
	middleware := grpc_service.NewMiddleware(
		token.NewToken(env.AccessSecret),
		log,
	)
	return grpc_service.NewGRPCServer(
		config,
		log,
		func(server *grpc.Server) {
			proto_user.RegisterUserServiceServer(server, userService)
			proto_session.RegisterSessionServiceServer(server, sessionService)
		},
		middleware.AuthorizationInterceptor(
			env.SecretService,
			func(action string, resource string) bool {
				hasPermission, err := cache.Get(resource + "." + action)
				if err != nil {
					return false
				}
				return hasPermission != nil && string(hasPermission) == "true"
			},
			func(at string) *user_context.UserContext {
				userData, err := cache.Get(at)
				if err != nil || userData == nil {
					return nil
				}
				uCtx := user_context.NewUserContext()
				uCtx.FromBytes(userData)
				return uCtx
			},
		),
	)
}
