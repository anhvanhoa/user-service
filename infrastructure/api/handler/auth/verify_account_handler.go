package handler

import (
	pkgjwt "cms-server/infrastructure/service/jwt"
	pkgres "cms-server/infrastructure/service/response"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
)

func (vah *authHandlerImpl) VerifyAccount(c *fiber.Ctx) error {
	t := c.Params("t")
	if t == "" {
		err := pkgres.NewErr("Dữ liệu không hợp lệ").BadReq()
		return vah.log.Log(c, err)
	}

	claims, err := vah.verifyAccountUc.VerifyRegister(t)
	if err == pkgjwt.ErrParseToken {
		err := pkgres.NewErr("Không thể lấy thông tin")
		return vah.log.Log(c, err)
	} else if err != nil {
		err := pkgres.NewErr("Xác thực không thành công").BadReq()
		return vah.log.Log(c, err)
	}
	user, err := vah.verifyAccountUc.GetUserById(claims.Id)
	if err == pg.ErrNoRows {
		err := pkgres.NewErr("Tài khoản không tồn tại").NotFound()
		return vah.log.Log(c, err)
	} else if err != nil {
		return vah.log.Log(c, err) // internal error
	} else if user.Veryfied != nil {
		err := pkgres.NewErr("Tài khoản đã được xác thực").BadReq()
		return vah.log.Log(c, err)
	} else if user.CodeVerify != claims.Code {
		err := pkgres.NewErr("Mã xác thực không hợp lệ").BadReq()
		return vah.log.Log(c, err)
	}

	if err := vah.verifyAccountUc.VerifyAccount(claims.Id); err != nil {
		return vah.log.Log(c, err)
	}
	res := pkgres.NewRes("Xác thực tài khoản thành công").Code(fiber.StatusOK)
	return c.Status(res.GetCode()).JSON(res)
}
