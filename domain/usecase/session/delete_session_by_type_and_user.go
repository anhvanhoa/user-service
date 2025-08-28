package sessionusecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"context"
)

type DeleteSessionByTypeAndUserUsecase interface {
	DeleteSessionByTypeAndUserID(ctx context.Context, sessionType entity.SessionType, userID string) error
}

type deleteSessionByTypeAndUserUsecase struct {
	sessionRepo repository.SessionRepository
}

func NewDeleteSessionByTypeAndUserUsecase(sessionRepo repository.SessionRepository) DeleteSessionByTypeAndUserUsecase {
	return &deleteSessionByTypeAndUserUsecase{
		sessionRepo: sessionRepo,
	}
}

func (d *deleteSessionByTypeAndUserUsecase) DeleteSessionByTypeAndUserID(ctx context.Context, sessionType entity.SessionType, userID string) error {
	return d.sessionRepo.DeleteSessionByTypeAndUserID(ctx, sessionType, userID)
}
