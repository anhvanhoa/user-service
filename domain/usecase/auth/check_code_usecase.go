package authUC

import (
	"cms-server/domain/repository"
	serviceError "cms-server/domain/service/error"
)

type CheckCodeUsecase interface {
	CheckCode(code, email string) (bool, error)
}

type checkCodeUsecaseImpl struct {
	userRepo repository.UserRepository
	session  repository.SessionRepository
}

func NewCheckCodeUsecase(userRepo repository.UserRepository, session repository.SessionRepository) CheckCodeUsecase {
	return &checkCodeUsecaseImpl{
		userRepo: userRepo,
		session:  session,
	}
}

func (c *checkCodeUsecaseImpl) CheckCode(code, email string) (bool, error) {
	user, err := c.userRepo.GetUserByEmail(email)
	if err != nil {
		return false, serviceError.NewErrorApp("Không tìm thấy người dùng với email này")
	}

	session, err := c.session.GetSessionForgotAliveByTokenAndIdUser(code, user.ID)
	if err != nil {
		return false, serviceError.NewErrorApp("Mã xác thực không hợp lệ")
	}
	if session.Token == "" {
		return false, serviceError.NewErrorApp("Mã xác thực không hợp lệ hoặc đã hết hạn")
	}
	return true, nil
}
