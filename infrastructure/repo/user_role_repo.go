package repo

import (
	"context"
	"time"
	"user-service/domain/entity"
	"user-service/domain/repository"

	"github.com/go-pg/pg/v10"
)

type userRoleRepositoryImpl struct {
	db pg.DBI
}

func NewUserRoleRepository(db *pg.DB) repository.UserRoleRepository {
	return &userRoleRepositoryImpl{
		db: db,
	}
}

func (urr *userRoleRepositoryImpl) CreateUserRole(userRole entity.UserRole) error {
	userRole.CreatedAt = time.Now()
	_, err := urr.db.Model(&userRole).Insert()
	if err != nil {
		return err
	}
	return nil
}

func (urr *userRoleRepositoryImpl) GetUserRolesByUserID(userID string) ([]entity.UserRole, error) {
	var userRoles []entity.UserRole
	err := urr.db.Model(&userRoles).Where("user_id = ?", userID).Select()
	if err != nil {
		return nil, err
	}
	return userRoles, nil
}

func (urr *userRoleRepositoryImpl) GetUserRolesByRoleID(roleID string) ([]entity.UserRole, error) {
	var userRoles []entity.UserRole
	err := urr.db.Model(&userRoles).Where("role_id = ?", roleID).Select()
	if err != nil {
		return nil, err
	}
	return userRoles, nil
}

func (urr *userRoleRepositoryImpl) UpdateUserRoles(userID string, roleIDs []string) error {
	// Use background context for this operation
	ctx := context.Background()

	// Delete existing user roles
	err := urr.DeleteUserRolesByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Create new user roles
	for _, roleID := range roleIDs {
		userRole := entity.UserRole{
			UserID:    userID,
			RoleID:    roleID,
			CreatedAt: time.Now(),
		}
		err = urr.CreateUserRole(userRole)
		if err != nil {
			return err
		}
	}
	return nil
}

func (urr *userRoleRepositoryImpl) DeleteUserRolesByUserID(ctx context.Context, userID string) error {
	_, err := urr.db.ModelContext(ctx, &entity.UserRole{}).Where("user_id = ?", userID).Delete()
	if err != nil {
		return err
	}
	return nil
}

func (urr *userRoleRepositoryImpl) DeleteUserRole(ctx context.Context, userID, roleID string) error {
	_, err := urr.db.ModelContext(ctx, &entity.UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete()
	if err != nil {
		return err
	}
	return nil
}

func (urr *userRoleRepositoryImpl) Tx(ctx context.Context) repository.UserRoleRepository {
	tx := getTx(ctx, urr.db)
	return &userRoleRepositoryImpl{
		db: tx,
	}
}
