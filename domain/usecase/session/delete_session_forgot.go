package sessionusecase

import (
	"context"
	"user-service/domain/repository"
)

type DeleteSessionForgotUsecase interface {
	Excute(ctx context.Context) error
}

type deleteSessionForgotUsecase struct {
	sessionRepo repository.SessionRepository
}

func NewDeleteSessionForgotUsecase(sessionRepo repository.SessionRepository) DeleteSessionForgotUsecase {
	return &deleteSessionForgotUsecase{
		sessionRepo: sessionRepo,
	}
}

func (d *deleteSessionForgotUsecase) Excute(ctx context.Context) error {
	return d.sessionRepo.DeleteAllSessionsForgot(ctx)
}
