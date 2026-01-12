package worker

import (
	"context"
	"strconv"
	"time"

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
	w.logger.Infof("Worker %d is started ",w.id)
    for{
        select {
        case <- ctx.Done():
            return 
        default :
          w.processOnce(ctx)
        } 
    }
}

func (w *Worker) processOnce(ctx context.Context) {
	taskID, err := w.queue.Pop(ctx)
	if err != nil {
		time.Sleep(500 * time.Millisecond)
		w.logger.Errorw("queue pop failed", "error", err)
		return
	}

	task, err := w.tasks.GetById(ctx, taskID)
	if err != nil || task == nil {
		return
	}

	claimed, err := w.tasks.ClaimTask(ctx, taskID, strconv.Itoa(w.id))
	if err != nil || !claimed {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			w.logger.Errorw("task panic", "task_id", taskID, "panic", r)
		}
	}()

	if err := w.executor.Execute(ctx, task); err != nil {
		retry, err2 := w.tasks.FailedExecution(ctx, task.ID)
		if err2 != nil {
			w.logger.Errorw("failed to mark task failed", "error", err2)
			return
		}

		if retry {
			_ = w.queue.Push(ctx, taskID)
		}
		return
	}

	if err := w.tasks.SuccessfullExecution(ctx, task); err != nil {
		w.logger.Errorw("failed to mark task completed", "error", err)
	}
}
