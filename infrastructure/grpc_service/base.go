package grpcservice

import (
	authUC "cms-server/domain/usecase/auth"
	"cms-server/infrastructure/repo"
	authpb "cms-server/proto"

	"github.com/go-pg/pg/v10"
)

type authService struct {
	authpb.UnimplementedAuthServiceServer
	checkTokenUc authUC.CheckTokenUsecase
}

func NewAuthService(db *pg.DB) authpb.AuthServiceServer {
	sessionRepo := repo.NewSessionRepository(db)
	return &authService{
		checkTokenUc: authUC.NewCheckTokenUsecase(sessionRepo),
	}
}
