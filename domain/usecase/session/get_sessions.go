package sessionusecase

import (
	"context"
	"user-service/domain/entity"
	"user-service/domain/repository"
)

type GetSessionsUsecase interface {
	GetAllSessions(ctx context.Context) ([]entity.Session, error)
	GetSessionsByUserID(ctx context.Context, userID string) ([]entity.Session, error)
	GetSessionsByType(ctx context.Context, sessionType entity.SessionType) ([]entity.Session, error)
}

type getSessionsUsecase struct {
	sessionRepo repository.SessionRepository
}

func NewGetSessionsUsecase(sessionRepo repository.SessionRepository) GetSessionsUsecase {
	return &getSessionsUsecase{
		sessionRepo: sessionRepo,
	}
}

func (g *getSessionsUsecase) GetAllSessions(ctx context.Context) ([]entity.Session, error) {
	return g.sessionRepo.GetAllSessions(ctx)
}

func (g *getSessionsUsecase) GetSessionsByUserID(ctx context.Context, userID string) ([]entity.Session, error) {
	return g.sessionRepo.GetSessionsByUserID(ctx, userID)
}

func (g *getSessionsUsecase) GetSessionsByType(ctx context.Context, sessionType entity.SessionType) ([]entity.Session, error) {
	return g.sessionRepo.GetSessionsByType(ctx, sessionType)
}
