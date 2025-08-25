package grpcservice

import (
	"context"
	"fmt"
	"net"

	loggerI "auth-service/domain/service/logger"

	"buf.build/go/protovalidate"
	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server *grpc.Server
	port   string
	log    loggerI.Log
}

func NewGRPCServer(port string, authService proto_auth.AuthServiceServer, log loggerI.Log) *GRPCServer {
	validator, err := protovalidate.New()
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			protovalidate_middleware.UnaryServerInterceptor(validator),
			LoggingInterceptor(log),
		),
	)

	proto_auth.RegisterAuthServiceServer(server, authService)

	reflection.Register(server)

	return &GRPCServer{
		server: server,
		port:   port,
		log:    log,
	}
}

func (s *GRPCServer) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.log.Info(fmt.Sprintf("gRPC server starting on port %s", s.port))

	go func() {
		if err := s.server.Serve(lis); err != nil {
			s.log.Error(fmt.Sprintf("failed to serve: %v", err))
		}
	}()

	<-ctx.Done()

	s.log.Info("Shutting down gRPC server...")
	s.server.GracefulStop()

	return nil
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}
