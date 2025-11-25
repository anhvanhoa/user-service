package user_server

import (
	"context"
	"user-service/domain/entity"

	proto_user "github.com/anhvanhoa/sf-proto/gen/user/v1"
)

func (s *userServer) GetUserMap(ctx context.Context, req *proto_user.GetUserMapRequest) (*proto_user.GetUserMapResponse, error) {
	users, err := s.userUsecase.GetUserMapUsecase.Excute(req.Ids)
	if err != nil {
		return nil, err
	}
	return &proto_user.GetUserMapResponse{
		UserMap: s.createProtoUserInfoMap(users),
	}, nil
}

func (s *userServer) createProtoUserInfoMap(users map[string]entity.User) map[string]*proto_user.UserInfo {
	protoUsers := make(map[string]*proto_user.UserInfo)
	for _, user := range users {
		protoUsers[user.ID] = s.createProtoUserInfo(user.GetInfor())
	}
	return protoUsers
}
