package grpcservice

import (
	proto "cms-server/proto/gen/auth/v1"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) Logout(ctx context.Context, req *proto.LogoutRequest) (*proto.LogoutResponse, error) {
	// Verify token
	if err := a.logoutUc.VerifyToken(req.GetToken()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Token không hợp lệ")
	}

	// Logout user
	if err := a.logoutUc.Logout(req.GetToken()); err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể đăng xuất")
	}

	return &proto.LogoutResponse{
		Message: "Đăng xuất thành công",
	}, nil
}
