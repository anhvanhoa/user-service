package handler

import (
	"cms-server/constants"
	pkgres "cms-server/infrastructure/service/response"
	"errors"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
)

func (lh *authHandlerImpl) Login(c *fiber.Ctx) error {
	var body loginReq

	if err := c.BodyParser(&body); err != nil {
		err := pkgres.NewErr("Dữ liệu không hợp lệ").BadReq()
		return lh.log.Log(c, err)
	}

	if err := lh.validate.ValidateStruct(&body); err != nil {
		err := pkgres.Err(err).BadReq().SetData(err.Data)
		return lh.log.Log(c, err)
	}

	user, err := lh.loginUc.GetUserByEmailOrPhone(body.Identifier)
	if errors.Is(err, pg.ErrNoRows) {
		err := pkgres.NewErr("Tài khoản không hợp lệ, hãy kiểm tra lại").Code(fiber.StatusBadRequest)
		return lh.log.Log(c, err)
	}

	if err != nil {
		return lh.log.Log(c, err)
	}

	if !lh.loginUc.CheckHashPassword(body.Password, user.Password) {
		err := pkgres.NewErr("Mật khẩu không chính xác, hãy kiểm tra lại").Code(fiber.StatusBadRequest)
		return lh.log.Log(c, err)
	}
	timeAccess := time.Now().Add(constants.AccessExpiredAt * time.Second)
	access, _ := lh.loginUc.GengerateAccessToken(user.ID, user.FullName, timeAccess)
	timeRefresh := time.Now().Add(constants.RefreshExpiredAt * time.Second)
	os := c.Get("User-Agent")
	refresh, _ := lh.loginUc.GengerateRefreshToken(user.ID, user.FullName, timeRefresh, os)

	c.Cookie(&fiber.Cookie{
		Name:     constants.KeyCookieAccessToken,
		Value:    access,
		Path:     "/",
		Domain:   lh.env.HOST_APP,
		Secure:   lh.env.IsProduction(),
		HTTPOnly: true,
		Expires:  timeAccess,
	})
	c.Cookie(&fiber.Cookie{
		Name:     constants.KeyCookieRefreshToken,
		Value:    refresh,
		Path:     "/",
		Domain:   lh.env.HOST_APP,
		Secure:   lh.env.IsProduction(),
		HTTPOnly: true,
		Expires:  timeRefresh,
	})

	return c.JSON(pkgres.ResData(user.GetInfor()).SetMessage("Đăng nhập thành công"))
}
