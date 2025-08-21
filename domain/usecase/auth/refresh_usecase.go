package authUC

import (
	"cms-server/domain/entity"
	"cms-server/domain/repository"
	"cms-server/domain/service/cache"
	serviceJwt "cms-server/domain/service/jwt"
	"time"
)

type RefreshUsecase interface {
	GetSessionByToken(token string) (entity.Session, error)
	ClearSessionExpired() error
	VerifyToken(token string) (*serviceJwt.AuthClaims, error)
	GengerateAccessToken(id string, fullName string, exp time.Time) (string, error)
	GengerateRefreshToken(id string, fullName string, exp time.Time, os string) (string, error)
}

type refreshUsecaseImpl struct {
	sessionRepo repository.SessionRepository
	jwtAccess   serviceJwt.JwtService
	jwtRefresh  serviceJwt.JwtService
	cache       cache.RedisConfigImpl
}

func NewRefreshUsecase(
	sessionRepo repository.SessionRepository,
	jwtAccess serviceJwt.JwtService,
	jwtRefresh serviceJwt.JwtService,
	cache cache.RedisConfigImpl,
) RefreshUsecase {
	return &refreshUsecaseImpl{
		sessionRepo: sessionRepo,
		jwtAccess:   jwtAccess,
		jwtRefresh:  jwtRefresh,
		cache:       cache,
	}
}

func (uc *refreshUsecaseImpl) GetSessionByToken(token string) (entity.Session, error) {
	session, err := uc.sessionRepo.GetSessionAliveByToken(entity.SessionTypeAuth, token)
	if err != nil {
		return entity.Session{}, err
	}
	if err := uc.sessionRepo.DeleteSessionAuthByToken(token); err != nil {
		return entity.Session{}, err
	}
	return session, nil
}

func (uc *refreshUsecaseImpl) ClearSessionExpired() error {
	if err := uc.sessionRepo.DeleteAllSessionsExpired(); err != nil {
		return err
	}
	return nil
}

func (uc *refreshUsecaseImpl) VerifyToken(token string) (*serviceJwt.AuthClaims, error) {
	claims, err := uc.jwtRefresh.VerifyAuthToken(token)
	if err != nil {
		return claims, err
	}
	return claims, nil
}

func (uc *refreshUsecaseImpl) GengerateAccessToken(id string, fullName string, exp time.Time) (string, error) {
	return uc.jwtAccess.GenAuthToken(id, fullName, exp)
}

func (uc *refreshUsecaseImpl) GengerateRefreshToken(id string, fullName string, exp time.Time, os string) (string, error) {
	token, err := uc.jwtRefresh.GenAuthToken(id, fullName, exp)
	if err != nil {
		return "", err
	}

	session := entity.Session{
		Token:     token,
		UserID:    id,
		Os:        os,
		Type:      entity.SessionTypeAuth,
		ExpiredAt: exp,
		CreatedAt: time.Now(),
	}

	if err := uc.cache.Set(token, []byte(id), time.Until(exp)); err != nil {
		if err := uc.sessionRepo.CreateSession(session); err != nil {
			return "", err
		}
	} else {
		go uc.sessionRepo.CreateSession(session)
	}
	return token, nil
}
