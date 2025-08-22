package grpcservice

import (
	"cms-server/domain/usecase"
	proto "cms-server/proto/gen/auth/v1"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *authService) ForgotPassword(ctx context.Context, req *proto.ForgotPasswordRequest) (*proto.ForgotPasswordResponse, error) {
	// Convert method to usecase type
	var method usecase.ForgotPasswordType
	switch req.GetMethod() {
	case proto.ForgotPasswordType_FORGOT_PASSWORD_TYPE_UNSPECIFIED:
		method = usecase.ForgotByCode
	case proto.ForgotPasswordType_FORGOT_PASSWORD_TYPE_TOKEN:
		method = usecase.ForgotByToken
	default:
		return nil, status.Errorf(codes.InvalidArgument, "Phương thức xác thực không hợp lệ")
	}

	// Process forgot password
	result, err := a.forgotPasswordUc.ForgotPassword(req.GetEmail(), req.GetOs(), method)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// Convert user to UserInfo
	userInfo := &proto.UserInfo{
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

	return &proto.ForgotPasswordResponse{
		User:    userInfo,
		Token:   result.Token,
		Code:    result.Code,
		Message: "Yêu cầu đặt lại mật khẩu đã được gửi. Vui lòng kiểm tra email.",
	}, nil
}
