package user_server

import (
	"context"
	"time"
	"user-service/domain/entity"

	"github.com/anhvanhoa/service-core/common"
	common_proto "github.com/anhvanhoa/sf-proto/gen/common/v1"
	proto_user "github.com/anhvanhoa/sf-proto/gen/user/v1"
)

func (s *userServer) GetUsers(ctx context.Context, req *proto_user.GetUsersRequest) (*proto_user.GetUsersResponse, error) {
	pagination := s.convertPagination(req.Pagination)
	filter := s.convertFilter(req.Filter)
	result, err := s.userUsecase.GetUsers(pagination, filter)
	if err != nil {
		return nil, err
	}
	return &proto_user.GetUsersResponse{
		Users: s.createProtoUsers(result.Data),
		Pagination: &common_proto.PaginationResponse{
			Total:      int32(result.Total),
			TotalPages: int32(result.TotalPages),
			Page:       int32(result.Page),
			PageSize:   int32(result.PageSize),
		},
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
	case proto_user.UserStatus_ACTIVE:
		active := entity.UserStatus(string(entity.UserStatusActive))
		filterResult.Status = &active
	case proto_user.UserStatus_INACTIVE:
		inactive := entity.UserStatus(string(entity.UserStatusInactive))
		filterResult.Status = &inactive
	case proto_user.UserStatus_DELETED:
		deleted := entity.UserStatus(string(entity.UserStatusDeleted))
		filterResult.Status = &deleted
	default:
		filterResult.Status = nil
	}
	if filter.FromDate != nil {
		fromDate, err := time.Parse(time.RFC3339, filter.FromDate.String())
		if err != nil {
			return &filterResult
		}
		filterResult.FromDate = &fromDate
	}
	if filter.ToDate != nil {
		toDate, err := time.Parse(time.RFC3339, filter.ToDate.String())
		if err != nil {
			return &filterResult
		}
		filterResult.ToDate = &toDate
	}
	return &filterResult
}

func (s *userServer) createProtoUsers(users []entity.User) []*proto_user.User {
	protoUsers := make([]*proto_user.User, len(users))
	for i, user := range users {
		protoUsers[i] = s.createProtoUser(user)
	}
	return protoUsers
}
