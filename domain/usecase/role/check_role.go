package roleusecase

import (
	"user-service/domain/repository"
)

type CheckRoleUsecase interface {
	CheckRoleExist(name string) (bool, error)
}

type checkRoleUsecase struct {
	roleRepo repository.RoleRepository
}

func NewCheckRoleUsecase(roleRepo repository.RoleRepository) CheckRoleUsecase {
	return &checkRoleUsecase{
		roleRepo: roleRepo,
	}
}

func (c *checkRoleUsecase) CheckRoleExist(name string) (bool, error) {
	return c.roleRepo.CheckRoleExist(name)
}
