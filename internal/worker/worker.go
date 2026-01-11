package worker

import (
	"context"

	"github.com/samualhalder/task-queue-go/internal/queue"
	"github.com/samualhalder/task-queue-go/internal/store"
	"go.uber.org/zap"
)

type Worker struct {
    id        int
    queue     queue.Queue
    tasks     store.TaskRepository
    executor  Executor
    logger    *zap.SugaredLogger
}





func NewWorker(
    id int,
    queue queue.Queue,
    tasks store.TaskRepository,
    exec Executor,
    logger *zap.SugaredLogger,
) *Worker {
    return &Worker{
        id:     id,
        queue: queue,
        tasks: tasks,
        executor:  exec,
        logger: logger,
    }
}

func(w *Worker) Start(ctx context.Context) {
	w.logger.Info("Worder %d is started ",w.id)
}