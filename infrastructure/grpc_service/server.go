package grpcservice

import (
	"context"
	"fmt"
	"net"

	proto "cms-server/proto/gen/auth/v1"

	"buf.build/go/protovalidate"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server *grpc.Server
	port   string
}

func NewGRPCServer(port string, authService proto.AuthServiceServer) *GRPCServer {
	validator, err := protovalidate.New()
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			protovalidate_middleware.UnaryServerInterceptor(validator),
			LoggingInterceptor(),
		),
	)

	proto.RegisterAuthServiceServer(server, authService)

	reflection.Register(server)

	return &GRPCServer{
		server: server,
		port:   port,
	}
}

func (s *GRPCServer) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	fmt.Printf("gRPC server starting on port %s\n", s.port)

	go func() {
		if err := s.server.Serve(lis); err != nil {
			fmt.Printf("failed to serve: %v", err)
		}
	}()

	<-ctx.Done()

	fmt.Println("Shutting down gRPC server...")
	s.server.GracefulStop()

	return nil
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}
