package worker

import (
	"context"

	taskerrors "github.com/samualhalder/task-queue-go/internal/errors"
	"github.com/samualhalder/task-queue-go/internal/model"
)

type Executor interface {
	Execute(ctx context.Context, task *model.Task) *taskerrors.TaskError
}
