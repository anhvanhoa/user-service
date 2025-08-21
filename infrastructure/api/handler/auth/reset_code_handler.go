package handler

import (
	pkgres "cms-server/infrastructure/service/response"

	"github.com/gofiber/fiber/v2"
)

func (h *authHandlerImpl) ResetCode(c *fiber.Ctx) error {
	var req resetPasswordByCodeReq
	if err := c.BodyParser(&req); err != nil {
		err := pkgres.NewErr("Dữ liệu không hợp lệ").BadReq()
		return h.log.Log(c, err)
	}

	if err := h.validate.ValidateStruct(req); err != nil {
		err := pkgres.NewErr("Dữ liệu không hợp lệ").SetData(err.Data).BadReq()
		return h.log.Log(c, err)
	}

	if req.Password != req.ConfirmPassword {
		err := pkgres.NewErr("Mật khẩu không khớp").BadReq()
		return h.log.Log(c, err)
	}

	userID, err := h.resetCodeUc.VerifySession(req.Code, req.Email)
	if err != nil {
		err = pkgres.Err(err).BadReq()
		return h.log.Log(c, err)
	}

	if err := h.resetCodeUc.ResetPass(userID, req.Password, req.ConfirmPassword); err != nil {
		err = pkgres.Err(err).BadReq()
		return h.log.Log(c, err)
	}

	return c.JSON(pkgres.NewRes("Cập nhật mật khẩu thành công"))
}
