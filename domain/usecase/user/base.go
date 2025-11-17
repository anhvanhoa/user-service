package user

type UserUsecase struct {
	CreateUserUsecase CreateUserUsecase
	DeleteUserUsecase DeleteUserUsecase
	GetUserUsecase    GetUserUsecase
	GetUsersUsecase   GetUsersUsecase
	UpdateUserUsecase UpdateUserUsecase
	LockUserUsecase   LockUserUsecase
	UnlockUserUsecase UnlockUserUsecase
}

func NewUserUsecase(
	CreateUserUsecase CreateUserUsecase,
	DeleteUserUsecase DeleteUserUsecase,
	GetUserUsecase GetUserUsecase,
	GetUsersUsecase GetUsersUsecase,
	UpdateUserUsecase UpdateUserUsecase,
	LockUserUsecase LockUserUsecase,
	UnlockUserUsecase UnlockUserUsecase,
) *UserUsecase {
	return &UserUsecase{
		CreateUserUsecase: CreateUserUsecase,
		DeleteUserUsecase: DeleteUserUsecase,
		GetUserUsecase:    GetUserUsecase,
		GetUsersUsecase:   GetUsersUsecase,
		UpdateUserUsecase: UpdateUserUsecase,
		LockUserUsecase:   LockUserUsecase,
		UnlockUserUsecase: UnlockUserUsecase,
	}
}
