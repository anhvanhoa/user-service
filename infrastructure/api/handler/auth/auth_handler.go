package handler

import (
	"cms-server/bootstrap"
	"cms-server/domain/usecase"
	pkglog "cms-server/infrastructure/service/logger"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	CheckCode(c *fiber.Ctx) error
	CheckToken(c *fiber.Ctx) error
	Forgot(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	ResetCode(c *fiber.Ctx) error
	ResetToken(c *fiber.Ctx) error
	VerifyAccount(c *fiber.Ctx) error
}

type authHandlerImpl struct {
	env             *bootstrap.Env
	validate        bootstrap.IValidator
	log             pkglog.Logger
	checkTokenUc    usecase.CheckTokenUsecase
	checkCodeUc     usecase.CheckCodeUsecase
	forgotUc        usecase.ForgotPasswordUsecase
	loginUc         usecase.LoginUsecase
	logoutUc        usecase.LogoutUsecase
	refreshUc       usecase.RefreshUsecase
	registerUc      usecase.RegisterUsecase
	resetCodeUc     usecase.ResetPasswordByCodeUsecase
	resetTokenUc    usecase.ResetPasswordByTokenUsecase
	verifyAccountUc usecase.VerifyAccountUsecase
}

func NewAuthHandler(
	checkTokenUc usecase.CheckTokenUsecase,
	checkCodeUc usecase.CheckCodeUsecase,
	forgotUc usecase.ForgotPasswordUsecase,
	loginUc usecase.LoginUsecase,
	logoutUc usecase.LogoutUsecase,
	refreshUc usecase.RefreshUsecase,
	registerUc usecase.RegisterUsecase,
	resetCodeUc usecase.ResetPasswordByCodeUsecase,
	resetTokenUc usecase.ResetPasswordByTokenUsecase,
	verifyAccountUc usecase.VerifyAccountUsecase,
	log pkglog.Logger,
	env *bootstrap.Env,
	validate bootstrap.IValidator,
) AuthHandler {
	return &authHandlerImpl{
		validate:        validate,
		env:             env,
		log:             log,
		checkCodeUc:     checkCodeUc,
		checkTokenUc:    checkTokenUc,
		forgotUc:        forgotUc,
		loginUc:         loginUc,
		logoutUc:        logoutUc,
		refreshUc:       refreshUc,
		registerUc:      registerUc,
		resetCodeUc:     resetCodeUc,
		resetTokenUc:    resetTokenUc,
		verifyAccountUc: verifyAccountUc,
	}
}
