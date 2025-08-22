package router

import (
	"cms-server/domain/usecase"
	handler "cms-server/infrastructure/api/handler/auth"
	"cms-server/infrastructure/repo"
	argonS "cms-server/infrastructure/service/argon"
	goidS "cms-server/infrastructure/service/goid"
	pkgjwt "cms-server/infrastructure/service/jwt"
)

func (r *Router) initAuthRouter() {
	authR := r.app.Group("/auth")
	sessionRepo := repo.NewSessionRepository(r.db)
	userRepo := repo.NewUserRepository(r.db)
	tx := repo.NewManagerTransaction(r.db)
	jwtForgot := pkgjwt.NewJWT(r.env.JWT_SECRET.Forgot)
	jwtAccess := pkgjwt.NewJWT(r.env.JWT_SECRET.Access)
	jwtRefresh := pkgjwt.NewJWT(r.env.JWT_SECRET.Refresh)
	jwtVerify := pkgjwt.NewJWT(r.env.JWT_SECRET.Verify)
	argon := argonS.NewArgon()
	goid := goidS.NewGoId()
	h := handler.NewAuthHandler(
		usecase.NewCheckTokenUsecase(sessionRepo),
		usecase.NewCheckCodeUsecase(userRepo, sessionRepo),
		usecase.NewForgotPasswordUsecase(userRepo, sessionRepo, tx, jwtForgot, r.cache),
		usecase.NewLoginUsecase(userRepo, sessionRepo, jwtAccess, jwtRefresh, argon, r.cache),
		usecase.NewLogoutUsecase(sessionRepo, jwtAccess, r.cache),
		usecase.NewRefreshUsecase(sessionRepo, jwtAccess, jwtRefresh, r.cache),
		usecase.NewRegisterUsecase(userRepo, sessionRepo, jwtVerify, tx, goid, argon, r.cache),
		usecase.NewResetPasswordCodeUsecase(userRepo, sessionRepo, r.cache, jwtForgot, argon),
		usecase.NewResetPasswordTokenUsecase(userRepo, sessionRepo, r.cache, jwtForgot, argon),
		usecase.NewVerifyAccountUsecase(userRepo, sessionRepo, jwtVerify, r.cache),
		r.log,
		r.env,
		r.valid,
	)
	authR.Post("/login", h.Login)
	authR.Post("/register", h.Register)
	authR.Post("/verify/:t", h.VerifyAccount)
	authR.Post("/forgot-password", h.Forgot)
	authR.Get("/forgot-password", h.CheckToken)
	authR.Post("/reset-password", h.ResetToken)
	authR.Post("/check-code/forgot-password", h.CheckCode)
	authR.Post("/reset-password/code", h.ResetCode)
	authR.Post("/refresh", h.Refresh)
	authR.Post("/logout", h.Logout)
}
