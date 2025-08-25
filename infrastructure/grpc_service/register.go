package grpcservice

import (
	"auth-service/constants"
	"auth-service/domain/service/queue"
	"auth-service/domain/usecase"
	"context"
	"time"

	proto_mail_history "github.com/anhvanhoa/sf-proto/gen/mail_history/v1"
	proto_mail_template "github.com/anhvanhoa/sf-proto/gen/mail_tmpl/v1"

	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *authService) Register(ctx context.Context, req *proto_auth.RegisterRequest) (*proto_auth.RegisterResponse, error) {
	if !isValidEmail(req.GetEmail()) {
		return nil, status.Errorf(codes.InvalidArgument, "email không đúng định dạng")
	}

	if len(req.GetFullName()) < 2 {
		return nil, status.Errorf(codes.InvalidArgument, "họ tên phải có ít nhất 2 ký tự")
	}

	if err := validatePasswordStrength(req.GetPassword()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := validatePasswordMatch(req.GetPassword(), req.GetConfirmPassword()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	existingUser, err := a.registerUc.CheckUserExist(req.GetEmail())
	if err == nil && existingUser.ID != "" {
		return nil, status.Error(codes.AlreadyExists, "Email đã được sử dụng")
	}

	code := a.registerUc.GengerateCode(6)
	registerReq := usecase.RegisterReq{
		Email:           req.GetEmail(),
		FullName:        req.GetFullName(),
		Password:        req.GetPassword(),
		ConfirmPassword: req.GetConfirmPassword(),
		Code:            code,
	}

	exp := time.Now().Add(15 * time.Minute)
	result, err := a.registerUc.Register(registerReq, "web", exp)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tmpl, err := a.mailService.Mtc.GetMailTmpl(ctx, &proto_mail_template.GetMailTmplRequest{
		Id: constants.TPL_REGISTER_MAIL,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data := map[string]any{
		"link": a.env.FRONTEND_URL + "/auth/verify/" + result.Token,
		"user": result.UserInfor,
	}
	payload := queue.Payload{
		Provider: tmpl.MailTmpl.ProviderEmail,
		Template: tmpl.MailTmpl.Id,
		Data:     data,
		Tos:      &[]string{result.UserInfor.Email},
	}

	Id, err := a.registerUc.SendMail(payload)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoData := make(map[string]*anypb.Any)
	for k := range data {
		if anyValue, err := anypb.New(timestamppb.New(time.Now())); err == nil {
			protoData[k] = anyValue
		}
	}

	a.mailService.Mhc.CreateMailHistory(ctx, &proto_mail_history.CreateMailHistoryRequest{
		Id:            Id,
		TemplateId:    constants.TPL_REGISTER_MAIL,
		Subject:       tmpl.MailTmpl.Subject,
		Body:          tmpl.MailTmpl.Body,
		Tos:           []string{result.UserInfor.Email},
		Data:          protoData,
		EmailProvider: tmpl.MailTmpl.ProviderEmail,
	})

	userInfo := &proto_auth.UserInfo{
		Id:       result.UserInfor.ID,
		Email:    result.UserInfor.Email,
		Phone:    result.UserInfor.Phone,
		FullName: result.UserInfor.FullName,
		Avatar:   result.UserInfor.Avatar,
		Bio:      result.UserInfor.Bio,
		Address:  result.UserInfor.Address,
	}

	if result.UserInfor.Birthday != nil {
		userInfo.Birthday = timestamppb.New(*result.UserInfor.Birthday)
	}

	return &proto_auth.RegisterResponse{
		User:    userInfo,
		Token:   result.Token,
		Message: "Đăng ký thành công. Vui lòng kiểm tra email để xác thực tài khoản.",
	}, nil
}
