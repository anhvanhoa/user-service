package grpcservice

import (
	authpb "cms-server/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) ResetPasswordByCode(ctx context.Context, req *authpb.ResetPasswordByCodeRequest) (*authpb.ResetPasswordByCodeResponse, error) {
	// Business logic validation: check if passwords match
	if req.GetNewPassword() != req.GetConfirmPassword() {
		return nil, status.Errorf(codes.InvalidArgument, "Mật khẩu mới và xác nhận mật khẩu không khớp")
	}

	// Verify session
	userID, err := a.resetCodeUc.VerifySession(req.GetCode(), req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Reset password
	if err := a.resetCodeUc.ResetPass(userID, req.GetNewPassword(), req.GetConfirmPassword()); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &authpb.ResetPasswordByCodeResponse{
		Message: "Đặt lại mật khẩu thành công",
	}, nil
}
