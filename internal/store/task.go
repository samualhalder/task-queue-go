package store

import (
	"context"

	"github.com/samualhalder/task-queue-go/internal/dto"
	"github.com/samualhalder/task-queue-go/internal/model"
)


type TaskStore struct{
	db DBTX
}

func(t *TaskStore) GetById(ctx context.Context,id int64) (*model.Task,error){
	print("task id ", id)
	return nil,nil
}
	

func(t *TaskStore) Create(ctx context.Context,task *dto.TaskResponse) error{
	query:=`INSERT INTO tasks(task_name,payload,scheduled_at) VALUES($1,$2, COALESCE($3,now())) RETURNING id,status,priority,scheduled_at,created_at`
    
	err:=t.db.QueryRowContext(ctx,query,task.TaskName,task.Payload,task.ScheduledAt).
			Scan(&task.ID,&task.Status,&task.Priority,&task.ScheduledAt,&task.CreatedAt)
	if err!=nil{
		return err
	}
	return nil
}