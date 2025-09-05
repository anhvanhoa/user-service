package user

import (
	"user-service/domain/entity"
	"user-service/domain/repository"
)

type GetUserUsecase interface {
	Excute(id string) (entity.User, error)
}

type getUserUsecase struct {
	userRepo repository.UserRepository
}

func NewGetUserUsecase(userRepo repository.UserRepository) GetUserUsecase {
	return &getUserUsecase{
		userRepo: userRepo,
	}
}

func (g *getUserUsecase) Excute(id string) (entity.User, error) {
	return g.userRepo.GetUserByID(id)
}
