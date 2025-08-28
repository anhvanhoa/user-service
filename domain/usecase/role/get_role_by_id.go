package roleusecase

import (
	"user-service/domain/entity"
	"user-service/domain/repository"
)

type GetRoleByIDUsecase interface {
	GetRoleByID(id string) (entity.Role, error)
}

type getRoleByIDUsecase struct {
	roleRepo repository.RoleRepository
}

func NewGetRoleByIDUsecase(roleRepo repository.RoleRepository) GetRoleByIDUsecase {
	return &getRoleByIDUsecase{
		roleRepo: roleRepo,
	}
}

func (g *getRoleByIDUsecase) GetRoleByID(id string) (entity.Role, error) {
	return g.roleRepo.GetRoleByID(id)
}
