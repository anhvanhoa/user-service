package session

import (
	"context"
	"user-service/domain/entity"
	"user-service/domain/repository"

	"github.com/anhvanhoa/service-core/domain/cache"
)

type DeleteSessionByTypeAndUserUsecase interface {
	DeleteTokenDb(ctx context.Context, sessionType entity.SessionType, userID string) error
	Excute(ctx context.Context, sessionType entity.SessionType, userID string) error
}

type deleteSessionByTypeAndUserUsecase struct {
	sessionRepo repository.SessionRepository
	cache       cache.CacheI
}

func NewDeleteSessionByTypeAndUserUsecase(sessionRepo repository.SessionRepository, cache cache.CacheI) DeleteSessionByTypeAndUserUsecase {
	return &deleteSessionByTypeAndUserUsecase{
		sessionRepo: sessionRepo,
		cache:       cache,
	}
}

func (d *deleteSessionByTypeAndUserUsecase) DeleteTokenDb(ctx context.Context, sessionType entity.SessionType, userID string) error {
	return d.sessionRepo.DeleteSessionByTypeAndUserID(ctx, sessionType, userID)
}

func (d *deleteSessionByTypeAndUserUsecase) Excute(ctx context.Context, sessionType entity.SessionType, userID string) error {
	tokens, err := d.sessionRepo.GetTokensByTypeAndUserID(ctx, sessionType, userID)
	if err != nil {
		return err
	}
	if len(tokens) != 0 && sessionType == entity.SessionTypeAuthZ {
		go func() {
			for _, token := range tokens {
				d.cache.Delete(token)
			}
		}()
	}
	return d.sessionRepo.DeleteSessionByTypeAndUserID(ctx, sessionType, userID)
}
