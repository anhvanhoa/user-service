package user

import (
	"user-service/domain/entity"
	"user-service/domain/repository"

	"github.com/anhvanhoa/service-core/common"
	"github.com/anhvanhoa/service-core/utils"
)

type GetUsersUsecase interface {
	Excute(pagination *common.Pagination, filter *entity.FilterUser) (*common.PaginationResult[entity.User], error)
	ExtractUserIds(users []entity.User) ([]string, error)
}

type getUsersUsecase struct {
	userRepo repository.UserRepository
	helper   utils.Helper
}

func NewGetUsersUsecase(userRepo repository.UserRepository, helper utils.Helper) GetUsersUsecase {
	return &getUsersUsecase{
		userRepo: userRepo,
		helper:   helper,
	}
}

func (g *getUsersUsecase) Excute(pagination *common.Pagination, filter *entity.FilterUser) (*common.PaginationResult[entity.User], error) {
	users, total, err := g.userRepo.GetUsers(pagination, filter)
	if err != nil {
		return nil, ErrGetUsers
	}
	totalPages := g.helper.CalculateTotalPages(int64(total), int64(pagination.PageSize))
	return &common.PaginationResult[entity.User]{
		Data:       users,
		Total:      int64(total),
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (g *getUsersUsecase) ExtractUserIds(users []entity.User) ([]string, error) {
	userIds := []string{}
	for _, user := range users {
		userIds = append(userIds, user.ID)
	}
	return userIds, nil
}
