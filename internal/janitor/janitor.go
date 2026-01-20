package janitor

import (
	"context"
	"time"

	"github.com/samualhalder/task-queue-go/internal/queue"
	"github.com/samualhalder/task-queue-go/internal/store"
	"go.uber.org/zap"
)

type Janitor struct {
	tasks store.TaskRepository
	queue queue.Queue

	interval time.Duration
	timeout  time.Duration
	logger   *zap.SugaredLogger
}

func NewJanitor(
	tasks store.TaskRepository,
	queue queue.Queue,
	interval time.Duration,
	timeout time.Duration,
	logger *zap.SugaredLogger,
) *Janitor {
	return &Janitor{
		tasks:    tasks,
		queue:    queue,
		interval: interval,
		timeout:  timeout,
		logger:   logger,
	}
}

func (j *Janitor) Start(ctx context.Context) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			j.runOnce(ctx)
		}
	}
}

func (j *Janitor) runOnce(ctx context.Context) {

	// TODO: pull task from bd need a method in task to do that
	j.logger.Infof("running janitor")
	ids, err := j.tasks.FetchStuckTasks(ctx, j.timeout)
	if err != nil {
		j.logger.Errorw("failed to fetch stuck tasks from db", "error", err)
		return
	}
	if len(ids) > 0 {
		err := j.queue.Signal(ctx)
		if err != nil {
			j.logger.Errorf("failed to push task in janitor", "error", err)
		}
	}
}
