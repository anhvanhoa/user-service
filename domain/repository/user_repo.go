package repository

import (
	"cms-server/domain/entity"
	"context"
)

type UserRepository interface {
	CreateUser(entity.User) (entity.UserInfor, error)
	GetUserByEmailOrPhone(val string) (entity.User, error)
	GetUserByID(id string) (entity.User, error)
	CheckUserExist(val string) (bool, error)
	GetUserByEmail(email string) (entity.User, error)
	UpdateUser(Id string, data entity.User) (entity.UserInfor, error)
	UpdateUserByEmail(email string, data entity.User) (bool, error)
	Tx(ctx context.Context) UserRepository
}
