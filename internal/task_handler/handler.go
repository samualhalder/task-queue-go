package taskhandler

import (
	"context"

	taskerrors "github.com/samualhalder/task-queue-go/internal/errors"
	"github.com/samualhalder/task-queue-go/internal/model"
)

type TaskHandler interface {
	Handle(context.Context, *model.Task) *taskerrors.TaskError
}
