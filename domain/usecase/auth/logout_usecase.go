package authUC

import (
	"cms-server/domain/entity"
	"cms-server/domain/repository"
	"cms-server/domain/service/cache"
	serviceError "cms-server/domain/service/error"
	serviceJwt "cms-server/domain/service/jwt"
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
		return serviceError.ErrNotFoundSession
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
	if err := l.sessionRepo.DeleteSessionAuthByToken(token); err != nil {
		return err
	}
	return nil
}
