package sessionusecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"context"
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
