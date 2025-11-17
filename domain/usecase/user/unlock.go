package user

import (
	"context"
	"user-service/domain/entity"
	"user-service/domain/repository"
)

type UnlockUserUsecase interface {
	Excute(ctx context.Context, id string) error
}

type unlockUserUsecase struct {
	userRepo repository.UserRepository
}

func NewUnlockUserUsecase(userRepo repository.UserRepository) UnlockUserUsecase {
	return &unlockUserUsecase{
		userRepo: userRepo,
	}
}

func (u *unlockUserUsecase) Excute(ctx context.Context, id string) error {
	user, err := u.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user.LockedAt == nil || user.Status != entity.UserStatusLocked {
		return ErrUserAlreadyUnlocked
	}
	return u.userRepo.UnlockUser(id)
}
