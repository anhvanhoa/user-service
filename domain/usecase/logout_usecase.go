package usecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"auth-service/domain/service/cache"
	se "auth-service/domain/service/error"
	serviceJwt "auth-service/domain/service/jwt"
	"context"
)

var (
	ErrNotFoundSession = se.NewErr("Không tìm thấy phiên làm việc")
)

type LogoutUsecase interface {
	VerifyToken(token string) error
	Logout(token string) error
}

type logoutUsecaseImpl struct {
	sessionRepo repository.SessionRepository
	jwt         serviceJwt.JwtService
	cache       cache.RedisConfigImpl
}

func NewLogoutUsecase(
	sessionRepo repository.SessionRepository,
	jwt serviceJwt.JwtService,
	cache cache.RedisConfigImpl,
) LogoutUsecase {
	return &logoutUsecaseImpl{
		sessionRepo,
		jwt,
		cache,
	}
}

func (l *logoutUsecaseImpl) VerifyToken(token string) error {
	_, err := l.sessionRepo.GetSessionAliveByToken(entity.SessionTypeAuth, token)
	if err != nil {
		return ErrNotFoundSession
	}
	_, err = l.jwt.VerifyAuthToken(token)
	if err != nil {
		return err
	}
	return nil
}

func (l *logoutUsecaseImpl) Logout(token string) error {
	if err := l.cache.Delete(token); err != nil {
		return err
	}
	if err := l.sessionRepo.DeleteSessionAuthByToken(context.Background(), token); err != nil {
		return err
	}
	return nil
}
