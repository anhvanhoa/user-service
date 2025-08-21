package main

import (
	"cms-server/bootstrap"
	grpcservice "cms-server/infrastructure/grpc_service"
	pkglog "cms-server/infrastructure/service/logger"
	authpb "cms-server/proto"
	"context"
	"net"

	"github.com/go-pg/pg/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer *grpc.Server
	env        *bootstrap.Env
	log        pkglog.Logger
}

func newGRPCServer(db *pg.DB, env *bootstrap.Env, log pkglog.Logger) *Server {
	s := grpc.NewServer()
	// health
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(s, healthServer)
	// reflection (dev only)
	reflection.Register(s)

	authpb.RegisterAuthServiceServer(s, grpcservice.NewAuthService(db))
	return &Server{grpcServer: s, env: env, log: log}
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.env.HOST_GRPC+":"+s.env.PORT_GRPC)
	if err != nil {
		s.log.Error("failed to listen gRPC", "error", err.Error())
		return err
	}
	go func() {
		<-ctx.Done()
		s.grpcServer.GracefulStop()
	}()
	s.log.Info("gRPC server starting", "port", s.env.PORT_GRPC)
	return s.grpcServer.Serve(lis)
}

func (s *Server) GetGRPC() *grpc.Server {
	return s.grpcServer
}
