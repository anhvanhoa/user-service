package user

import (
	"user-service/domain/repository"
)

type DeleteUserUsecase interface {
	Excute(id string) error
}

type deleteUserUsecase struct {
	userRepo repository.UserRepository
}

func NewDeleteUserUsecase(userRepo repository.UserRepository) DeleteUserUsecase {
	return &deleteUserUsecase{
		userRepo: userRepo,
	}
}

func (d *deleteUserUsecase) Excute(id string) error {
	return d.userRepo.DeleteByID(id)
}
