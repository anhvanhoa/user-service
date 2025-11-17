package user

import (
	"user-service/domain/entity"
	"user-service/domain/repository"
)

type LockUserUsecase interface {
	Excute(id string, reason string, by string) error
}

type lockUserUsecase struct {
	userRepo repository.UserRepository
}

func NewLockUserUsecase(userRepo repository.UserRepository) LockUserUsecase {
	return &lockUserUsecase{
		userRepo: userRepo,
	}

}

func (l *lockUserUsecase) Excute(id string, reason string, by string) error {
	user, err := l.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user.LockedAt != nil || user.Status == entity.UserStatusLocked {
		return ErrUserAlreadyLocked
	}
	return l.userRepo.LockUser(id, reason, by)
}
