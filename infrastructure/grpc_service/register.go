package grpcservice

import (
	"auth-service/domain/usecase"
	"context"
	"regexp"
	"time"

	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *authService) Register(ctx context.Context, req *proto_auth.RegisterRequest) (*proto_auth.RegisterResponse, error) {
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Business logic validation: check password match
	if err := validatePasswordMatch(req.GetPassword(), req.GetConfirmPassword()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Business logic validation: check verification code format
	if !isValidVerificationCode(req.GetCode()) {
		return nil, status.Error(codes.InvalidArgument, "mã xác thực phải là 6 chữ số")
	}

	// Check if user already exists
	existingUser, err := a.registerUc.CheckUserExist(req.GetEmail())
	if err == nil && existingUser.ID != "" {
		return nil, status.Error(codes.AlreadyExists, "Email đã được sử dụng")
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
	exp := time.Now().Add(15 * time.Minute)
	result, err := a.registerUc.Register(registerReq, req.GetOs(), exp)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	userInfo := &proto_auth.UserInfo{
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

	return &proto_auth.RegisterResponse{
		User:    userInfo,
		Token:   result.Token,
		Message: "Đăng ký thành công. Vui lòng kiểm tra email để xác thực tài khoản.",
	}, nil
}

func isValidVerificationCode(code string) bool {
	codeRegex := regexp.MustCompile(`^[0-9]{6}$`)
	return codeRegex.MatchString(code)
}
