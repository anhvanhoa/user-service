package repository

import (
	"cms-server/domain/entity"
	"context"
)

type SessionRepository interface {
	CreateSession(data entity.Session) error
	GetSessionAliveByToken(typeSession entity.SessionType, token string) (entity.Session, error)
	GetSessionAliveByTokenAndIdUser(typeSession entity.SessionType, token, idUser string) (entity.Session, error)
	GetSessionForgotAliveByTokenAndIdUser(token, idUser string) (entity.Session, error)
	TokenExists(token string) bool
	DeleteSessionByTypeAndUserID(sessionType entity.SessionType, userID string) error
	DeleteSessionByTypeAndToken(sessionType entity.SessionType, token string) error
	DeleteSessionVerifyByUserID(userID string) error
	DeleteSessionAuthByToken(token string) error
	DeleteSessionVerifyByToken(token string) error
	DeleteSessionForgotByToken(token string) error
	DeleteAllSessionsExpired() error
	DeleteAllSessionsForgot() error
	DeleteSessionForgotByTokenAndIdUser(token, idUser string) error
	Tx(ctx context.Context) SessionRepository
}
