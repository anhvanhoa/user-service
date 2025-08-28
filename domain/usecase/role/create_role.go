package roleusecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
)

type CreateRoleUsecase interface {
	CreateRole(role entity.Role) error
}

type createRoleUsecase struct {
	roleRepo repository.RoleRepository
}

func NewCreateRoleUsecase(roleRepo repository.RoleRepository) CreateRoleUsecase {
	return &createRoleUsecase{
		roleRepo: roleRepo,
	}
}

func (c *createRoleUsecase) CreateRole(role entity.Role) error {
	return c.roleRepo.CreateRole(role)
}
