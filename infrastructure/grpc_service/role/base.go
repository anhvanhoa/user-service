package role_server

import (
	"context"
	"user-service/domain/entity"
	role "user-service/domain/usecase/role"
	"user-service/infrastructure/repo"

	proto_role "github.com/anhvanhoa/sf-proto/gen/role/v1"
	"github.com/go-pg/pg/v10"
)

type roleServer struct {
	proto_role.UnimplementedRoleServiceServer
	roleUsecase role.RoleUsecaseI
}

func NewRoleServer(db *pg.DB) proto_role.RoleServiceServer {
	roleRepo := repo.NewRoleRepository(db)
	roleUC := role.NewRoleUsecase(roleRepo)
	return &roleServer{
		roleUsecase: roleUC,
	}
}

func (s *roleServer) GetAllRoles(ctx context.Context, req *proto_role.GetAllRolesRequest) (*proto_role.GetAllRolesResponse, error) {
	roles, err := s.roleUsecase.GetAllRoles()
	if err != nil {
		return nil, err
	}
	return &proto_role.GetAllRolesResponse{
		Roles: s.createProtoRoles(roles),
	}, nil
}

func (s *roleServer) GetRoleById(ctx context.Context, req *proto_role.GetRoleByIdRequest) (*proto_role.GetRoleByIdResponse, error) {
	role, err := s.roleUsecase.GetRoleById(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto_role.GetRoleByIdResponse{
		Role: s.createProtoRole(role),
	}, nil
}

func (s *roleServer) CreateRole(ctx context.Context, req *proto_role.CreateRoleRequest) (*proto_role.CreateRoleResponse, error) {
	role := s.createEntityRole(req.Name, req.Description, req.Variant)
	err := s.roleUsecase.CreateRole(role)
	if err != nil {
		return nil, err
	}
	return &proto_role.CreateRoleResponse{
		Message: "Create role successfully",
		Success: true,
	}, nil
}

func (s *roleServer) UpdateRole(ctx context.Context, req *proto_role.UpdateRoleRequest) (*proto_role.UpdateRoleResponse, error) {
	role := s.createEntityRole(req.Name, req.Description, req.Variant)
	updatedRole, err := s.roleUsecase.UpdateRole(req.Id, role)
	if err != nil {
		return nil, err
	}
	return &proto_role.UpdateRoleResponse{
		Role: s.createProtoRole(updatedRole),
	}, nil
}

func (s *roleServer) DeleteRole(ctx context.Context, req *proto_role.DeleteRoleRequest) (*proto_role.DeleteRoleResponse, error) {

	err := s.roleUsecase.DeleteRole(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto_role.DeleteRoleResponse{
		Message: "Delete role successfully",
		Success: true,
	}, nil
}

func (s *roleServer) createEntityRole(
	name string,
	description string,
	variant string,
) entity.Role {
	return entity.Role{
		Name:        name,
		Description: description,
		Variant:     variant,
	}
}

func (s *roleServer) createProtoRoles(roles []entity.Role) []*proto_role.Role {
	protoRoles := make([]*proto_role.Role, len(roles))
	for i, role := range roles {
		protoRoles[i] = s.createProtoRole(role)
	}
	return protoRoles
}

func (s *roleServer) createProtoRole(role entity.Role) *proto_role.Role {
	return &proto_role.Role{
		Id:   role.ID,
		Name: role.Name,
	}
}
