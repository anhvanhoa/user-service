package grpcservice

import (
	"cms-server/domain/usecase"
	proto "cms-server/proto/gen/auth/v1"
	"context"
	"regexp"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *authService) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	// Business logic validation: check email format
	if !isValidEmail(req.GetEmail()) {
		return nil, status.Errorf(codes.InvalidArgument, "email không đúng định dạng")
	}

	// Business logic validation: check full name length
	if len(req.GetFullName()) < 2 {
		return nil, status.Errorf(codes.InvalidArgument, "họ tên phải có ít nhất 2 ký tự")
	}

	// Business logic validation: check password strength
	if err := validatePasswordStrength(req.GetPassword()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Business logic validation: check password match
	if err := validatePasswordMatch(req.GetPassword(), req.GetConfirmPassword()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Business logic validation: check verification code format
	if !isValidVerificationCode(req.GetCode()) {
		return nil, status.Errorf(codes.InvalidArgument, "mã xác thực phải là 6 chữ số")
	}

	// Check if user already exists
	existingUser, err := a.registerUc.CheckUserExist(req.GetEmail())
	if err == nil && existingUser.ID != "" {
		return nil, status.Errorf(codes.AlreadyExists, "Email đã được sử dụng")
	}

	// Create register request
	registerReq := usecase.RegisterReq{
		Email:           req.GetEmail(),
		FullName:        req.GetFullName(),
		Password:        req.GetPassword(),
		ConfirmPassword: req.GetConfirmPassword(),
		Code:            req.GetCode(),
	}

	// Register user
	exp := time.Now().Add(15 * time.Minute) // Registration token expires in 15 minutes
	result, err := a.registerUc.Register(registerReq, req.GetOs(), exp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// Convert user to UserInfo
	userInfo := &proto.UserInfo{
		Id:       result.UserInfor.ID,
		Email:    result.UserInfor.Email,
		Phone:    result.UserInfor.Phone,
		FullName: result.UserInfor.FullName,
		Avatar:   result.UserInfor.Avatar,
		Bio:      result.UserInfor.Bio,
		Address:  result.UserInfor.Address,
	}

	if result.UserInfor.Birthday != nil {
		userInfo.Birthday = timestamppb.New(*result.UserInfor.Birthday)
	}

	return &proto.RegisterResponse{
		User:    userInfo,
		Token:   result.Token,
		Message: "Đăng ký thành công. Vui lòng kiểm tra email để xác thực tài khoản.",
	}, nil
}

// isValidVerificationCode validates verification code format
func isValidVerificationCode(code string) bool {
	codeRegex := regexp.MustCompile(`^[0-9]{6}$`)
	return codeRegex.MatchString(code)
}
