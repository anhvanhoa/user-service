package sessionusecase

import (
	"context"
	"user-service/domain/entity"
)

type SessionUsecaseI interface {
	GetSessions(ctx context.Context) ([]entity.Session, error)
	DeleteSessionByTypeAndToken(ctx context.Context, sessionType entity.SessionType, token string) error
	DeleteSessionByTypeAndUser(ctx context.Context, sessionType entity.SessionType, userID string) error
	DeleteSessionExpired(ctx context.Context) error
	DeleteSessionForgot(ctx context.Context) error
}

type sessionUsecase struct {
	getSessionsUsecase                 GetSessionsUsecase
	deleteSessionByTypeAndTokenUsecase DeleteSessionByTypeAndTokenUsecase
	deleteSessionByTypeAndUserUsecase  DeleteSessionByTypeAndUserUsecase
	deleteSessionExpiredUsecase        DeleteSessionExpiredUsecase
	deleteSessionForgotUsecase         DeleteSessionForgotUsecase
}

func NewSessionUsecase(
	getSessionsUsecase GetSessionsUsecase,
	deleteSessionByTypeAndTokenUsecase DeleteSessionByTypeAndTokenUsecase,
	deleteSessionByTypeAndUserUsecase DeleteSessionByTypeAndUserUsecase,
	deleteSessionExpiredUsecase DeleteSessionExpiredUsecase,
	deleteSessionForgotUsecase DeleteSessionForgotUsecase,
) SessionUsecaseI {
	return &sessionUsecase{
		getSessionsUsecase:                 getSessionsUsecase,
		deleteSessionByTypeAndTokenUsecase: deleteSessionByTypeAndTokenUsecase,
		deleteSessionByTypeAndUserUsecase:  deleteSessionByTypeAndUserUsecase,
		deleteSessionExpiredUsecase:        deleteSessionExpiredUsecase,
		deleteSessionForgotUsecase:         deleteSessionForgotUsecase,
	}
}

func (s *sessionUsecase) GetSessions(ctx context.Context) ([]entity.Session, error) {
	return s.getSessionsUsecase.Excute(ctx)
}

func (s *sessionUsecase) DeleteSessionByTypeAndToken(ctx context.Context, sessionType entity.SessionType, token string) error {
	return s.deleteSessionByTypeAndTokenUsecase.Excute(ctx, sessionType, token)
}

func (s *sessionUsecase) DeleteSessionByTypeAndUser(ctx context.Context, sessionType entity.SessionType, userID string) error {
	return s.deleteSessionByTypeAndUserUsecase.Excute(ctx, sessionType, userID)
}

func (s *sessionUsecase) DeleteSessionExpired(ctx context.Context) error {
	return s.deleteSessionExpiredUsecase.Excute(ctx)
}

func (s *sessionUsecase) DeleteSessionForgot(ctx context.Context) error {
	return s.deleteSessionForgotUsecase.Excute(ctx)
}
