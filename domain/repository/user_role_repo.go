package repository

import (
	"context"
	"user-service/domain/entity"
)

type UserRoleRepository interface {
	CreateUserRole(userRole entity.UserRole) error
	GetUserRolesByUserID(userID string) ([]entity.UserRole, error)
	GetUserRolesByRoleID(roleID string) ([]entity.UserRole, error)
	UpdateUserRoles(userID string, roleIDs []string) error
	DeleteUserRolesByUserID(ctx context.Context, userID string) error
	DeleteUserRole(ctx context.Context, userID, roleID string) error
	Tx(ctx context.Context) UserRoleRepository
}
