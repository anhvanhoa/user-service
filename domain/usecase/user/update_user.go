package user

import (
	"time"
	"user-service/domain/entity"
	"user-service/domain/repository"
)

type UpdateUserUsecase interface {
	Excute(id string, data entity.User, roleIDs []string) (entity.UserInfor, error)
}

type updateUserUsecase struct {
	userRepo     repository.UserRepository
	userRoleRepo repository.UserRoleRepository
}

func NewUpdateUserUsecase(userRepo repository.UserRepository, userRoleRepo repository.UserRoleRepository) UpdateUserUsecase {
	return &updateUserUsecase{
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
	}
}

func (u *updateUserUsecase) Excute(id string, data entity.User, roleIDs []string) (entity.UserInfor, error) {
	now := time.Now()
	data.UpdatedAt = &now

	userInfo, err := u.userRepo.UpdateUser(id, data)
	if err != nil {
		return entity.UserInfor{}, err
	}

	if len(roleIDs) > 0 {
		err = u.userRoleRepo.UpdateUserRoles(id, roleIDs)
		if err != nil {
			return entity.UserInfor{}, err
		}
	}

	return userInfo, nil
}
