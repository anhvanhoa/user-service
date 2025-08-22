package grpcservice

import (
	authpb "cms-server/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) ResetPasswordByToken(ctx context.Context, req *authpb.ResetPasswordByTokenRequest) (*authpb.ResetPasswordByTokenResponse, error) {
	// Business logic validation: check if passwords match
	if req.GetNewPassword() != req.GetConfirmPassword() {
		return nil, status.Errorf(codes.InvalidArgument, "Mật khẩu mới và xác nhận mật khẩu không khớp")
	}

	// Verify session
	userID, err := a.resetTokenUc.VerifySession(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Reset password
	if err := a.resetTokenUc.ResetPass(userID, req.GetNewPassword(), req.GetConfirmPassword()); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &authpb.ResetPasswordByTokenResponse{
		Message: "Đặt lại mật khẩu thành công",
	}, nil
}
