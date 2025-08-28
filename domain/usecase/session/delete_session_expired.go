package sessionusecase

import (
	"auth-service/domain/repository"
	"context"
)

type DeleteSessionExpiredUsecase interface {
	DeleteAllSessionsExpired(ctx context.Context) error
}

type deleteSessionExpiredUsecase struct {
	sessionRepo repository.SessionRepository
}

func NewDeleteSessionExpiredUsecase(sessionRepo repository.SessionRepository) DeleteSessionExpiredUsecase {
	return &deleteSessionExpiredUsecase{
		sessionRepo: sessionRepo,
	}
}

func (d *deleteSessionExpiredUsecase) DeleteAllSessionsExpired(ctx context.Context) error {
	return d.sessionRepo.DeleteAllSessionsExpired(ctx)
}
