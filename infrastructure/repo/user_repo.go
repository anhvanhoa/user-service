package repo

import (
	"context"
	"user-service/constants"
	"user-service/domain/entity"
	"user-service/domain/repository"

	"github.com/anhvanhoa/service-core/common"
	"github.com/anhvanhoa/service-core/utils"
	"github.com/go-pg/pg/v10"
)

type userRepository struct {
	db     pg.DBI
	helper utils.Helper
}

func NewUserRepository(db *pg.DB, helper utils.Helper) repository.UserRepository {
	return &userRepository{
		db:     db,
		helper: helper,
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

func (ur *userRepository) GetUsers(pagination *common.Pagination, filter *entity.FilterUser) ([]entity.User, int, error) {
	var users []entity.User
	query := ur.db.Model(&users)
	if filter != nil && filter.Status != nil {
		column := constants.MapStatusColumn[*filter.Status]
		if filter.FromDate != nil {
			query = query.Where(column+" >= ?", filter.FromDate)
		}
		if filter.ToDate != nil {
			query = query.Where(column+" <= ?", filter.ToDate)
		}
		query = query.Where("status = ?", filter.Status)
	}
	if pagination != nil {
		if pagination.Search != "" {
			query = query.Where("full_name ILIKE ? OR email ILIKE ? OR phone ILIKE ? OR address ILIKE ? OR bio ILIKE ?", "%"+pagination.Search+"%", "%"+pagination.Search+"%", "%"+pagination.Search+"%", "%"+pagination.Search+"%", "%"+pagination.Search+"%")
		}
		if pagination.SortBy != "" {
			query = query.Order(pagination.SortBy)
		}
		if pagination.SortOrder != "" {
			query = query.Order(pagination.SortOrder)
		}
	}
	total, err := query.Count()
	if err != nil {
		return nil, 0, err
	}

	if pagination != nil {
		if pagination.Page > 0 {
			offset := ur.helper.CalculateOffset(pagination.Page, pagination.PageSize)
			query = query.Offset(offset)
		}
		if pagination.PageSize > 0 {
			query = query.Limit(pagination.PageSize)
		}
	}

	err = query.Select()
	if err != nil {
		return nil, 0, err
	}
	return users, total, err
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
