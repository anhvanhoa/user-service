package handler

import (
	"cms-server/constants"
	serviceError "cms-server/domain/service/error"
	pkgres "cms-server/infrastructure/service/response"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func (l *authHandlerImpl) Logout(c *fiber.Ctx) error {
	token := c.Cookies(constants.KeyCookieAccessToken)
	if token == "" {
		err := pkgres.NewErr("Phiên làm việc đã hết hạn, vui lòng đăng nhập").Unauthorized()
		return l.log.Log(c, err)
	}

	if err := l.logoutUc.VerifyToken(token); err != nil {
		if errors.Is(err, serviceError.ErrNotFoundSession) {
			err := pkgres.Err(err).Unauthorized()
			return l.log.Log(c, err)
		}
		err := pkgres.NewErr("Phiên làm việc không hợp lệ").Unauthorized()
		return l.log.Log(c, err)
	}

	if err := l.logoutUc.Logout(token); err != nil {
		err := pkgres.Err(err).InternalServerError()
		return l.log.Log(c, err)
	}

	return c.JSON(pkgres.NewErr("Đăng xuất thành công"))
}
