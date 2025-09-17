package user

import (
	"user-service/domain/entity"
)

type UserUsecaseI interface {
	GetUserById(id string) (entity.User, error)
	DeleteUserById(id string) error
	UpdateUserById(id string, data entity.User, roleIDs []string) (entity.UserInfor, error)
	UpdateUserRolesById(id string, roleIDs []string) error
}

type userUsecase struct {
	deleteUserUsecase      DeleteUserUsecase
	getUserUsecase         GetUserUsecase
	updateUserUsecase      UpdateUserUsecase
	updateUserRolesUsecase UpdateUserRolesUsecase
}

func NewUserUsecase(
	deleteUserUsecase DeleteUserUsecase,
	getUserUsecase GetUserUsecase,
	updateUserUsecase UpdateUserUsecase,
	updateUserRolesUsecase UpdateUserRolesUsecase,
) UserUsecaseI {
	return &userUsecase{
		deleteUserUsecase:      deleteUserUsecase,
		getUserUsecase:         getUserUsecase,
		updateUserUsecase:      updateUserUsecase,
		updateUserRolesUsecase: updateUserRolesUsecase,
	}
}

func (u *userUsecase) GetUserById(id string) (entity.User, error) {
	return u.getUserUsecase.Excute(id)
}

func (u *userUsecase) DeleteUserById(id string) error {
	return u.deleteUserUsecase.Excute(id)
}

func (u *userUsecase) UpdateUserById(id string, data entity.User, roleIDs []string) (entity.UserInfor, error) {
	return u.updateUserUsecase.Excute(id, data, roleIDs)
}

func (u *userUsecase) UpdateUserRolesById(id string, roleIDs []string) error {
	return u.updateUserRolesUsecase.Excute(id, roleIDs)
}
