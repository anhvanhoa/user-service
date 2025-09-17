package user_server

import (
	"context"
	"errors"
	"time"
	"user-service/domain/entity"
	"user-service/domain/usecase/user"
	"user-service/infrastructure/repo"

	"github.com/anhvanhoa/service-core/common"
	proto_user "github.com/anhvanhoa/sf-proto/gen/user/v1"
	"github.com/go-pg/pg/v10"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrUserNotFound = errors.New("không tìm thấy người dùng")
)

type userServer struct {
	proto_user.UnsafeUserServiceServer
	userUsecase user.UserUsecaseI
}

func NewUserServer(db *pg.DB) proto_user.UserServiceServer {
	userRepo := repo.NewUserRepository(db)
	userRoleRepo := repo.NewUserRoleRepository(db)
	userUC := user.NewUserUsecase(
		user.NewDeleteUserUsecase(userRepo),
		user.NewGetUserUsecase(userRepo),
		user.NewUpdateUserUsecase(userRepo, userRoleRepo),
		user.NewUpdateUserRolesUsecase(userRoleRepo),
	)
	return &userServer{
		userUsecase: userUC,
	}
}

func (s *userServer) GetUserById(ctx context.Context, req *proto_user.GetUserByIdRequest) (*proto_user.GetUserByIdResponse, error) {
	user, err := s.userUsecase.GetUserById(req.Id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &proto_user.GetUserByIdResponse{
		User: s.createProtoUser(user),
	}, nil
}

func (s *userServer) DeleteUser(ctx context.Context, req *proto_user.DeleteUserRequest) (*proto_user.DeleteUserResponse, error) {
	err := s.userUsecase.DeleteUserById(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto_user.DeleteUserResponse{
		Message: "Delete user successfully",
		Success: true,
	}, nil
}

func (s *userServer) UpdateUser(ctx context.Context, req *proto_user.UpdateUserRequest) (*proto_user.UpdateUserResponse, error) {
	user := s.createEntityUser(req)
	updatedUser, err := s.userUsecase.UpdateUserById(req.Id, user, req.RoleIds)
	if err != nil {
		return nil, err
	}
	return &proto_user.UpdateUserResponse{
		UserInfo: s.createProtoUserInfo(updatedUser),
	}, nil
}

func (s *userServer) createProtoUser(user entity.User) *proto_user.User {
	var birthday *timestamppb.Timestamp
	if user.Birthday != nil {
		birthday = timestamppb.New(*user.Birthday)
	}
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt != nil {
		updatedAt = timestamppb.New(*user.UpdatedAt)
	}
	return &proto_user.User{
		Id:        user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		FullName:  user.FullName,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		Address:   user.Address,
		Status:    string(user.Status),
		CreatedBy: user.CreatedBy,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
		Birthday:  birthday,
	}
}

func (s *userServer) createProtoUserInfo(user entity.UserInfor) *proto_user.UserInfo {
	var birthday *timestamppb.Timestamp
	if user.Birthday != nil {
		birthday = timestamppb.New(*user.Birthday)
	}
	return &proto_user.UserInfo{
		Id:       user.ID,
		Email:    user.Email,
		Phone:    user.Phone,
		FullName: user.FullName,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
		Address:  user.Address,
		Birthday: birthday,
	}
}

func (s *userServer) createEntityUser(req *proto_user.UpdateUserRequest) entity.User {
	var birthday *time.Time
	if req.Birthday != nil {
		birthdayTime, err := time.Parse(time.RFC3339, req.Birthday.String())
		if err != nil {
			return entity.User{}
		}
		birthday = &birthdayTime
	}
	return entity.User{
		ID:       req.Id,
		Email:    req.Email,
		Phone:    req.Phone,
		FullName: req.FullName,
		Avatar:   req.Avatar,
		Bio:      req.Bio,
		Address:  req.Address,
		Status:   common.Status(req.Status),
		Birthday: birthday,
	}
}
