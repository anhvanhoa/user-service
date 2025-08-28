package userusecase

import (
	"user-service/domain/repository"
)

type UpdateUserRolesUsecase interface {
	UpdateUserRoles(userID string, roleIDs []string) error
}

type updateUserRolesUsecase struct {
	userRoleRepo repository.UserRoleRepository
}

func NewUpdateUserRolesUsecase(userRoleRepo repository.UserRoleRepository) UpdateUserRolesUsecase {
	return &updateUserRolesUsecase{
		userRoleRepo: userRoleRepo,
	}
}

func (u *updateUserRolesUsecase) UpdateUserRoles(userID string, roleIDs []string) error {
	return u.userRoleRepo.UpdateUserRoles(userID, roleIDs)
}
