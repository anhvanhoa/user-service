package grpcservice

import (
	"context"
	"fmt"
	"net"

	authpb "cms-server/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCServer represents the gRPC server
type GRPCServer struct {
	server *grpc.Server
	port   string
}

// NewGRPCServer creates a new gRPC server with all middleware
func NewGRPCServer(port string, authService authpb.AuthServiceServer) *GRPCServer {
	// Create server with interceptors
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			// Chain interceptors in order
			ValidationInterceptor(),
			LoggingInterceptor(),
			ErrorHandlerInterceptor(),
		),
	)

	// Register services
	authpb.RegisterAuthServiceServer(server, authService)

	// Enable reflection for debugging
	reflection.Register(server)

	return &GRPCServer{
		server: server,
		port:   port,
	}
}

// Start starts the gRPC server
func (s *GRPCServer) Start(ctx context.Context) error {
	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	fmt.Printf("gRPC server starting on port %s\n", s.port)

	// Start server in a goroutine
	go func() {
		if err := s.server.Serve(lis); err != nil {
			fmt.Printf("failed to serve: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	fmt.Println("Shutting down gRPC server...")
	s.server.GracefulStop()

	return nil
}

// Stop stops the gRPC server
func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

// GetServer returns the underlying gRPC server
func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}
