package grpcservice

import (
	authUC "cms-server/domain/usecase/auth"
	authpb "cms-server/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *authService) ForgotPassword(ctx context.Context, req *authpb.ForgotPasswordRequest) (*authpb.ForgotPasswordResponse, error) {
	// Convert method to usecase type
	var method authUC.ForgotPasswordType
	switch req.GetMethod() {
	case authpb.ForgotPasswordType_FORGOT_BY_CODE:
		method = authUC.ForgotByCode
	case authpb.ForgotPasswordType_FORGOT_BY_TOKEN:
		method = authUC.ForgotByToken
	default:
		return nil, status.Errorf(codes.InvalidArgument, "Phương thức xác thực không hợp lệ")
	}

	// Process forgot password
	result, err := a.forgotPasswordUc.ForgotPassword(req.GetEmail(), req.GetOs(), method)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// Convert user to UserInfo
	userInfo := &authpb.UserInfo{
		Id:       result.User.ID,
		Email:    result.User.Email,
		Phone:    result.User.Phone,
		FullName: result.User.FullName,
		Avatar:   result.User.Avatar,
		Bio:      result.User.Bio,
		Address:  result.User.Address,
	}

	if result.User.Birthday != nil {
		userInfo.Birthday = timestamppb.New(*result.User.Birthday)
	}

	return &authpb.ForgotPasswordResponse{
		User:    userInfo,
		Token:   result.Token,
		Code:    result.Code,
		Message: "Yêu cầu đặt lại mật khẩu đã được gửi. Vui lòng kiểm tra email.",
	}, nil
}
