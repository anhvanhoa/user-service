package role

import (
	"user-service/domain/entity"
	"user-service/domain/repository"
)

type GetAllRolesUsecase interface {
	Excute() ([]entity.Role, error)
}

type getAllRolesUsecase struct {
	roleRepo repository.RoleRepository
}

func NewGetAllRolesUsecase(roleRepo repository.RoleRepository) GetAllRolesUsecase {
	return &getAllRolesUsecase{
		roleRepo: roleRepo,
	}
}

func (g *getAllRolesUsecase) Excute() ([]entity.Role, error) {
	return g.roleRepo.GetAllRoles()
}
