package saga

import (
	loggerI "auth-service/domain/service/logger"
	"context"
	"fmt"
)

// SagaStep định nghĩa một bước trong saga
type SagaStep struct {
	Name       string
	Execute    func(ctx context.Context) error
	Compensate func(ctx context.Context) error
	IsExecuted bool
	IsExternal bool // Đánh dấu step này có call external service không
}

// SagaTransaction quản lý distributed transaction
type SagaTransaction struct {
	ID       string
	Steps    []*SagaStep
	Context  context.Context
	Logger   loggerI.Log
	executed []int // Track các step đã execute
}

// SagaManager interface để quản lý saga transactions
type SagaManager interface {
	NewTransaction(id string, ctx context.Context) *SagaTransaction
}

type sagaManager struct {
	logger loggerI.Log
}

func NewSagaManager(logger loggerI.Log) SagaManager {
	return &sagaManager{
		logger: logger,
	}
}

func (sm *sagaManager) NewTransaction(id string, ctx context.Context) *SagaTransaction {
	return &SagaTransaction{
		ID:       id,
		Steps:    make([]*SagaStep, 0),
		Context:  ctx,
		Logger:   sm.logger,
		executed: make([]int, 0),
	}
}

func (st *SagaTransaction) AddStep(step *SagaStep) *SagaTransaction {
	st.Steps = append(st.Steps, step)
	return st
}

func (st *SagaTransaction) Execute() error {
	st.Logger.Info(fmt.Sprintf("Starting saga transaction: %s", st.ID))

	for i, step := range st.Steps {
		st.Logger.Info(fmt.Sprintf("Executing step %d: %s", i, step.Name))

		if err := step.Execute(st.Context); err != nil {
			st.Logger.Error(fmt.Sprintf("Step %d (%s) failed: %v", i, step.Name, err))

			// Rollback tất cả steps đã executed
			if rollbackErr := st.rollbackExecutedSteps(); rollbackErr != nil {
				st.Logger.Error(fmt.Sprintf("Rollback failed: %v", rollbackErr))
				return fmt.Errorf("step %d failed: %w, và rollback cũng failed: %v", i, err, rollbackErr)
			}

			return fmt.Errorf("step %d (%s) failed: %w", i, step.Name, err)
		}

		step.IsExecuted = true
		st.executed = append(st.executed, i)
		st.Logger.Info(fmt.Sprintf("Step %d (%s) completed successfully", i, step.Name))
	}

	st.Logger.Info(fmt.Sprintf("Saga transaction %s completed successfully", st.ID))
	return nil
}

func (st *SagaTransaction) Rollback() error {
	return st.rollbackExecutedSteps()
}

func (st *SagaTransaction) rollbackExecutedSteps() error {
	st.Logger.Info(fmt.Sprintf("Starting rollback for saga transaction: %s", st.ID))

	// Rollback theo thứ tự ngược lại
	for i := len(st.executed) - 1; i >= 0; i-- {
		stepIndex := st.executed[i]
		step := st.Steps[stepIndex]

		if step.Compensate == nil {
			st.Logger.Warn(fmt.Sprintf("Step %d (%s) không có compensation action", stepIndex, step.Name))
			continue
		}

		st.Logger.Info(fmt.Sprintf("Compensating step %d: %s", stepIndex, step.Name))

		if err := step.Compensate(st.Context); err != nil {
			st.Logger.Error(fmt.Sprintf("Compensation failed for step %d (%s): %v", stepIndex, step.Name, err))
			return fmt.Errorf("compensation failed for step %d (%s): %w", stepIndex, step.Name, err)
		}

		st.Logger.Info(fmt.Sprintf("Step %d (%s) compensated successfully", stepIndex, step.Name))
	}

	st.Logger.Info(fmt.Sprintf("Rollback completed for saga transaction: %s", st.ID))
	return nil
}

func NewStep(name string, execute, compensate func(ctx context.Context) error) *SagaStep {
	return &SagaStep{
		Name:       name,
		Execute:    execute,
		Compensate: compensate,
		IsExternal: false,
	}
}

func NewExternalStep(name string, execute, compensate func(ctx context.Context) error) *SagaStep {
	return &SagaStep{
		Name:       name,
		Execute:    execute,
		Compensate: compensate,
		IsExternal: true,
	}
}
