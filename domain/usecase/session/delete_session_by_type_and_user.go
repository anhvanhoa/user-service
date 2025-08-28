package sessionusecase

import (
	"context"
	"user-service/domain/entity"
	"user-service/domain/repository"
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
