package usecase

import (
	"auth-service/constants"
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"auth-service/domain/service/cache"
	se "auth-service/domain/service/error"
	serviceJwt "auth-service/domain/service/jwt"
	"math/rand"
	"strconv"
	"time"
)

type ForgotPasswordType string

const (
	ForgotByCode  ForgotPasswordType = "ForgotByCode"
	ForgotByToken ForgotPasswordType = "ForgotByToken"
)

var (
	ErrValidateForgotPassword = se.NewErr("Phương thức xác thực không hợp lệ, vui lòng chọn code hoặc token")
)

type ForgotPasswordRes struct {
	User  entity.UserInfor
	Token string
	Code  string
}

type ForgotPasswordUsecase interface {
	ForgotPassword(email, os string, method ForgotPasswordType) (ForgotPasswordRes, error)
	saveCodeOrToken(typeForgot ForgotPasswordType, userID, codeOrToken, os string, exp time.Time) error
	SendEmailForgotPassword(user entity.UserInfor, code, link string) error
	generateRandomCode(length int) string
}

type forgotPasswordUsecaseImpl struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	tx          repository.ManagerTransaction
	jwtService  serviceJwt.JwtService
	cache       cache.RedisConfigImpl
}

func NewForgotPasswordUsecase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	tx repository.ManagerTransaction,
	jwtService serviceJwt.JwtService,
	cache cache.RedisConfigImpl,
) ForgotPasswordUsecase {
	return &forgotPasswordUsecaseImpl{
		userRepo,
		sessionRepo,
		tx,
		jwtService,
		cache,
	}
}

func (uc *forgotPasswordUsecaseImpl) saveCodeOrToken(typeForgot ForgotPasswordType, userID, codeOrToken, os string, exp time.Time) error {
	session := entity.Session{
		Token:     codeOrToken,
		UserID:    userID,
		Type:      entity.SessionTypeForgot,
		Os:        os,
		ExpiredAt: exp,
		CreatedAt: time.Now(),
	}
	key := codeOrToken
	if typeForgot == ForgotByCode && len(codeOrToken) == 6 {
		key = codeOrToken + userID
	}
	if err := uc.cache.Set(key, []byte(codeOrToken), constants.ForgotExpiredAt*time.Minute); err != nil {
		if err := uc.sessionRepo.CreateSession(session); err != nil {
			return se.NewErr("không thể tạo phiên làm việc")
		}
	} else {
		go uc.sessionRepo.CreateSession(session)
	}
	return nil
}

func (uc *forgotPasswordUsecaseImpl) SendEmailForgotPassword(user entity.UserInfor, code, link string) error {
	go uc.sessionRepo.DeleteAllSessionsForgot()
	return nil
}

func (uc *forgotPasswordUsecaseImpl) ForgotPassword(email, os string, method ForgotPasswordType) (ForgotPasswordRes, error) {
	var resForgotPassword ForgotPasswordRes
	user, err := uc.userRepo.GetUserByEmail(email)
	if err != nil {
		return resForgotPassword, err
	}
	resForgotPassword.User = user.GetInfor()
	exp := time.Now().Add(constants.ForgotExpiredAt * time.Minute)
	switch method {
	case ForgotByCode:
		resForgotPassword.Code = uc.generateRandomCode(6)
		if err := uc.saveCodeOrToken(ForgotByCode, user.ID, resForgotPassword.Code, os, exp); err != nil {
			return resForgotPassword, err
		}
		return resForgotPassword, nil
	case ForgotByToken:
		code := uc.generateRandomCode(6)
		resForgotPassword.Token, err = uc.jwtService.GenForgotPasswordToken(user.ID, code, exp)
		if err != nil {
			return resForgotPassword, err
		}
		if err := uc.saveCodeOrToken(ForgotByToken, user.ID, resForgotPassword.Token, os, exp); err != nil {
			return resForgotPassword, err
		}
		return resForgotPassword, nil
	}
	return resForgotPassword, ErrValidateForgotPassword
}

func (uc *forgotPasswordUsecaseImpl) generateRandomCode(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	min := int64(1)
	for i := 1; i < length; i++ {
		min *= 10
	}
	max := min*10 - 1
	num := r.Int63n(max-min+1) + min
	return strconv.FormatInt(num, 10)
}
