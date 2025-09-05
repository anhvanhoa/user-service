package session_server

import (
	"context"
	"user-service/domain/entity"
	"user-service/domain/usecase/session"
	"user-service/infrastructure/repo"

	"github.com/anhvanhoa/service-core/domain/cache"
	proto_session "github.com/anhvanhoa/sf-proto/gen/session/v1"
	"github.com/go-pg/pg/v10"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type sessionServer struct {
	proto_session.UnimplementedSessionServiceServer
	sessionUsecase session.SessionUsecaseI
}

func NewSessionServer(
	db *pg.DB,
	cache cache.CacheI,
) proto_session.SessionServiceServer {
	sessionRepo := repo.NewSessionRepository(db)
	sessionUC := session.NewSessionUsecase(
		session.NewGetSessionsUsecase(sessionRepo),
		session.NewGetSessionByUser(sessionRepo),
		session.NewDeleteSessionByTypeAndTokenUsecase(sessionRepo, cache),
		session.NewDeleteSessionByTypeAndUserUsecase(sessionRepo, cache),
		session.NewDeleteSessionExpiredUsecase(sessionRepo),
	)
	return &sessionServer{
		sessionUsecase: sessionUC,
	}
}

func (s *sessionServer) GetSessions(ctx context.Context, req *proto_session.GetAllSessionsRequest) (*proto_session.GetAllSessionsResponse, error) {
	sessions, err := s.sessionUsecase.GetSessions(ctx)
	if err != nil {
		return nil, err
	}
	return &proto_session.GetAllSessionsResponse{
		Sessions: s.createProtoSessions(sessions),
	}, nil
}

func (s *sessionServer) GetSessionsByUserId(ctx context.Context, req *proto_session.GetSessionsByUserIdRequest) (*proto_session.GetSessionsByUserIdResponse, error) {
	sessions, err := s.sessionUsecase.GetSessionsByUserId(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &proto_session.GetSessionsByUserIdResponse{
		Sessions: s.createProtoSessions(sessions),
	}, nil
}

func (s *sessionServer) DeleteSessionByTypeAndToken(ctx context.Context, req *proto_session.DeleteSessionByTypeAndTokenRequest) (*proto_session.DeleteSessionByTypeAndTokenResponse, error) {
	err := s.sessionUsecase.DeleteSessionByTypeAndToken(ctx, entity.SessionType(req.Type), req.Token)
	if err != nil {
		return nil, err
	}
	return &proto_session.DeleteSessionByTypeAndTokenResponse{
		Message: "Delete session by type and token successfully",
		Success: true,
	}, nil
}

func (s *sessionServer) DeleteSessionByTypeAndUser(ctx context.Context, req *proto_session.DeleteSessionByTypeAndUserRequest) (*proto_session.DeleteSessionByTypeAndUserResponse, error) {
	err := s.sessionUsecase.DeleteSessionByTypeAndUser(ctx, entity.SessionType(req.Type), req.UserId)
	if err != nil {
		return nil, err
	}
	return &proto_session.DeleteSessionByTypeAndUserResponse{
		Message: "Delete session by type and user successfully",
		Success: true,
	}, nil
}

func (s *sessionServer) DeleteSessionExpired(ctx context.Context, req *proto_session.DeleteSessionExpiredRequest) (*proto_session.DeleteSessionExpiredResponse, error) {
	err := s.sessionUsecase.DeleteSessionExpired(ctx)
	if err != nil {
		return nil, err
	}
	return &proto_session.DeleteSessionExpiredResponse{
		Message: "Delete session expired successfully",
		Success: true,
	}, nil
}

func (s *sessionServer) createProtoSessions(sessions []entity.Session) []*proto_session.Session {
	protoSessions := make([]*proto_session.Session, len(sessions))
	for i, session := range sessions {
		protoSessions[i] = &proto_session.Session{
			Token:     session.Token,
			UserId:    session.UserID,
			Type:      string(session.Type),
			Os:        session.Os,
			ExpiredAt: timestamppb.New(session.ExpiredAt),
			CreatedAt: timestamppb.New(session.CreatedAt),
		}
	}
	return protoSessions
}
