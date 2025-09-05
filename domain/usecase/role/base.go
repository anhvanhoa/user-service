package roleusecase

import (
	"context"
	"user-service/domain/entity"
)

type RoleUsecaseI interface {
	GetAllRoles() ([]entity.Role, error)
	GetRoleById(id string) (entity.Role, error)
	CreateRole(role entity.Role) error
	UpdateRole(id string, role entity.Role) (entity.Role, error)
	DeleteRole(ctx context.Context, id string) error
	CheckRole(name string) (bool, error)
}

type roleUsecase struct {
	getAllRolesUsecase GetAllRolesUsecase
	getRoleByIdUsecase GetRoleByIDUsecase
	createRoleUsecase  CreateRoleUsecase
	updateRoleUsecase  UpdateRoleUsecase
	deleteRoleUsecase  DeleteRoleUsecase
	checkRoleUsecase   CheckRoleUsecase
}

func NewRoleUsecase(
	getAllRolesUsecase GetAllRolesUsecase,
	getRoleByIdUsecase GetRoleByIDUsecase,
	createRoleUsecase CreateRoleUsecase,
	updateRoleUsecase UpdateRoleUsecase,
	deleteRoleUsecase DeleteRoleUsecase,
	checkRoleUsecase CheckRoleUsecase,
) RoleUsecaseI {
	return &roleUsecase{
		getAllRolesUsecase: getAllRolesUsecase,
		getRoleByIdUsecase: getRoleByIdUsecase,
		createRoleUsecase:  createRoleUsecase,
		updateRoleUsecase:  updateRoleUsecase,
		deleteRoleUsecase:  deleteRoleUsecase,
		checkRoleUsecase:   checkRoleUsecase,
	}
}

func (r *roleUsecase) GetAllRoles() ([]entity.Role, error) {
	return r.getAllRolesUsecase.Excute()
}

func (r *roleUsecase) GetRoleById(id string) (entity.Role, error) {
	return r.getRoleByIdUsecase.Excute(id)
}

func (r *roleUsecase) CreateRole(role entity.Role) error {
	return r.createRoleUsecase.Excute(role)
}

func (r *roleUsecase) UpdateRole(id string, role entity.Role) (entity.Role, error) {
	return r.updateRoleUsecase.Excute(id, role)
}

func (r *roleUsecase) DeleteRole(ctx context.Context, id string) error {
	return r.deleteRoleUsecase.Excute(ctx, id)
}

func (r *roleUsecase) CheckRole(name string) (bool, error) {
	return r.checkRoleUsecase.Excute(name)
}
