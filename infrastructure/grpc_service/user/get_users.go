package user_server

import (
	"context"
	"user-service/domain/entity"

	"github.com/anhvanhoa/service-core/common"
	common_proto "github.com/anhvanhoa/sf-proto/gen/common/v1"
	proto_role "github.com/anhvanhoa/sf-proto/gen/role/v1"
	proto_user "github.com/anhvanhoa/sf-proto/gen/user/v1"
	proto_user_role "github.com/anhvanhoa/sf-proto/gen/user_role/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *userServer) GetUsers(ctx context.Context, req *proto_user.GetUsersRequest) (*proto_user.GetUsersResponse, error) {
	pagination := s.convertPagination(req.Pagination)
	filter := s.convertFilter(req.Filter)

	// 1️⃣ Lấy danh sách user
	result, err := s.userUsecase.GetUsersUsecase.Excute(pagination, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 2️⃣ Lấy userIds
	userIds, err := s.userUsecase.GetUsersUsecase.ExtractUserIds(result.Data)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 3️⃣ Lấy roles cho tất cả user
	userRoles := &proto_user_role.GetUserRolesResponse{
		UserRoleMap: make(map[string]*proto_user_role.RoleList),
	}
	if s.permissionClient.UserRoleServiceClient != nil {
		userRoles, err = s.permissionClient.UserRoleServiceClient.GetUserRoles(ctx, &proto_user_role.GetUserRolesRequest{
			UserIds: userIds,
		})
		if err != nil {
			s.log.Error("Error getting user roles: " + err.Error())
		}
	}

	// 4️⃣ Tạo proto user kết hợp role
	return &proto_user.GetUsersResponse{
		Users:      s.createProtoUsers(result.Data, userRoles.UserRoleMap),
		Pagination: s.convertPaginationResponse(result),
	}, nil
}

func (s *userServer) convertPagination(pagination *common_proto.PaginationRequest) *common.Pagination {
	if pagination == nil {
		return nil
	}
	return &common.Pagination{
		Page:      int(pagination.Page),
		PageSize:  int(pagination.PageSize),
		Search:    pagination.Search,
		SortBy:    pagination.SortBy,
		SortOrder: pagination.SortOrder,
	}
}

func (s *userServer) convertFilter(filter *proto_user.UserFilter) *entity.FilterUser {
	if filter == nil {
		return nil
	}
	var filterResult entity.FilterUser
	switch filter.Status {
	case proto_user.UserStatus_active:
		active := entity.UserStatus(string(entity.UserStatusActive))
		filterResult.Status = &active
	case proto_user.UserStatus_inactive:
		inactive := entity.UserStatus(string(entity.UserStatusInactive))
		filterResult.Status = &inactive
	case proto_user.UserStatus_locked:
		locked := entity.UserStatus(string(entity.UserStatusLocked))
		filterResult.Status = &locked
	default:
		filterResult.Status = nil
	}
	if filter.FromDate != nil {
		fromDate := filter.FromDate.AsTime()
		filterResult.FromDate = &fromDate
	}
	if filter.ToDate != nil {
		toDate := filter.ToDate.AsTime()
		filterResult.ToDate = &toDate
	}
	return &filterResult
}

func (s *userServer) createProtoUsers(users []entity.User, userRolesMap map[string]*proto_user_role.RoleList) []*proto_user.User {
	protoUsers := make([]*proto_user.User, len(users))

	for i, user := range users {
		u := s.createProtoUser(user)
		u.Roles = s.mapRoles(userRolesMap[user.ID])
		protoUsers[i] = u
	}

	return protoUsers
}

func (s *userServer) mapRoles(roles *proto_user_role.RoleList) []*proto_role.Role {
	if roles == nil || len(roles.Roles) == 0 {
		return nil
	}

	protoRoles := make([]*proto_role.Role, len(roles.Roles))
	for i, role := range roles.Roles {
		protoRoles[i] = &proto_role.Role{
			Id:      role.Id,
			Name:    role.Name,
			Variant: role.Variant,
		}
	}
	return protoRoles
}

func (s *userServer) convertPaginationResponse(result *common.PaginationResult[entity.User]) *common_proto.PaginationResponse {
	return &common_proto.PaginationResponse{
		Total:      int32(result.Total),
		TotalPages: int32(result.TotalPages),
		Page:       int32(result.Page),
		PageSize:   int32(result.PageSize),
	}
}
