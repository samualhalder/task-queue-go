package worker

import (
	"context"

	"github.com/samualhalder/task-queue-go/internal/model"
	"go.uber.org/zap"
)

type Executor interface{
	Execute(ctx context.Context, task *model.Task) error
}

type WorkExec struct{
	logger *zap.SugaredLogger
}

func NewWorkerExec(logger zap.SugaredLogger) *WorkExec{
	return &WorkExec{
		logger: &logger,
	}
}

func(w *WorkExec) Execute(ctx context.Context,task *model.Task) error{
	w.logger.Info("executor is executed for", task.ID)
	return nil
}