package userusecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"time"
)

type UpdateUserUsecase interface {
	UpdateUser(id string, data entity.User) (entity.UserInfor, error)
	UpdateUserByEmail(email string, data entity.User) (bool, error)
	UpdateUserWithRoles(id string, data entity.User, roleIDs []string) (entity.UserInfor, error)
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

func (u *updateUserUsecase) UpdateUser(id string, data entity.User) (entity.UserInfor, error) {
	now := time.Now()
	data.UpdatedAt = &now

	return u.userRepo.UpdateUser(id, data)
}

func (u *updateUserUsecase) UpdateUserByEmail(email string, data entity.User) (bool, error) {
	now := time.Now()
	data.UpdatedAt = &now

	return u.userRepo.UpdateUserByEmail(email, data)
}

func (u *updateUserUsecase) UpdateUserWithRoles(id string, data entity.User, roleIDs []string) (entity.UserInfor, error) {
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
