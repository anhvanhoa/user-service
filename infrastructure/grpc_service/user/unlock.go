package user_server

import (
	"context"

	proto_user "github.com/anhvanhoa/sf-proto/gen/user/v1"
)

func (s *userServer) UnlockUser(ctx context.Context, req *proto_user.UnlockUserRequest) (*proto_user.UnlockUserResponse, error) {
	err := s.userUsecase.UnlockUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &proto_user.UnlockUserResponse{
		Message: "Mở khóa người dùng thành công",
	}, nil
}
