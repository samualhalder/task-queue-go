package taskexecutor

import (
	"context"
	"fmt"

	taskerrors "github.com/samualhalder/task-queue-go/internal/errors"
	"github.com/samualhalder/task-queue-go/internal/model"
	taskhandler "github.com/samualhalder/task-queue-go/internal/task_handler"
	"go.uber.org/zap"
)

type TaskExecutor struct {
	Handlers map[string]taskhandler.TaskHandler
	Logger   *zap.SugaredLogger
}

func (t *TaskExecutor) Execute(ctx context.Context, task *model.Task) *taskerrors.TaskError {
	handler, ok := t.Handlers[task.TaskName]
	if !ok {
		return taskerrors.NonRetryable(fmt.Errorf("no task of type: %s", task.TaskName))
	}
	return handler.Handle(ctx, task)
}
