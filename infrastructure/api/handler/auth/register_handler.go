package handler

import (
	"auth-service/constants"
	"auth-service/domain/usecase"
	pkgres "auth-service/infrastructure/service/response"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
)

func (rh *authHandlerImpl) Register(c *fiber.Ctx) error {
	var body registerReq
	if err := c.BodyParser(&body); err != nil {
		err := pkgres.NewErr("Dữ liệu không hợp lệ").BadReq()
		return rh.log.Log(c, err)
	}

	if err := rh.validate.ValidateStruct(body); err != nil {
		err := pkgres.NewErr("Dữ liệu không hợp lệ").BadReq()
		return rh.log.Log(c, err)
	}

	if body.Password != body.ConfirmPassword {
		err := pkgres.NewErr("Mật khẩu không khớp").BadReq()
		return rh.log.Log(c, err)
	}

	if u, err := rh.registerUc.CheckUserExist(body.Email); err != nil && err != pg.ErrNoRows {
		return rh.log.Log(c, err)
	} else if u.ID != "" && u.Veryfied != nil {
		err := pkgres.NewErr("Tài khoản đã tồn tại, vui lòng thử lại !").Code(fiber.StatusBadRequest)
		return rh.log.Log(c, err)
	}

	expAt := time.Now().Add(time.Second * constants.VerifyExpiredAt)
	body.Code = rh.registerUc.GengerateCode(6)
	os := c.Get("User-Agent")
	dataRegister := usecase.RegisterReq{
		Email:           body.Email,
		FullName:        body.FullName,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
		Code:            body.Code,
	}
	res, err := rh.registerUc.Register(dataRegister, os, expAt)
	if pg.ErrNoRows == err {
		err := pkgres.NewErr("Không tìm thấy mẫu email").NotFound()
		return rh.log.Log(c, err)
	} else if err != nil {
		return rh.log.Log(c, err)
	}

	err = rh.registerUc.SendMail(res.UserInfor, rh.env.FRONTEND_URL+"/auth/verify/"+res.Token)
	if pg.ErrNoRows == err {
		err := pkgres.NewErr("Không tìm thấy mẫu email").NotFound()
		return rh.log.Log(c, err)
	} else if err != nil {
		return rh.log.Log(c, err)
	}
	return c.JSON(pkgres.ResData(res.UserInfor).SetMessage("Đăng ký thành công"))
}
