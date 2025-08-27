package common

import (
	"auth-service/domain/service/saga"
	"context"
)

type ExecuteSaga func(ctx context.Context, sagaTx *saga.SagaTransaction) error
