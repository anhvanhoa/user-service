package grpcservice

import (
	proto "cms-server/proto/gen/auth/v1"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) CheckToken(ctx context.Context, req *proto.CheckTokenRequest) (*proto.CheckTokenResponse, error) {
	ok, err := a.checkTokenUc.CheckToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &proto.CheckTokenResponse{
		Data:    ok,
		Message: "Mã truyền lên hợp lệ",
	}, nil
}
