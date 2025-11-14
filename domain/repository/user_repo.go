package repository

import (
	"context"
	"user-service/domain/entity"

	"github.com/anhvanhoa/service-core/common"
)

type UserRepository interface {
	CreateUser(entity.User) (entity.UserInfor, error)
	GetUserByEmailOrPhone(val string) (entity.User, error)
	GetUserByID(id string) (entity.User, error)
	CheckUserExist(val string) (bool, error)
	GetUserByEmail(email string) (entity.User, error)
	GetUsers(pagination *common.Pagination, filter *entity.FilterUser) ([]entity.User, int, error)
	UpdateUser(Id string, data entity.User) (entity.UserInfor, error)
	UpdateUserByEmail(email string, data entity.User) (bool, error)
	DeleteByID(id string) error
	Tx(ctx context.Context) UserRepository
}
