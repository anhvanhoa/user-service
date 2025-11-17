package grpc_client

import (
	gc "github.com/anhvanhoa/service-core/domain/grpc_client"
	proto_permission "github.com/anhvanhoa/sf-proto/gen/permission/v1"
	proto_user_role "github.com/anhvanhoa/sf-proto/gen/user_role/v1"
)

type PermissionClientImpl struct {
	PermissionServiceClient proto_permission.PermissionServiceClient
	UserRoleServiceClient   proto_user_role.UserRoleServiceClient
}

func NewPermissionClient(client *gc.Client) *PermissionClientImpl {
	return &PermissionClientImpl{
		PermissionServiceClient: proto_permission.NewPermissionServiceClient(client.GetConnection()),
		UserRoleServiceClient:   proto_user_role.NewUserRoleServiceClient(client.GetConnection()),
	}
}
