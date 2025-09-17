package session

import (
	"context"
	"user-service/domain/entity"
)

type SessionUsecaseI interface {
	GetSessions(ctx context.Context) ([]entity.Session, error)
	GetSessionsByUserId(ctx context.Context, userID string) ([]entity.Session, error)
	DeleteSessionByTypeAndToken(ctx context.Context, sessionType entity.SessionType, token string) error
	DeleteSessionByTypeAndUser(ctx context.Context, sessionType entity.SessionType, userID string) error
	DeleteSessionExpired(ctx context.Context) error
	GetSessionsByType(ctx context.Context, sessionType entity.SessionType) ([]entity.Session, error)
}

type sessionUsecase struct {
	getSessionsUsecase                 GetSessionsUsecase
	getSessionsByUserIdUsecase         GetSessionByUserI
	deleteSessionByTypeAndTokenUsecase DeleteSessionByTypeAndTokenUsecase
	deleteSessionByTypeAndUserUsecase  DeleteSessionByTypeAndUserUsecase
	deleteSessionExpiredUsecase        DeleteSessionExpiredUsecase
	getSessionsByTypeUsecase           GetSessionsByTypeUsecase
}

func NewSessionUsecase(
	getSessionsUsecase GetSessionsUsecase,
	getSessionsByUserIdUsecase GetSessionByUserI,
	deleteSessionByTypeAndTokenUsecase DeleteSessionByTypeAndTokenUsecase,
	deleteSessionByTypeAndUserUsecase DeleteSessionByTypeAndUserUsecase,
	deleteSessionExpiredUsecase DeleteSessionExpiredUsecase,
	getSessionsByTypeUsecase GetSessionsByTypeUsecase,
) SessionUsecaseI {
	return &sessionUsecase{
		getSessionsUsecase:                 getSessionsUsecase,
		getSessionsByUserIdUsecase:         getSessionsByUserIdUsecase,
		deleteSessionByTypeAndTokenUsecase: deleteSessionByTypeAndTokenUsecase,
		deleteSessionByTypeAndUserUsecase:  deleteSessionByTypeAndUserUsecase,
		deleteSessionExpiredUsecase:        deleteSessionExpiredUsecase,
		getSessionsByTypeUsecase:           getSessionsByTypeUsecase,
	}
}

func (s *sessionUsecase) GetSessions(ctx context.Context) ([]entity.Session, error) {
	return s.getSessionsUsecase.Excute(ctx)
}

func (s *sessionUsecase) GetSessionsByUserId(ctx context.Context, userID string) ([]entity.Session, error) {
	return s.getSessionsByUserIdUsecase.Excute(ctx, userID)
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

func (s *sessionUsecase) GetSessionsByType(ctx context.Context, sessionType entity.SessionType) ([]entity.Session, error) {
	return s.getSessionsByTypeUsecase.Excute(ctx, sessionType)
}
