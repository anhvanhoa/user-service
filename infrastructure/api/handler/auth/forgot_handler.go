package handler

import (
	authUC "cms-server/domain/usecase/auth"
	pkgres "cms-server/infrastructure/service/response"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func (h *authHandlerImpl) Forgot(c *fiber.Ctx) error {
	var body forgotPasswordReq
	if err := c.BodyParser(&body); err != nil {
		err := pkgres.NewErr("Dữ liệu không hợp lệ").BadReq()
		return h.log.Log(c, err)
	}

	if err := h.validate.ValidateStruct(body); err != nil {
		err := pkgres.NewErr("Dữ liệu không hợp lệ").SetData(err.Data).BadReq()
		return h.log.Log(c, err)
	}

	os := c.Get("User-Agent")
	resForpass, err := h.forgotUc.ForgotPassword(body.Email, os, body.Type)

	if errors.Is(err, authUC.ErrValidateForgotPassword) {
		err := pkgres.Err(err).BadReq()
		return h.log.Log(c, err)
	} else if err != nil {
		err := pkgres.NewErr("Không tìm thấy tài khoản, hãy kiểm tra lại").BadReq()
		return h.log.Log(c, err)
	}

	var link string
	if body.Type == authUC.ForgotByToken {
		link = h.env.FRONTEND_URL + "/auth/forgot-password?code=" + resForpass.Token
	}
	if err := h.forgotUc.SendEmailForgotPassword(resForpass.User, resForpass.Code, link); err != nil {
		err := pkgres.Err(err).Code(fiber.StatusInternalServerError)
		return h.log.Log(c, err)
	}
	return c.JSON(pkgres.NewRes("Yêu cầu đặt lại mật khẩu đã được gửi đến email của bạn"))
}
