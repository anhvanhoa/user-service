package handler

import (
	"cms-server/constants"
	pkgres "cms-server/infrastructure/service/response"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (rh *authHandlerImpl) Refresh(c *fiber.Ctx) error {
	refresh := c.Cookies(constants.KeyCookieRefreshToken, "")
	session, err := rh.refreshUc.GetSessionByToken(refresh)
	if err != nil {
		err := pkgres.NewErr("Phiên làm việc không hợp lệ, hãy đăng nhập lại").Unauthorized()
		return rh.log.Log(c, err)
	}

	go rh.refreshUc.ClearSessionExpired()

	claims, err := rh.refreshUc.VerifyToken(refresh)
	if err != nil {
		err := pkgres.NewErr("Phiên làm việc không hợp lệ, hãy đăng nhập lại").Unauthorized()
		return rh.log.Log(c, err)
	}
	expAccess := time.Now().Add(constants.AccessExpiredAt * time.Second)
	access, err := rh.refreshUc.GengerateAccessToken(session.UserID, claims.FullName, expAccess)
	if err != nil {
		err := pkgres.NewErr("Không thể tạo token mới").InternalServerError()
		return rh.log.Log(c, err)
	}
	expRefresh := time.Now().Add(constants.RefreshExpiredAt * time.Second)
	os := c.Get("User-Agent")
	refreshToken, err := rh.refreshUc.GengerateRefreshToken(session.UserID, claims.FullName, expRefresh, os)
	if err != nil {
		err := pkgres.NewErr("Không thể tạo token mới").InternalServerError()
		return rh.log.Log(c, err)
	}
	expR := time.Now().Add(constants.RefreshExpiredAt * time.Second)
	expA := time.Now().Add(constants.AccessExpiredAt * time.Second)
	c.Cookie(&fiber.Cookie{
		Name:     constants.KeyCookieRefreshToken,
		Value:    refreshToken,
		Path:     "/",
		Domain:   rh.env.HOST_APP,
		Secure:   rh.env.IsProduction(),
		HTTPOnly: true,
		Expires:  expR,
	})
	c.Cookie(&fiber.Cookie{
		Name:     constants.KeyCookieAccessToken,
		Value:    access,
		Path:     "/",
		Domain:   rh.env.HOST_APP,
		Secure:   rh.env.IsProduction(),
		HTTPOnly: true,
		Expires:  expA,
	})
	return c.JSON(pkgres.ResData(nil).SetMessage("Làm mới phiên làm việc thành công"))
}
