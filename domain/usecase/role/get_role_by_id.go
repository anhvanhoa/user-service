package role

import (
	"user-service/domain/entity"
	"user-service/domain/repository"
)

type GetRoleByIDUsecase interface {
	Excute(id string) (entity.Role, error)
}

type getRoleByIDUsecase struct {
	roleRepo repository.RoleRepository
}

func NewGetRoleByIDUsecase(roleRepo repository.RoleRepository) GetRoleByIDUsecase {
	return &getRoleByIDUsecase{
		roleRepo: roleRepo,
	}
}

func (g *getRoleByIDUsecase) Excute(id string) (entity.Role, error) {
	return g.roleRepo.GetRoleByID(id)
}
