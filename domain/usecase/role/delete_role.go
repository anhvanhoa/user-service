package roleusecase

import (
	"context"
	"user-service/domain/repository"
)

type DeleteRoleUsecase interface {
	Excute(ctx context.Context, id string) error
}

type deleteRoleUsecase struct {
	roleRepo repository.RoleRepository
}

func NewDeleteRoleUsecase(roleRepo repository.RoleRepository) DeleteRoleUsecase {
	return &deleteRoleUsecase{
		roleRepo: roleRepo,
	}
}

func (d *deleteRoleUsecase) Excute(ctx context.Context, id string) error {
	return d.roleRepo.DeleteByID(ctx, id)
}
