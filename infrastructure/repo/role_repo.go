package repo

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-pg/pg/v10"
)

type roleRepository struct {
	db pg.DBI
}

func NewRoleRepository(db *pg.DB) repository.RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (rr *roleRepository) CreateRole(role entity.Role) error {
	_, err := rr.db.Model(&role).Insert()
	return err
}

func (rr *roleRepository) GetRoleByID(id string) (entity.Role, error) {
	var role entity.Role
	err := rr.db.Model(&role).Where("id = ?", id).Select()
	return role, err
}

func (rr *roleRepository) GetRoleByName(name string) (entity.Role, error) {
	var role entity.Role
	err := rr.db.Model(&role).Where("name = ?", name).Select()
	return role, err
}

func (rr *roleRepository) GetAllRoles() ([]entity.Role, error) {
	var roles []entity.Role
	err := rr.db.Model(&roles).Select()
	return roles, err
}

func (rr *roleRepository) UpdateRole(id string, role entity.Role) (entity.Role, error) {
	var setClauses []string
	var params []any

	v := reflect.ValueOf(role)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if field.Name == "ID" || field.Name == "CreatedAt" {
			continue
		}

		if !value.IsZero() {
			columnName := field.Tag.Get("pg")
			setClauses = append(setClauses, fmt.Sprintf("%s = ?", columnName))
			params = append(params, value.Interface())
		}
	}

	if len(setClauses) == 0 {
		return role, nil
	}

	setQuery := strings.Join(setClauses, ", ")

	if _, err := rr.db.Model(&role).Where("id = ?", id).Set(setQuery, params...).Update(); err != nil {
		return entity.Role{}, err
	}

	return role, nil
}

func (rr *roleRepository) DeleteByID(ctx context.Context, id string) error {
	var role entity.Role
	_, err := rr.db.ModelContext(ctx, &role).Where("id = ?", id).Delete()
	return err
}

func (rr *roleRepository) CheckRoleExist(name string) (bool, error) {
	var role entity.Role
	count, err := rr.db.Model(&role).Where("name = ?", name).Count()
	isExist := count > 0
	return isExist, err
}

func (rr *roleRepository) Tx(ctx context.Context) repository.RoleRepository {
	tx := getTx(ctx, rr.db)
	return &roleRepository{
		db: tx,
	}
}
