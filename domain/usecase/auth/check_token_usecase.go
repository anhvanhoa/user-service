package authUC

import (
	"cms-server/domain/entity"
	"cms-server/domain/repository"
	serviceError "cms-server/domain/service/error"
)

type CheckTokenUsecase interface {
	CheckToken(token string) (bool, error)
}

type checkTokenUsecaseImpl struct {
	sessionRepo repository.SessionRepository
}

func NewCheckTokenUsecase(sessionRepo repository.SessionRepository) CheckTokenUsecase {
	return &checkTokenUsecaseImpl{
		sessionRepo: sessionRepo,
	}
}

func (c *checkTokenUsecaseImpl) CheckToken(token string) (bool, error) {
	session, err := c.sessionRepo.GetSessionAliveByToken(entity.SessionTypeForgot, token)
	if err != nil {
		return false, serviceError.NewErrorApp("Phiên làm việc không hợp lệ hoặc đã hết hạn")
	}
	if session.Token == "" {
		return false, serviceError.NewErrorApp("Phiên làm việc không hợp lệ hoặc đã hết hạn")
	}
	return true, nil
}
