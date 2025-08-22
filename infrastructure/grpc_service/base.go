package grpcservice

import (
	"cms-server/bootstrap"
	authUC "cms-server/domain/usecase/auth"
	"cms-server/infrastructure/repo"
	argonS "cms-server/infrastructure/service/argon"
	goidS "cms-server/infrastructure/service/goid"
	"cms-server/infrastructure/service/jwt"
	proto "cms-server/proto/gen/auth/v1"

	"github.com/go-pg/pg/v10"
)

type authService struct {
	proto.UnimplementedAuthServiceServer
	checkTokenUc     authUC.CheckTokenUsecase
	loginUc          authUC.LoginUsecase
	registerUc       authUC.RegisterUsecase
	refreshUc        authUC.RefreshUsecase
	logoutUc         authUC.LogoutUsecase
	verifyAccountUc  authUC.VerifyAccountUsecase
	forgotPasswordUc authUC.ForgotPasswordUsecase
	resetCodeUc      authUC.ResetPasswordByCodeUsecase
	resetTokenUc     authUC.ResetPasswordByTokenUsecase
	checkCodeUc      authUC.CheckCodeUsecase
}

func NewAuthService(db *pg.DB, env *bootstrap.Env) proto.AuthServiceServer {
	// Initialize repositories
	userRepo := repo.NewUserRepository(db)
	sessionRepo := repo.NewSessionRepository(db)
	tx := repo.NewManagerTransaction(db)

	// Initialize services
	argonService := argonS.NewArgon()

	// Initialize cache service
	configRedis := bootstrap.NewRedisConfig(
		env.DB_CACHE.Addr,
		env.DB_CACHE.Password,
		env.DB_CACHE.DB,
		env.DB_CACHE.Network,
		env.DB_CACHE.MaxIdle,
		env.DB_CACHE.MaxActive,
		env.DB_CACHE.IdleTimeout,
	)
	cacheService := bootstrap.NewRedis(configRedis)

	goidService := goidS.NewGoId()
	jwtAccessService := jwt.NewJWT(env.JWT_SECRET.Access)
	jwtRefreshService := jwt.NewJWT(env.JWT_SECRET.Refresh)
	jwtRegisterService := jwt.NewJWT(env.JWT_SECRET.Verify)
	jwtForgotService := jwt.NewJWT(env.JWT_SECRET.Forgot)

	return &authService{
		checkTokenUc:     authUC.NewCheckTokenUsecase(sessionRepo),
		loginUc:          authUC.NewLoginUsecase(userRepo, sessionRepo, jwtAccessService, jwtRefreshService, argonService, cacheService),
		registerUc:       authUC.NewRegisterUsecase(userRepo, sessionRepo, jwtRegisterService, tx, goidService, argonService, cacheService),
		refreshUc:        authUC.NewRefreshUsecase(sessionRepo, jwtAccessService, jwtRefreshService, cacheService),
		logoutUc:         authUC.NewLogoutUsecase(sessionRepo, jwtAccessService, cacheService),
		verifyAccountUc:  authUC.NewVerifyAccountUsecase(userRepo, sessionRepo, jwtRegisterService, cacheService),
		forgotPasswordUc: authUC.NewForgotPasswordUsecase(userRepo, sessionRepo, tx, jwtForgotService, cacheService),
		resetCodeUc:      authUC.NewResetPasswordCodeUsecase(userRepo, sessionRepo, cacheService, jwtForgotService, argonService),
		resetTokenUc:     authUC.NewResetPasswordTokenUsecase(userRepo, sessionRepo, cacheService, jwtForgotService, argonService),
		checkCodeUc:      authUC.NewCheckCodeUsecase(userRepo, sessionRepo),
	}
}
