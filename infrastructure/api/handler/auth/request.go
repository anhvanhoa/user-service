package handler

import authUC "cms-server/domain/usecase/auth"

type resetPasswordByTokenReq struct {
	Token           string `validate:"required"`
	Password        string `validate:"required"`
	ConfirmPassword string `validate:"required"`
}

type resetPasswordByCodeReq struct {
	Code            string `validate:"required"`
	Email           string `validate:"required,email"`
	Password        string `validate:"required,min=6"`
	ConfirmPassword string `validate:"required,min=6,eqfield=Password"`
}

type checkCodeReq struct {
	Code  string `validate:"required"`
	Email string `validate:"required"`
}

type registerReq struct {
	Email           string `validate:"required,email"`
	FullName        string `validate:"required"`
	Password        string `validate:"required,min=6"`
	ConfirmPassword string `validate:"required,min=6,eqfield=Password"`
	Code            string
}

type loginReq struct {
	Identifier string `validate:"required,email_or_tell=vi"`
	Password   string `validate:"required,min=6"`
}

type forgotPasswordReq struct {
	Email string                    `validate:"required,email"`
	Type  authUC.ForgotPasswordType `validate:"required,in=ForgotByCode ForgotByToken"`
}
