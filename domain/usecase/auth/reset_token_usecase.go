package authUC

import (
	"cms-server/domain/entity"
	"cms-server/domain/repository"
	"cms-server/domain/service/argon"
	"cms-server/domain/service/cache"
	serviceError "cms-server/domain/service/error"
	serviceJwt "cms-server/domain/service/jwt"
)

type ResetPasswordByTokenUsecase interface {
	VerifySession(token string) (string, error)
	ResetPass(IdUser, Password, NewPassword string) error
}

type ResetPasswordByTokenUsecaseImpl struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	cache       cache.RedisConfigImpl
	jwt         serviceJwt.JwtService
	argon       argon.Argon
}

func NewResetPasswordTokenUsecase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	cache cache.RedisConfigImpl,
	jwt serviceJwt.JwtService,
	argon argon.Argon,
) ResetPasswordByTokenUsecase {
	return &ResetPasswordByTokenUsecaseImpl{
		userRepo,
		sessionRepo,
		cache,
		jwt,
		argon,
	}
}

func (uc *ResetPasswordByTokenUsecaseImpl) VerifySession(token string) (string, error) {
	if _, err := uc.cache.Get(token); err != nil {
		if isExist := uc.sessionRepo.TokenExists(token); !isExist {
			return "", serviceError.NewErrorApp("Phiên làm việc không hợp lệ hoặc đã hết hạn")
		}
	}
	go func() {
		uc.sessionRepo.DeleteSessionForgotByToken(token)
		uc.cache.Delete(token)
	}()
	claim, err := uc.jwt.VerifyForgotPasswordToken(token)
	if err != nil {
		return "", err
	}

	return claim.Id, nil
}

func (uc *ResetPasswordByTokenUsecaseImpl) ResetPass(IdUser, Password, ConfirmPassword string) error {
	user, err := uc.userRepo.GetUserByID(IdUser)
	if err != nil {
		return serviceError.NewErrorApp("Không tìm thấy người dùng")
	}

	ConfirmPassword, err = uc.argon.HashPassword(ConfirmPassword)
	if err != nil {
		return serviceError.NewErrorApp("Không thể mã hóa mật khẩu")
	}

	if _, err = uc.userRepo.UpdateUser(user.ID, entity.User{Password: ConfirmPassword}); err != nil {
		return serviceError.NewErrorApp("Không thể cập nhật mật khẩu")
	}
	return nil
}
