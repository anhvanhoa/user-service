package usecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"auth-service/domain/service/argon"
	"auth-service/domain/service/cache"
	serviceJwt "auth-service/domain/service/jwt"
	"time"
)

type LoginUsecase interface {
	GetUserByEmailOrPhone(val string) (entity.User, error)
	CheckHashPassword(password, hash string) bool
	GengerateAccessToken(id string, fullName string, exp time.Time) (string, error)
	GengerateRefreshToken(id string, fullName string, exp time.Time, os string) (string, error)
}

type loginUsecaseImpl struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	jwtAccess   serviceJwt.JwtService
	jwtRefresh  serviceJwt.JwtService
	argon       argon.Argon
	cache       cache.RedisConfigImpl
}

func NewLoginUsecase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	jwtAccess serviceJwt.JwtService,
	jwtRefresh serviceJwt.JwtService,
	argon argon.Argon,
	cache cache.RedisConfigImpl,
) LoginUsecase {
	return &loginUsecaseImpl{
		userRepo,
		sessionRepo,
		jwtAccess,
		jwtRefresh,
		argon,
		cache,
	}
}

func (uc *loginUsecaseImpl) GetUserByEmailOrPhone(val string) (entity.User, error) {
	return uc.userRepo.GetUserByEmailOrPhone(val)
}

func (uc *loginUsecaseImpl) CheckHashPassword(password, hash string) bool {
	mach, err := uc.argon.VerifyPassword(hash, password)
	if err != nil {
		return false
	}
	return mach
}

func (uc *loginUsecaseImpl) GengerateAccessToken(id string, fullName string, exp time.Time) (string, error) {
	return uc.jwtAccess.GenAuthToken(id, fullName, exp)
}

func (uc *loginUsecaseImpl) GengerateRefreshToken(id string, fullName string, exp time.Time, os string) (string, error) {
	token, err := uc.jwtRefresh.GenAuthToken(id, fullName, exp)
	if err != nil {
		return "", err
	}

	session := entity.Session{
		Token:     token,
		UserID:    id,
		Os:        os,
		Type:      (entity.SessionTypeAuth),
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
