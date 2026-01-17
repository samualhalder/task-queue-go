package taskhandler

import (
	"context"

	"github.com/samualhalder/task-queue-go/internal/model"
)

type TaskHandler interface{
	Handle(context.Context,*model.Task) error
}