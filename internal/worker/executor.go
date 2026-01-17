package worker

import (
	"context"

	"github.com/samualhalder/task-queue-go/internal/model"
)

type Executor interface{
	Execute(ctx context.Context, task *model.Task) error
}
