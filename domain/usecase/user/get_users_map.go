package user

import (
	"user-service/domain/entity"
	"user-service/domain/repository"
)

type GetUserMapUsecase interface {
	Excute(userIds []string) (map[string]entity.User, error)
}

type getUserMapUsecase struct {
	userRepo repository.UserRepository
}

func NewGetUserMapUsecase(userRepo repository.UserRepository) GetUserMapUsecase {
	return &getUserMapUsecase{
		userRepo: userRepo,
	}
}

func (g *getUserMapUsecase) Excute(userIds []string) (map[string]entity.User, error) {
	users, err := g.userRepo.GetUserMap(userIds)
	if err != nil {
		return nil, err
	}
	return users, nil
}
