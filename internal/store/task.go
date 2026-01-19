package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samualhalder/task-queue-go/internal/dto"
	"github.com/samualhalder/task-queue-go/internal/model"
)

type TaskStore struct {
	db DBTX
}

func (t *TaskStore) GetById(ctx context.Context, id uuid.UUID) (*model.Task, error) {

	query := `
	SELECT
		id, task_name,payload, status, attempts, max_attempts,
		locked_by, locked_at, created_at, updated_at
	FROM tasks
	WHERE id = $1
	`

	task := &model.Task{}
	err := t.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.TaskName,
		&task.Payload,
		&task.Status,
		&task.Attempts,
		&task.MaxAttempts,
		&task.LockedBy,
		&task.LockedAt,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return task, nil
}

func (t *TaskStore) Create(ctx context.Context, task *dto.TaskResponse) error {
	query := `INSERT INTO tasks(task_name,payload,scheduled_at) VALUES($1,$2, COALESCE($3,now())) RETURNING id,status,priority,scheduled_at,created_at`

	err := t.db.QueryRowContext(ctx, query, task.TaskName, task.Payload, task.ScheduledAt).
		Scan(&task.ID, &task.Status, &task.Priority, &task.ScheduledAt, &task.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskStore) FailedExecution(ctx context.Context, taskID uuid.UUID) (bool, error) {
	query := `
	UPDATE tasks
	SET
		attempts = attempts + 1,
		locked_by = NULL,
		locked_at = NULL,
		status = CASE
			WHEN attempts + 1 >= max_attempts THEN 'failed'
			ELSE 'pending'
		END
	WHERE id = $1 AND status = 'running'
	RETURNING attempts < max_attempts
	`

	var retry bool
	err := t.db.QueryRowContext(ctx, query, taskID).Scan(&retry)
	return retry, err
}

func (t *TaskStore) SuccessfullExecution(ctx context.Context, task *model.Task) error {
	query := `UPDATE tasks SET attempts= attempts+1,locked_by=NULL,locked_at=NULL,status='completed' WHERE id=$1 AND status='running'`

	_, err := t.db.ExecContext(ctx, query, task.ID)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskStore) ClaimTask(ctx context.Context, taskId uuid.UUID, workerId string) (bool, error) {
	query := `UPDATE tasks SET locked_by=$1,status='running',locked_at=NOW() WHERE id=$2 AND locked_by IS NULL AND status='pending'`
	res, err := t.db.ExecContext(ctx, query, workerId, taskId)
	if err != nil {
		return false, err
	}
	effRows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return effRows == 1, nil
}

func (t *TaskStore) FetchStuckTasks(ctx context.Context, timeout time.Duration) ([]uuid.UUID, error) {
	query := `UPDATE tasks
    SET
        attempts = attempts + 1,
        locked_by = NULL,
        locked_at = NULL,
        status = CASE
            WHEN attempts + 1 >= max_attempts THEN 'failed'
            ELSE 'pending'
        END
    WHERE status = 'running'
      AND locked_at < now() - $1::interval
    RETURNING id;
    `

	rows, err := t.db.QueryContext(ctx, query, timeout.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids = make([]uuid.UUID, 0)
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

func (t *TaskStore) GetAndClaimEligibleTask(
	ctx context.Context,
	workerID string,
) (*model.Task, error) {

	query := `
		SELECT
			id,
			task_name,
			payload,
			status,
			attempts,
			max_attempts,
			next_attempt,
			created_at,
			updated_at
		FROM tasks
		WHERE status= 'pending'
		  AND next_attempt <= NOW()
		ORDER BY next_attempt
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	`

	task := &model.Task{}
	err := t.db.QueryRowContext(ctx, query).Scan(
		&task.ID,
		&task.TaskName,
		&task.Payload,
		&task.Status,
		&task.Attempts,
		&task.MaxAttempts,
		&task.NextAttempt,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Claim immediately (same transaction)
	_, err = t.db.ExecContext(ctx, `
		UPDATE tasks
		SET status = 'processing',
		    locked_by = $1,
		    locked_at = NOW(),
		    updated_at = NOW()
		WHERE id = $2
	`, workerID, task.ID)

	if err != nil {
		return nil, err
	}

	return task, nil
}
