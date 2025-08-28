package roleusecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
)

type GetRoleByNameUsecase interface {
	GetRoleByName(name string) (entity.Role, error)
}

type getRoleByNameUsecase struct {
	roleRepo repository.RoleRepository
}

func NewGetRoleByNameUsecase(roleRepo repository.RoleRepository) GetRoleByNameUsecase {
	return &getRoleByNameUsecase{
		roleRepo: roleRepo,
	}
}

func (g *getRoleByNameUsecase) GetRoleByName(name string) (entity.Role, error) {
	return g.roleRepo.GetRoleByName(name)
}
