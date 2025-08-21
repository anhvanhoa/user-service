package authUC

import (
	"cms-server/domain/entity"
	"cms-server/domain/repository"
	"cms-server/domain/service/argon"
	"cms-server/domain/service/cache"
	serviceError "cms-server/domain/service/error"
	serviceJwt "cms-server/domain/service/jwt"
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
		return "", serviceError.NewErrorApp("Không tìm thấy người dùng")
	}
	key := code + user.ID
	if _, err := uc.cache.Get(key); err != nil {
		if _, err := uc.sessionRepo.GetSessionForgotAliveByTokenAndIdUser(code, user.ID); err != nil {
			return "", serviceError.NewErrorApp("Phiên làm việc không hợp lệ hoặc đã hết hạn")
		}
	}
	go func() {
		uc.sessionRepo.DeleteSessionForgotByTokenAndIdUser(code, user.ID)
		uc.cache.Delete(key)
	}()
	return user.ID, nil
}

func (uc *ResetPasswordByCodeUsecaseImpl) ResetPass(IdUser, Password, ConfirmPassword string) error {
	ConfirmPassword, err := uc.argon.HashPassword(ConfirmPassword)
	if err != nil {
		return serviceError.NewErrorApp("Không thể mã hóa mật khẩu")
	}

	if _, err = uc.userRepo.UpdateUser(IdUser, entity.User{Password: ConfirmPassword}); err != nil {
		return serviceError.NewErrorApp("Không thể cập nhật mật khẩu")
	}
	return nil
}
