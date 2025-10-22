package user

import (
	"user-service/domain/entity"
	"user-service/domain/repository"
)

type UpdateUserUsecase interface {
	Excute(id string, data entity.User) (entity.UserInfor, error)
}

type updateUserUsecase struct {
	userRepo repository.UserRepository
}

func NewUpdateUserUsecase(userRepo repository.UserRepository) UpdateUserUsecase {
	return &updateUserUsecase{
		userRepo: userRepo,
	}
}

func (u *updateUserUsecase) Excute(id string, data entity.User) (entity.UserInfor, error) {
	userInfo, err := u.userRepo.UpdateUser(id, data)
	if err != nil {
		return entity.UserInfor{}, err
	}

	return userInfo, nil
}
