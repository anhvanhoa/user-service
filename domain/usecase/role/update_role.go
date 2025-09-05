package roleusecase

import (
	"time"
	"user-service/domain/entity"
	"user-service/domain/repository"
)

type UpdateRoleUsecase interface {
	Excute(id string, role entity.Role) (entity.Role, error)
}

type updateRoleUsecase struct {
	roleRepo repository.RoleRepository
}

func NewUpdateRoleUsecase(roleRepo repository.RoleRepository) UpdateRoleUsecase {
	return &updateRoleUsecase{
		roleRepo: roleRepo,
	}
}

func (u *updateRoleUsecase) Excute(id string, role entity.Role) (entity.Role, error) {
	now := time.Now()
	role.UpdatedAt = &now

	return u.roleRepo.UpdateRole(id, role)
}
