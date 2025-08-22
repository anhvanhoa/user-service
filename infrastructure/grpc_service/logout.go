package grpcservice

import (
	authpb "cms-server/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	// Verify token
	if err := a.logoutUc.VerifyToken(req.GetToken()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Token không hợp lệ")
	}

	// Logout user
	if err := a.logoutUc.Logout(req.GetToken()); err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể đăng xuất")
	}

	return &authpb.LogoutResponse{
		Message: "Đăng xuất thành công",
	}, nil
}
