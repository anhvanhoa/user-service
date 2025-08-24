package grpcservice

import (
	"context"

	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) Logout(ctx context.Context, req *proto_auth.LogoutRequest) (*proto_auth.LogoutResponse, error) {
	// Verify token
	if err := a.logoutUc.VerifyToken(req.GetToken()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Token không hợp lệ")
	}

	// Logout user
	if err := a.logoutUc.Logout(req.GetToken()); err != nil {
		return nil, status.Error(codes.Internal, "Không thể đăng xuất")
	}

	return &proto_auth.LogoutResponse{
		Message: "Đăng xuất thành công",
	}, nil
}
