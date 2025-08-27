package repository

import (
	"auth-service/domain/service/saga"
	"context"
)

type ManagerTransaction interface {
	RunInTransaction(fn func(ctx context.Context) error) error
	Begin() (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	RunSagaTransaction(sagaID string, setupSteps func(ctx context.Context, sagaTx *saga.SagaTransaction) error) error
	GetSagaManager() saga.SagaManager
}
