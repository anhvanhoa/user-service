package authUC

import (
	"cms-server/domain/entity"
	"cms-server/domain/repository"
	"cms-server/domain/service/cache"
	serviceJwt "cms-server/domain/service/jwt"
	pkgjwt "cms-server/infrastructure/service/jwt"
	"time"
)

type VerifyAccountUsecase interface {
	VerifyRegister(t string) (*serviceJwt.VerifyClaims, error)
	GetUserById(id string) (entity.User, error)
	VerifyAccount(id string) error
}

type verifyAccountUsecaseImpl struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	jwt         serviceJwt.JwtService
	cache       cache.RedisConfigImpl
}

func NewVerifyAccountUsecase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	jwt serviceJwt.JwtService,
	cache cache.RedisConfigImpl,
) VerifyAccountUsecase {
	return &verifyAccountUsecaseImpl{
		userRepo,
		sessionRepo,
		jwt,
		cache,
	}
}

func (u *verifyAccountUsecaseImpl) VerifyRegister(t string) (*serviceJwt.VerifyClaims, error) {
	if _, err := u.cache.Get(t); err != nil {
		if isExist := u.sessionRepo.TokenExists(t); !isExist {
			return nil, pkgjwt.ErrTokenNotFound
		}
	} else {
		go func() {
			u.sessionRepo.DeleteSessionAuthByToken(t)
			u.cache.Delete(t)
		}()
	}

	data, err := u.jwt.VerifyRegisterToken(t)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *verifyAccountUsecaseImpl) GetUserById(id string) (entity.User, error) {
	user, err := u.userRepo.GetUserByID(id)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *verifyAccountUsecaseImpl) VerifyAccount(id string) error {
	t := time.Now()
	user := entity.User{
		ID:         id,
		CodeVerify: "",
		Veryfied:   &t,
	}
	_, err := u.userRepo.UpdateUser(id, user)
	return err
}
