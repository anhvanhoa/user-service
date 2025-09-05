package role

import (
	"user-service/domain/repository"
)

type DeleteRoleUsecase interface {
	Excute(id string) error
}

type deleteRoleUsecase struct {
	roleRepo repository.RoleRepository
}

func NewDeleteRoleUsecase(roleRepo repository.RoleRepository) DeleteRoleUsecase {
	return &deleteRoleUsecase{
		roleRepo: roleRepo,
	}
}

func (d *deleteRoleUsecase) Excute(id string) error {
	return d.roleRepo.DeleteByID(id)
}
