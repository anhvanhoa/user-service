package sessionusecase

import (
	"context"
	"user-service/domain/repository"
)

type DeleteSessionExpiredUsecase interface {
	Excute(ctx context.Context) error
}

type deleteSessionExpiredUsecase struct {
	sessionRepo repository.SessionRepository
}

func NewDeleteSessionExpiredUsecase(sessionRepo repository.SessionRepository) DeleteSessionExpiredUsecase {
	return &deleteSessionExpiredUsecase{
		sessionRepo: sessionRepo,
	}
}

func (d *deleteSessionExpiredUsecase) Excute(ctx context.Context) error {
	return d.sessionRepo.DeleteAllSessionsExpired(ctx)
}
