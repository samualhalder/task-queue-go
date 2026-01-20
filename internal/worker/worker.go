package worker

import (
	"context"
	"strconv"
	"time"

	"github.com/samualhalder/task-queue-go/internal/model"
	"github.com/samualhalder/task-queue-go/internal/queue"
	"github.com/samualhalder/task-queue-go/internal/store"
	"go.uber.org/zap"
)

type Worker struct {
	id       int
	queue    queue.Queue
	store    store.Store
	executor Executor
	logger   *zap.SugaredLogger
	interval time.Duration
}

func NewWorker(
	id int,
	queue queue.Queue,
	store store.Store,
	exec Executor,
	logger *zap.SugaredLogger,
) *Worker {
	return &Worker{
		id:       id,
		queue:    queue,
		store:    store,
		executor: exec,
		logger:   logger,
	}
}

func (w *Worker) Start(ctx context.Context) {
	w.logger.Infof("Worker %d is started ", w.id)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			w.logger.Infof("worker %d is stopped ", w.id)
			return
		case <-ticker.C:
			w.processOnce(ctx)
		default:
			_, err := w.queue.Pop(ctx)
			if err != nil {
				continue
			}
			w.processOnce(ctx)
		}
	}
}

func (w *Worker) processOnce(ctx context.Context) {
	w.logger.Infof("worker %d is processing a task", w.id)

	var task *model.Task

	err := w.store.ExecTx(ctx, func(txStore *store.Store) error {
		var err error
		task, err = txStore.Task.GetAndClaimEligibleTask(ctx, strconv.Itoa(w.id))

		return err
	})
	if err != nil || task == nil {
		return
	}

	if err := w.executor.Execute(ctx, task); err != nil {

		w.store.ExecTx(ctx, func(txStore *store.Store) error {
			return w.store.Task.HandleFailure(ctx, task, err.Error())
		})
		return
	}
	w.store.ExecTx(ctx, func(txStore *store.Store) error {
		return w.store.Task.HandleSuccess(ctx, task)
	})
}
