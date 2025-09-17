package session

import (
	"context"
	"user-service/domain/entity"
	"user-service/domain/repository"
)

type GetSessionsByTypeUsecase interface {
	Excute(ctx context.Context, sessionType entity.SessionType) ([]entity.Session, error)
}

type getSessionsByTypeUsecase struct {
	sessionRepo repository.SessionRepository
}

func NewGetSessionsByTypeUsecase(sessionRepo repository.SessionRepository) GetSessionsByTypeUsecase {
	return &getSessionsByTypeUsecase{
		sessionRepo: sessionRepo,
	}
}

func (g *getSessionsByTypeUsecase) Excute(ctx context.Context, sessionType entity.SessionType) ([]entity.Session, error) {
	return g.sessionRepo.GetSessionsByType(ctx, sessionType)
}
