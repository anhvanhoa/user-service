package user

import (
	"context"
	"user-service/domain/entity"
	"user-service/domain/repository"
)

type LockUserUsecase interface {
	Excute(ctx context.Context, id string, reason string, by string) error
}

type lockUserUsecase struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewLockUserUsecase(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) LockUserUsecase {
	return &lockUserUsecase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}

}

func (l *lockUserUsecase) Excute(ctx context.Context, id string, reason string, by string) error {
	user, err := l.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user.LockedAt != nil || user.Status == entity.UserStatusLocked {
		return ErrUserAlreadyLocked
	}
	err = l.userRepo.LockUser(id, reason, by)
	if err != nil {
		return err
	}
	return l.sessionRepo.DeleteAllSessionsByUserID(ctx, id)
}
