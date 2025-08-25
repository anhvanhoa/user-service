package grpcservice

import (
	"auth-service/bootstrap"
	loggerI "auth-service/domain/service/logger"
	"auth-service/domain/service/queue"
	"auth-service/domain/usecase"
	"auth-service/infrastructure/grpc_client"
	"auth-service/infrastructure/repo"
	argonS "auth-service/infrastructure/service/argon"
	goidS "auth-service/infrastructure/service/goid"
	"auth-service/infrastructure/service/jwt"

	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"

	"github.com/go-pg/pg/v10"
)

type authService struct {
	proto_auth.UnimplementedAuthServiceServer
	env              *bootstrap.Env
	mailService      *grpc_client.MailService
	checkTokenUc     usecase.CheckTokenUsecase
	loginUc          usecase.LoginUsecase
	registerUc       usecase.RegisterUsecase
	refreshUc        usecase.RefreshUsecase
	logoutUc         usecase.LogoutUsecase
	verifyAccountUc  usecase.VerifyAccountUsecase
	forgotPasswordUc usecase.ForgotPasswordUsecase
	resetCodeUc      usecase.ResetPasswordByCodeUsecase
	resetTokenUc     usecase.ResetPasswordByTokenUsecase
	checkCodeUc      usecase.CheckCodeUsecase
}

func NewAuthService(
	db *pg.DB,
	env *bootstrap.Env,
	log loggerI.Log,
	mailService *grpc_client.MailService,
	queueClient queue.QueueClient,
) proto_auth.AuthServiceServer {
	userRepo := repo.NewUserRepository(db)
	sessionRepo := repo.NewSessionRepository(db)
	tx := repo.NewManagerTransaction(db)
	argonService := argonS.NewArgon()
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
		env:          env,
		mailService:  mailService,
		checkTokenUc: usecase.NewCheckTokenUsecase(sessionRepo),
		loginUc: usecase.NewLoginUsecase(
			userRepo,
			sessionRepo,
			jwtAccessService,
			jwtRefreshService,
			argonService,
			cacheService,
		),
		registerUc: usecase.NewRegisterUsecase(
			userRepo,
			sessionRepo,
			jwtRegisterService,
			tx,
			goidService,
			argonService,
			cacheService,
			queueClient,
		),
		refreshUc: usecase.NewRefreshUsecase(
			sessionRepo,
			jwtAccessService,
			jwtRefreshService,
			cacheService,
		),
		logoutUc: usecase.NewLogoutUsecase(
			sessionRepo,
			jwtAccessService,
			cacheService,
		),
		verifyAccountUc: usecase.NewVerifyAccountUsecase(
			userRepo,
			sessionRepo,
			jwtRegisterService,
			cacheService,
		),
		forgotPasswordUc: usecase.NewForgotPasswordUsecase(
			userRepo,
			sessionRepo,
			tx,
			jwtForgotService,
			cacheService,
		),
		resetCodeUc: usecase.NewResetPasswordCodeUsecase(
			userRepo,
			sessionRepo,
			cacheService,
			jwtForgotService,
			argonService,
		),
		resetTokenUc: usecase.NewResetPasswordTokenUsecase(
			userRepo,
			sessionRepo,
			cacheService,
			jwtForgotService,
			argonService,
		),
		checkCodeUc: usecase.NewCheckCodeUsecase(
			userRepo,
			sessionRepo,
		),
	}
}
