package worker

import (
	"context"
	"fmt"
	"math/rand"

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
	num:=rand.Intn(2)
	if(num==0){
		return fmt.Errorf("failed to execute task %d",task.ID)
	}
	return nil
}