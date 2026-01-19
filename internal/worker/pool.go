package worker

import (
	"context"

	"github.com/samualhalder/task-queue-go/internal/queue"
	"github.com/samualhalder/task-queue-go/internal/store"
	"go.uber.org/zap"
)


type Pool struct{
	WorkerCount int
	Workers []*Worker
	logger *zap.SugaredLogger
	queue queue.Queue
    store store.Store
    exec Executor
}


func NewPool(workerCount int,logger *zap.SugaredLogger,queue queue.Queue,store store.Store,exec Executor) *Pool{
	return &Pool{
		WorkerCount: workerCount,
		logger: logger,
		queue: queue,
		store: store,
		exec: exec,
	}
}

func(p *Pool) Start(ctx context.Context){
	p.Workers = make([]*Worker,0,p.WorkerCount)
	for i:=0 ; i<p.WorkerCount;i++{
		worker:=NewWorker(i,p.queue,p.store,p.exec,p.logger)
		
		p.Workers = append(p.Workers, worker)
		go worker.Start(ctx)
	}
}