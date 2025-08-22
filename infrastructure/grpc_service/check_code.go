package grpcservice

import (
	authpb "cms-server/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) CheckCode(ctx context.Context, req *authpb.CheckCodeRequest) (*authpb.CheckCodeResponse, error) {
	// Check code
	valid, err := a.checkCodeUc.CheckCode(req.GetCode(), req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if !valid {
		return &authpb.CheckCodeResponse{
			Valid:   false,
			Message: "Mã xác thực không hợp lệ hoặc đã hết hạn",
		}, nil
	}

	return &authpb.CheckCodeResponse{
		Valid:   true,
		Message: "Mã xác thực hợp lệ",
	}, nil
}
