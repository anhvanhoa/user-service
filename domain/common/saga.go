package common

import (
	"context"

	"github.com/anhvanhoa/service-core/domain/saga"
)

type ExecuteSaga func(ctx context.Context, sagaTx saga.SagaTransactionI) error
