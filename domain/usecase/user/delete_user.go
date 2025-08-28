package userusecase

import (
	"auth-service/domain/repository"
	"context"
)

type DeleteUserUsecase interface {
	DeleteUser(ctx context.Context, id string) error
}

type deleteUserUsecase struct {
	userRepo repository.UserRepository
}

func NewDeleteUserUsecase(userRepo repository.UserRepository) DeleteUserUsecase {
	return &deleteUserUsecase{
		userRepo: userRepo,
	}
}

func (d *deleteUserUsecase) DeleteUser(ctx context.Context, id string) error {
	return d.userRepo.DeleteByID(ctx, id)
}
