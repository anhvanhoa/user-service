package handler

import (
	"cms-server/bootstrap"
	authUC "cms-server/domain/usecase/auth"
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
	checkTokenUc    authUC.CheckTokenUsecase
	checkCodeUc     authUC.CheckCodeUsecase
	forgotUc        authUC.ForgotPasswordUsecase
	loginUc         authUC.LoginUsecase
	logoutUc        authUC.LogoutUsecase
	refreshUc       authUC.RefreshUsecase
	registerUc      authUC.RegisterUsecase
	resetCodeUc     authUC.ResetPasswordByCodeUsecase
	resetTokenUc    authUC.ResetPasswordByTokenUsecase
	verifyAccountUc authUC.VerifyAccountUsecase
}

func NewAuthHandler(
	checkTokenUc authUC.CheckTokenUsecase,
	checkCodeUc authUC.CheckCodeUsecase,
	forgotUc authUC.ForgotPasswordUsecase,
	loginUc authUC.LoginUsecase,
	logoutUc authUC.LogoutUsecase,
	refreshUc authUC.RefreshUsecase,
	registerUc authUC.RegisterUsecase,
	resetCodeUc authUC.ResetPasswordByCodeUsecase,
	resetTokenUc authUC.ResetPasswordByTokenUsecase,
	verifyAccountUc authUC.VerifyAccountUsecase,
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
