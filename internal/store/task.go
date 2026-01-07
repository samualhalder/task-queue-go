package store

import (
	"context"

	"github.com/samualhalder/task-queue-go/internal/model"
)


type TaskStore struct{
	db DBTX
}

func(t *TaskStore) GetById(ctx context.Context,id int64) (*model.Task,error){
	print("task id ", id)
	return nil,nil
}