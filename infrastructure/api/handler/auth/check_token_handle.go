package handler

import (
	pkgres "cms-server/infrastructure/service/response"

	"github.com/gofiber/fiber/v2"
)

func (h *authHandlerImpl) CheckToken(c *fiber.Ctx) error {
	token := c.Query("token", "")
	if token == "" {
		err := pkgres.NewErr("Dữ liệu không hợp lệ").BadReq()
		return h.log.Log(c, err)
	}

	ok, err := h.checkTokenUc.CheckToken(token)
	if err != nil {
		err := pkgres.Err(err).BadReq()
		return h.log.Log(c, err)
	}

	if !ok {
		err := pkgres.NewErr("Phiên làm việc không hợp lệ").BadReq()
		return h.log.Log(c, err)
	}
	return c.JSON(pkgres.NewRes("Phiên làm việc hợp lệ"))
}
