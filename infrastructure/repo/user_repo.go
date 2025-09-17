package repo

import (
	"context"
	"user-service/domain/entity"
	"user-service/domain/repository"

	"github.com/go-pg/pg/v10"
)

type userRepository struct {
	db pg.DBI
}

func NewUserRepository(db *pg.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) CreateUser(user entity.User) (entity.UserInfor, error) {
	_, err := ur.db.Model(&user).Insert()
	return user.GetInfor(), err
}

func (ur *userRepository) GetUserByEmailOrPhone(val string) (entity.User, error) {
	var user entity.User
	err := ur.db.Model(&user).Where("email = ?", val).WhereOr("phone = ?", val).Select()
	return user, err
}

func (ur *userRepository) CheckUserExist(val string) (bool, error) {
	var user entity.User
	count, err := ur.db.Model(&user).Where("email = ?", val).Count()
	isExist := count > 0
	return isExist, err
}

func (ur *userRepository) UpdateUser(id string, user entity.User) (entity.UserInfor, error) {
	_, err := ur.db.Model(&user).Where("id = ?", id).UpdateNotZero(&user)
	return user.GetInfor(), err
}

func (ur *userRepository) GetUserByID(id string) (entity.User, error) {
	var user entity.User
	err := ur.db.Model(&user).Where("id = ?", id).Select()
	return user, err
}

func (ur *userRepository) GetUserByEmail(email string) (entity.User, error) {
	var user entity.User
	err := ur.db.Model(&user).Where("email = ?", email).Select()
	return user, err
}

func (ur *userRepository) UpdateUserByEmail(email string, user entity.User) (bool, error) {
	r, err := ur.db.Model(&user).Where("email = ?", email).UpdateNotZero(&user)
	return r.RowsAffected() != -1, err
}

func (ur *userRepository) DeleteByID(id string) error {
	var user entity.User
	_, err := ur.db.Model(&user).Where("id = ?", id).Delete()
	return err
}

func (ur *userRepository) Tx(ctx context.Context) repository.UserRepository {
	tx := getTx(ctx, ur.db)
	return &userRepository{
		db: tx,
	}
}
