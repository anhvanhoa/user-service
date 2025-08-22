package grpcservice

import (
	authpb "cms-server/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) VerifyAccount(ctx context.Context, req *authpb.VerifyAccountRequest) (*authpb.VerifyAccountResponse, error) {
	// Verify register token
	claims, err := a.verifyAccountUc.VerifyRegister(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Token không hợp lệ hoặc đã hết hạn")
	}

	// Get user by ID
	_, err = a.verifyAccountUc.GetUserById(claims.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Không tìm thấy người dùng")
	}

	// Verify account
	if err := a.verifyAccountUc.VerifyAccount(claims.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể xác thực tài khoản")
	}

	return &authpb.VerifyAccountResponse{
		Message: "Xác thực tài khoản thành công",
	}, nil
}
