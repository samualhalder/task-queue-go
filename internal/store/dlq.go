package store

import (
	"context"
	"fmt"

	"github.com/samualhalder/task-queue-go/internal/model"
)

type DlqStore struct {
	DB DBTX
}

func (d *DlqStore) Push(ctx context.Context, task *model.Task, errorType string) error {
	fmt.Print("hti store", task.LastError)
	query := `
	INSERT INTO tasks_dlq (
		task_id,
		error_text,
		error_type,
		attempts,
		payload
	)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (task_id) DO NOTHING
	`

	_, err := d.DB.ExecContext(ctx, query, task.ID, task.LastError, errorType, task.Attempts, task.Payload)
	if err != nil {
		fmt.Print("err ", err.Error())
		return err
	}
	return nil
}
