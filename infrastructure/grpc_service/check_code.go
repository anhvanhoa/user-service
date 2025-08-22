package grpcservice

import (
	proto "cms-server/proto/gen/auth/v1"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) CheckCode(ctx context.Context, req *proto.CheckCodeRequest) (*proto.CheckCodeResponse, error) {
	// Check code
	valid, err := a.checkCodeUc.CheckCode(req.GetCode(), req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if !valid {
		return &proto.CheckCodeResponse{
			Valid:   false,
			Message: "Mã xác thực không hợp lệ hoặc đã hết hạn",
		}, nil
	}

	return &proto.CheckCodeResponse{
		Valid:   true,
		Message: "Mã xác thực hợp lệ",
	}, nil
}
