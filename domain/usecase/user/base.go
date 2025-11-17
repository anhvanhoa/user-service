package user

import (
	"user-service/domain/entity"

	"github.com/anhvanhoa/service-core/common"
)

type UserUsecaseI interface {
	CreateUser(data *entity.User) (entity.User, error)
	GetUserById(id string) (entity.User, error)
	GetUsers(pagination *common.Pagination, filter *entity.FilterUser) (*common.PaginationResult[entity.User], error)
	DeleteUserById(id string) error
	UpdateUserById(id string, data *entity.User) (entity.UserInfor, error)
}

type userUsecase struct {
	createUserUsecase CreateUserUsecase
	deleteUserUsecase DeleteUserUsecase
	getUserUsecase    GetUserUsecase
	getUsersUsecase   GetUsersUsecase
	updateUserUsecase UpdateUserUsecase
}

func NewUserUsecase(
	createUserUsecase CreateUserUsecase,
	deleteUserUsecase DeleteUserUsecase,
	getUserUsecase GetUserUsecase,
	getUsersUsecase GetUsersUsecase,
	updateUserUsecase UpdateUserUsecase,
) UserUsecaseI {
	return &userUsecase{
		createUserUsecase: createUserUsecase,
		deleteUserUsecase: deleteUserUsecase,
		getUserUsecase:    getUserUsecase,
		getUsersUsecase:   getUsersUsecase,
		updateUserUsecase: updateUserUsecase,
	}
}

func (u *userUsecase) CreateUser(data *entity.User) (entity.User, error) {
	return u.createUserUsecase.Excute(data)
}

func (u *userUsecase) GetUserById(id string) (entity.User, error) {
	return u.getUserUsecase.Excute(id)
}

func (u *userUsecase) GetUsers(pagination *common.Pagination, filter *entity.FilterUser) (*common.PaginationResult[entity.User], error) {
	return u.getUsersUsecase.Excute(pagination, filter)
}

func (u *userUsecase) DeleteUserById(id string) error {
	return u.deleteUserUsecase.Excute(id)
}

func (u *userUsecase) UpdateUserById(id string, data *entity.User) (entity.UserInfor, error) {
	return u.updateUserUsecase.Excute(id, data)
}
