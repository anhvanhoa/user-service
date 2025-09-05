package sessionusecase

import (
	"context"
	"user-service/domain/entity"
	"user-service/domain/repository"
)

type DeleteSessionByTypeAndTokenUsecase interface {
	Excute(ctx context.Context, sessionType entity.SessionType, token string) error
}

type deleteSessionByTypeAndTokenUsecase struct {
	sessionRepo repository.SessionRepository
}

func NewDeleteSessionByTypeAndTokenUsecase(sessionRepo repository.SessionRepository) DeleteSessionByTypeAndTokenUsecase {
	return &deleteSessionByTypeAndTokenUsecase{
		sessionRepo: sessionRepo,
	}
}

func (d *deleteSessionByTypeAndTokenUsecase) Excute(ctx context.Context, sessionType entity.SessionType, token string) error {
	return d.sessionRepo.DeleteSessionByTypeAndToken(ctx, sessionType, token)
}
