package user

import (
	"user-service/domain/repository"
)

type UpdateUserRolesUsecase interface {
	Excute(userID string, roleIDs []string) error
}

type updateUserRolesUsecase struct {
	userRoleRepo repository.UserRoleRepository
}

func NewUpdateUserRolesUsecase(userRoleRepo repository.UserRoleRepository) UpdateUserRolesUsecase {
	return &updateUserRolesUsecase{
		userRoleRepo: userRoleRepo,
	}
}

func (u *updateUserRolesUsecase) Excute(userID string, roleIDs []string) error {
	return u.userRoleRepo.UpdateUserRoles(userID, roleIDs)
}
