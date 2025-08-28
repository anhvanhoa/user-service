package sessionusecase

import (
	"auth-service/domain/repository"
	"context"
)

type DeleteSessionForgotUsecase interface {
	DeleteAllSessionsForgot(ctx context.Context) error
}

type deleteSessionForgotUsecase struct {
	sessionRepo repository.SessionRepository
}

func NewDeleteSessionForgotUsecase(sessionRepo repository.SessionRepository) DeleteSessionForgotUsecase {
	return &deleteSessionForgotUsecase{
		sessionRepo: sessionRepo,
	}
}

func (d *deleteSessionForgotUsecase) DeleteAllSessionsForgot(ctx context.Context) error {
	return d.sessionRepo.DeleteAllSessionsForgot(ctx)
}
