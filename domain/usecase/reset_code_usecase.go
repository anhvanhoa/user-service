package usecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"auth-service/domain/service/argon"
	"auth-service/domain/service/cache"
	se "auth-service/domain/service/error"
	serviceJwt "auth-service/domain/service/jwt"
	"context"
)

type ResetPasswordByCodeUsecase interface {
	VerifySession(code, email string) (string, error)
	ResetPass(IdUser, Password, NewPassword string) error
}

type ResetPasswordByCodeUsecaseImpl struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	cache       cache.RedisConfigImpl
	jwt         serviceJwt.JwtService
	argon       argon.Argon
}

var (
	ErrNotFoundUser   = se.NewErr("Không tìm thấy người dùng")
	ErrHashPassword   = se.NewErr("Không thể mã hóa mật khẩu")
	ErrUpdatePassword = se.NewErr("Không thể cập nhật mật khẩu")
)

func NewResetPasswordCodeUsecase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	cache cache.RedisConfigImpl,
	jwt serviceJwt.JwtService,
	argon argon.Argon,
) ResetPasswordByCodeUsecase {
	return &ResetPasswordByCodeUsecaseImpl{
		userRepo,
		sessionRepo,
		cache,
		jwt,
		argon,
	}
}

func (uc *ResetPasswordByCodeUsecaseImpl) VerifySession(code, email string) (string, error) {
	user, err := uc.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", ErrNotFoundUser
	}
	key := code + user.ID
	if _, err := uc.cache.Get(key); err != nil {
		if _, err := uc.sessionRepo.GetSessionForgotAliveByTokenAndIdUser(code, user.ID); err != nil {
			return "", ErrNotFoundSession
		}
	}
	go func() {
		uc.sessionRepo.DeleteSessionForgotByTokenAndIdUser(context.Background(), code, user.ID)
		uc.cache.Delete(key)
	}()
	return user.ID, nil
}

func (uc *ResetPasswordByCodeUsecaseImpl) ResetPass(IdUser, Password, ConfirmPassword string) error {
	ConfirmPassword, err := uc.argon.HashPassword(ConfirmPassword)
	if err != nil {
		return ErrHashPassword
	}

	if _, err = uc.userRepo.UpdateUser(IdUser, entity.User{Password: ConfirmPassword}); err != nil {
		return ErrUpdatePassword
	}
	return nil
}
