package grpcservice

import (
	authpb "cms-server/proto"
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	// Get session by token
	if _, err := a.refreshUc.GetSessionByToken(req.GetRefreshToken()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Phiên làm việc không hợp lệ")
	}

	// Verify token
	claims, err := a.refreshUc.VerifyToken(req.GetRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Token không hợp lệ")
	}

	// Clear expired sessions
	if err := a.refreshUc.ClearSessionExpired(); err != nil {
		// Log error but continue
	}

	// Generate new access token
	accessExp := time.Now().Add(15 * time.Minute) // Access token expires in 15 minutes
	accessToken, err := a.refreshUc.GengerateAccessToken(claims.Id, claims.FullName, accessExp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể tạo access token")
	}

	// Generate new refresh token
	refreshExp := time.Now().Add(7 * 24 * time.Hour) // Refresh token expires in 7 days
	refreshToken, err := a.refreshUc.GengerateRefreshToken(claims.Id, claims.FullName, refreshExp, req.GetOs())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể tạo refresh token")
	}

	return &authpb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "Làm mới token thành công",
	}, nil
}
