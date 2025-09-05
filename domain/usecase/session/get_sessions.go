package sessionusecase

import (
	"context"
	"user-service/domain/entity"
	"user-service/domain/repository"
)

type GetSessionsUsecase interface {
	Excute(ctx context.Context) ([]entity.Session, error)
}

type getSessionsUsecase struct {
	sessionRepo repository.SessionRepository
}

func NewGetSessionsUsecase(sessionRepo repository.SessionRepository) GetSessionsUsecase {
	return &getSessionsUsecase{
		sessionRepo: sessionRepo,
	}
}

func (g *getSessionsUsecase) Excute(ctx context.Context) ([]entity.Session, error) {
	return g.sessionRepo.GetAllSessions(ctx)
}
