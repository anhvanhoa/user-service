package grpcservice

import (
	authpb "cms-server/proto"
	"context"
	"regexp"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *authService) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	// Business logic validation: check if email_or_phone is valid format
	identifier := req.GetEmailOrPhone()
	if !isValidEmail(identifier) && !isValidPhone(identifier) {
		return nil, status.Errorf(codes.InvalidArgument, "email hoặc số điện thoại không đúng định dạng")
	}

	// Business logic validation: check password strength
	if err := validatePasswordStrength(req.GetPassword()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Get user by email or phone
	user, err := a.loginUc.GetUserByEmailOrPhone(req.GetEmailOrPhone())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Không tìm thấy người dùng")
	}

	// Check password
	if !a.loginUc.CheckHashPassword(req.GetPassword(), user.Password) {
		return nil, status.Errorf(codes.InvalidArgument, "Mật khẩu không chính xác")
	}

	// Generate tokens
	exp := time.Now().Add(15 * time.Minute) // Access token expires in 15 minutes
	accessToken, err := a.loginUc.GengerateAccessToken(user.ID, user.FullName, exp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể tạo access token")
	}

	refreshExp := time.Now().Add(7 * 24 * time.Hour) // Refresh token expires in 7 days
	refreshToken, err := a.loginUc.GengerateRefreshToken(user.ID, user.FullName, refreshExp, req.GetOs())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể tạo refresh token")
	}

	// Convert user to UserInfo
	userInfo := &authpb.UserInfo{
		Id:       user.ID,
		Email:    user.Email,
		Phone:    user.Phone,
		FullName: user.FullName,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
		Address:  user.Address,
	}

	if user.Birthday != nil {
		userInfo.Birthday = timestamppb.New(*user.Birthday)
	}

	return &authpb.LoginResponse{
		User:         userInfo,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "Đăng nhập thành công",
	}, nil
}

// isValidPhone validates phone number format (Vietnamese format)
func isValidPhone(phone string) bool {
	// Vietnamese phone number regex
	phoneRegex := regexp.MustCompile(`^(0|\+84)(3[2-9]|5[689]|7[06-9]|8[1-689]|9[0-46-9])[0-9]{7}$`)
	return phoneRegex.MatchString(phone)
}
