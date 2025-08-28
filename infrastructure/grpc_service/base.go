package grpcservice

import (
	"user-service/bootstrap"
	"user-service/infrastructure/grpc_client"

	"github.com/anhvanhoa/service-core/domain/cache"
	"github.com/anhvanhoa/service-core/domain/log"
	"github.com/anhvanhoa/service-core/domain/queue"
	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"

	"github.com/go-pg/pg/v10"
)

type authService struct {
	proto_auth.UnimplementedAuthServiceServer
	env *bootstrap.Env
	log *log.LogGRPCImpl
}

func NewAuthService(
	db *pg.DB,
	env *bootstrap.Env,
	log *log.LogGRPCImpl,
	mailService *grpc_client.MailService,
	queueClient queue.QueueClient,
	cache cache.CacheI,
) proto_auth.AuthServiceServer {
	return &authService{
		env: env,
		log: log,
	}
}
