package user

import (
	"user-service/domain/entity"
	"user-service/domain/repository"

	hashpass "github.com/anhvanhoa/service-core/domain/hash_pass"
)

type CreateUserUsecase interface {
	Excute(data *entity.User) (entity.User, error)
}

type createUserUsecase struct {
	userRepo    repository.UserRepository
	hashService hashpass.HashPassI
}

func NewCreateUserUsecase(userRepo repository.UserRepository, hashService hashpass.HashPassI) CreateUserUsecase {
	return &createUserUsecase{
		userRepo:    userRepo,
		hashService: hashService,
	}
}

func (c *createUserUsecase) Excute(data *entity.User) (entity.User, error) {
	if data == nil {
		return entity.User{}, ErrCreateUser
	}

	isExist, err := c.userRepo.CheckUserExist(data.Email)
	if err != nil {
		return entity.User{}, ErrCreateUser
	}
	if isExist {
		return entity.User{}, ErrUserAlreadyExists
	}

	hashedPassword, err := c.hashService.HashPassword(data.Password)
	if err != nil {
		return entity.User{}, ErrCreateUser
	}
	data.Password = hashedPassword

	user, err := c.userRepo.CreateUser(*data)
	if err != nil {
		return entity.User{}, ErrCreateUser
	}

	return user, nil
}
