package repository

import (
	"context"
	"user-service/domain/entity"
)

type RoleRepository interface {
	CreateRole(role entity.Role) error
	GetRoleByID(id string) (entity.Role, error)
	GetRoleByName(name string) (entity.Role, error)
	GetAllRoles() ([]entity.Role, error)
	UpdateRole(id string, role entity.Role) (entity.Role, error)
	DeleteByID(id string) error
	CheckRoleExist(name string) (bool, error)
	Tx(ctx context.Context) RoleRepository
}
