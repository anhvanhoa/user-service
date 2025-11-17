package user_server

import (
	"context"
	"errors"
	"time"
	"user-service/domain/entity"
	"user-service/domain/usecase/user"
	"user-service/infrastructure/repo"

	hashpass "github.com/anhvanhoa/service-core/domain/hash_pass"
	"github.com/anhvanhoa/service-core/utils"
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

func NewUserServer(db *pg.DB, helper utils.Helper) proto_user.UserServiceServer {
	userRepo := repo.NewUserRepository(db, helper)
	hashService := hashpass.NewArgon()
	userUC := user.NewUserUsecase(
		user.NewCreateUserUsecase(userRepo, hashService),
		user.NewDeleteUserUsecase(userRepo),
		user.NewGetUserUsecase(userRepo),
		user.NewGetUsersUsecase(userRepo, helper),
		user.NewUpdateUserUsecase(userRepo),
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
	user := s.convertReqUpdateToEntity(req)
	updatedUser, err := s.userUsecase.UpdateUserById(req.Id, user)
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
	var deletedAt *timestamppb.Timestamp
	if user.DeletedAt != nil {
		deletedAt = timestamppb.New(*user.DeletedAt)
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
		IsSystem:  user.IsSystem,
		DeletedAt: deletedAt,
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

func (s *userServer) convertReqUpdateToEntity(req *proto_user.UpdateUserRequest) *entity.User {
	var birthday *time.Time
	if req.Birthday != nil {
		birthdayTime := req.Birthday.AsTime()
		birthday = &birthdayTime
	}
	return &entity.User{
		ID:       req.Id,
		Email:    req.Email,
		Phone:    req.Phone,
		FullName: req.FullName,
		Avatar:   req.Avatar,
		Bio:      req.Bio,
		Address:  req.Address,
		Status:   entity.UserStatus(req.Status),
		Birthday: birthday,
	}
}

func (s *userServer) CreateUser(ctx context.Context, req *proto_user.CreateUserRequest) (*proto_user.CreateUserResponse, error) {
	user := s.convertReqCreateToEntity(req)
	createdUser, err := s.userUsecase.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return &proto_user.CreateUserResponse{
		User: s.createProtoUser(createdUser),
	}, nil
}

func (s *userServer) convertReqCreateToEntity(req *proto_user.CreateUserRequest) *entity.User {
	u := &entity.User{
		Email:    req.Email,
		Phone:    req.Phone,
		FullName: req.FullName,
		Avatar:   req.Avatar,
		Bio:      req.Bio,
		Address:  req.Address,
		Password: req.Password,
	}
	if req.Birthday != nil {
		birthdayTime := req.Birthday.AsTime()
		u.Birthday = &birthdayTime
	}
	return u
}
