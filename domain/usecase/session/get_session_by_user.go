package session

import (
	"context"
	"user-service/domain/entity"
	"user-service/domain/repository"
)

type GetSessionByUserI interface {
	Excute(ctx context.Context, userID string) ([]entity.Session, error)
}

type getSessionByUser struct {
	sessionRepo repository.SessionRepository
}

func NewGetSessionByUser(sessionRepo repository.SessionRepository) GetSessionByUserI {
	return &getSessionByUser{
		sessionRepo: sessionRepo,
	}
}

func (g *getSessionByUser) Excute(ctx context.Context, userID string) ([]entity.Session, error) {
	return g.sessionRepo.GetSessionsByUserID(ctx, userID)
}
