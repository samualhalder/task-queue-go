package worker

import (
	"context"
	"strconv"
	"time"

	taskerrors "github.com/samualhalder/task-queue-go/internal/errors"
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
	ticker := time.NewTicker(time.Second)
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

	defer func() {
		if r := recover(); r != nil {
			w.logger.Errorf("worker panic", "worker_id", w.id, "panic", r)
		}
	}()

	// TODO: generate the thoken here (uuid)
	
	var task *model.Task

	err := w.store.ExecTx(ctx, func(txStore *store.Store) error {
		var err error
		task, err = txStore.Task.GetAndClaimEligibleTask(ctx, strconv.Itoa(w.id))

		return err
	})
	if err != nil || task == nil {
		return
	}
	if w.id == 1 {
		time.Sleep(time.Minute * 5)
	}

	if err := w.executor.Execute(ctx, task); err != nil {

		if err.Type == taskerrors.ErrRetryable {
			w.store.ExecTx(ctx, func(txStore *store.Store) error {
				return w.store.Task.HandleFailure(ctx, task, err.Error())
			})
		} else {
			w.store.ExecTx(ctx, func(txStore *store.Store) error {
				msg := err.Error()
				task.LastError = &msg
				err1 := w.store.Task.HandleFailed(ctx, task, err.Error())
				if err1 != nil {
					return err
				}
				err1 = w.store.Dlq.Push(ctx, task, string(err.Type))
				if err1 != nil {
					return err
				}
				return nil
			})
		}
		return
	}
	w.store.ExecTx(ctx, func(txStore *store.Store) error {
		return w.store.Task.HandleSuccess(ctx, task)
	})
}
