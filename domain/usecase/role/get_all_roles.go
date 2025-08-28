package roleusecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
)

type GetAllRolesUsecase interface {
	GetAllRoles() ([]entity.Role, error)
}

type getAllRolesUsecase struct {
	roleRepo repository.RoleRepository
}

func NewGetAllRolesUsecase(roleRepo repository.RoleRepository) GetAllRolesUsecase {
	return &getAllRolesUsecase{
		roleRepo: roleRepo,
	}
}

func (g *getAllRolesUsecase) GetAllRoles() ([]entity.Role, error) {
	return g.roleRepo.GetAllRoles()
}
