package sessionusecase

import (
	"context"
	"user-service/domain/entity"
	"user-service/domain/repository"
)

type DeleteSessionByTypeAndTokenUsecase interface {
	DeleteSessionByTypeAndToken(ctx context.Context, sessionType entity.SessionType, token string) error
}

type deleteSessionByTypeAndTokenUsecase struct {
	sessionRepo repository.SessionRepository
}

func NewDeleteSessionByTypeAndTokenUsecase(sessionRepo repository.SessionRepository) DeleteSessionByTypeAndTokenUsecase {
	return &deleteSessionByTypeAndTokenUsecase{
		sessionRepo: sessionRepo,
	}
}

func (d *deleteSessionByTypeAndTokenUsecase) DeleteSessionByTypeAndToken(ctx context.Context, sessionType entity.SessionType, token string) error {
	return d.sessionRepo.DeleteSessionByTypeAndToken(ctx, sessionType, token)
}
