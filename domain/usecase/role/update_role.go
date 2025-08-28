package roleusecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"time"
)

type UpdateRoleUsecase interface {
	UpdateRole(id string, role entity.Role) (entity.Role, error)
}

type updateRoleUsecase struct {
	roleRepo repository.RoleRepository
}

func NewUpdateRoleUsecase(roleRepo repository.RoleRepository) UpdateRoleUsecase {
	return &updateRoleUsecase{
		roleRepo: roleRepo,
	}
}

func (u *updateRoleUsecase) UpdateRole(id string, role entity.Role) (entity.Role, error) {
	now := time.Now()
	role.UpdatedAt = &now

	return u.roleRepo.UpdateRole(id, role)
}
