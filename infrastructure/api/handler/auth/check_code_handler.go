package handler

import (
	pkgres "cms-server/infrastructure/service/response"

	"github.com/gofiber/fiber/v2"
)

func (h *authHandlerImpl) CheckCode(c *fiber.Ctx) error {
	var req checkCodeReq
	if err := c.BodyParser(&req); err != nil {
		err := pkgres.NewErr("Dữ liệu không hợp lệ").BadReq()
		return h.log.Log(c, err)
	}

	if err := h.validate.ValidateStruct(req); err != nil {
		err := pkgres.NewErr("Dữ liệu không hợp lệ").BadReq()
		return h.log.Log(c, err)
	}

	ok, err := h.checkCodeUc.CheckCode(req.Code, req.Email)
	if err != nil {
		err := pkgres.Err(err).BadReq()
		return h.log.Log(c, err)
	}

	if !ok {
		err := pkgres.NewErr("Mã xác thực không hợp lệ").BadReq()
		return h.log.Log(c, err)
	}
	return c.JSON(pkgres.NewRes("Mã xác thực hợp lệ"))
}
