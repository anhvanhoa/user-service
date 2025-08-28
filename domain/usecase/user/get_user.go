package userusecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
)

type GetUserUsecase interface {
	GetUserByID(id string) (entity.User, error)
	GetUserByEmail(email string) (entity.User, error)
	GetUserByEmailOrPhone(val string) (entity.User, error)
}

type getUserUsecase struct {
	userRepo repository.UserRepository
}

func NewGetUserUsecase(userRepo repository.UserRepository) GetUserUsecase {
	return &getUserUsecase{
		userRepo: userRepo,
	}
}

func (g *getUserUsecase) GetUserByID(id string) (entity.User, error) {
	return g.userRepo.GetUserByID(id)
}

func (g *getUserUsecase) GetUserByEmail(email string) (entity.User, error) {
	return g.userRepo.GetUserByEmail(email)
}

func (g *getUserUsecase) GetUserByEmailOrPhone(val string) (entity.User, error) {
	return g.userRepo.GetUserByEmailOrPhone(val)
}
