package session

import (
	"context"
	"user-service/domain/entity"
	"user-service/domain/repository"

	"github.com/anhvanhoa/service-core/domain/cache"
)

type DeleteSessionByTypeAndTokenUsecase interface {
	Excute(ctx context.Context, sessionType entity.SessionType, token string) error
}

type deleteSessionByTypeAndTokenUsecase struct {
	sessionRepo repository.SessionRepository
	cache       cache.CacheI
}

func NewDeleteSessionByTypeAndTokenUsecase(sessionRepo repository.SessionRepository, cache cache.CacheI) DeleteSessionByTypeAndTokenUsecase {
	return &deleteSessionByTypeAndTokenUsecase{
		sessionRepo: sessionRepo,
		cache:       cache,
	}
}

func (d *deleteSessionByTypeAndTokenUsecase) Excute(ctx context.Context, sessionType entity.SessionType, token string) error {
	go func() {
		d.cache.Delete(token)
	}()
	return d.sessionRepo.DeleteSessionByTypeAndToken(ctx, sessionType, token)
}
