package taskexecutor

import (
	"context"
	"fmt"

	"github.com/samualhalder/task-queue-go/internal/model"
	taskhandler "github.com/samualhalder/task-queue-go/internal/task_handler"
	"go.uber.org/zap"
)

type TaskExecutor struct{
 Handlers map[string] taskhandler.TaskHandler
 Logger *zap.SugaredLogger
}

func(t *TaskExecutor) Execute(ctx context.Context,task *model.Task) error{
	handler,ok:=t.Handlers[task.TaskName]
	if !ok{
		return fmt.Errorf("no task of type: %s", task.TaskName)
	}
	return handler.Handle(ctx,task)
}