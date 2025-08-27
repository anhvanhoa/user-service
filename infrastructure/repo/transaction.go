package repo

import (
	"auth-service/domain/repository"
	loggerI "auth-service/domain/service/logger"
	"auth-service/domain/service/saga"
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
)

type managerTransaction struct {
	db          *pg.DB
	sagaManager saga.SagaManager
	logger      loggerI.Log
}

var ErrTxContextKey error = errors.New("no transaction in context")

func NewManagerTransaction(db *pg.DB, logger loggerI.Log) repository.ManagerTransaction {
	return &managerTransaction{
		db:          db,
		sagaManager: saga.NewSagaManager(logger),
		logger:      logger,
	}
}

func (mt *managerTransaction) RunInTransaction(fn func(ctx context.Context) error) error {
	tx, err := mt.db.BeginContext(mt.db.Context())
	if err != nil {
		return err
	}
	txCtx := context.WithValue(mt.db.Context(), txContextKey{}, tx)
	if err := fn(txCtx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (mt *managerTransaction) Begin() (context.Context, error) {
	ctx := mt.db.Context()
	tx, err := mt.db.BeginContext(ctx)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, txContextKey{}, tx), nil
}

func (mt *managerTransaction) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(txContextKey{}).(*pg.Tx)
	if !ok {
		return ErrTxContextKey
	}
	return tx.Commit()
}

func (mt *managerTransaction) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(txContextKey{}).(*pg.Tx)
	if !ok {
		return ErrTxContextKey
	}
	return tx.Rollback()
}

func (mt *managerTransaction) RunSagaTransaction(sagaID string, setupSteps func(ctx context.Context, sagaTx *saga.SagaTransaction) error) error {
	sagaTx := mt.sagaManager.NewTransaction(sagaID, mt.db.Context())
	if err := setupSteps(sagaTx.Context, sagaTx); err != nil {
		mt.logger.Error("Failed to setup saga steps: " + err.Error())
		return err
	}
	return sagaTx.Execute()
}

func (mt *managerTransaction) GetSagaManager() saga.SagaManager {
	return mt.sagaManager
}
