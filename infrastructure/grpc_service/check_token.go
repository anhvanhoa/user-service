package grpcservice

import (
	authpb "cms-server/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *authService) CheckToken(ctx context.Context, req *authpb.CheckTokenRequest) (*authpb.CheckTokenResponse, error) {
	if req.GetToken() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Không tìm thấy mã truyền lên")
	}
	ok, err := a.checkTokenUc.CheckToken(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &authpb.CheckTokenResponse{
		Data:    ok,
		Message: "Mã truyền lên hợp lệ",
	}, nil
}
