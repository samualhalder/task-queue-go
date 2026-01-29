package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samualhalder/task-queue-go/internal/dto"
	"github.com/samualhalder/task-queue-go/internal/model"
	"github.com/samualhalder/task-queue-go/internal/payloads"
	"github.com/samualhalder/task-queue-go/internal/retries"
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
	query := `INSERT INTO tasks(task_name,payload,scheduled_at,max_attempts) VALUES($1,$2, COALESCE($3,now()),$4) RETURNING id,status,priority,scheduled_at,created_at`

	err := t.db.QueryRowContext(ctx, query, task.TaskName, task.Payload, task.ScheduledAt, task.Max_Attempts).
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

func (t *TaskStore) ClaimTask(ctx context.Context, task *model.Task, workerID string) error {
	_, err := t.db.ExecContext(ctx, `
		UPDATE tasks
		SET status = 'running',
		    locked_by = $1,
		    locked_at = NOW(),
		    updated_at = NOW()
		WHERE id = $2
	`, workerID, task.ID)

	if err != nil {
		return err
	}

	return nil
}

func (t *TaskStore) FetchStuckTasks(ctx context.Context, timeout time.Duration) ([]uuid.UUID, error) {
	timeDuration := time.Now().Add(timeout)

	query := `UPDATE tasks
    SET
        attempts = attempts + 1,
        locked_by = NULL,
        locked_at = NULL,
		next_attempt=now(),
		last_error='worker crashed',
        status = CASE
            WHEN attempts + 1 >= max_attempts THEN 'failed'
            ELSE 'pending'
        END
    WHERE status = 'running'
      AND locked_at < $1
    RETURNING id;
    `

	rows, err := t.db.QueryContext(ctx, query, timeDuration)
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
			updated_at,
			last_error
		FROM tasks
		WHERE status= 'pending'
		  AND next_attempt <= NOW()
		  AND scheduled_at<=NOW()
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
		&task.LastError,
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
		SET status = 'running',
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

func (t *TaskStore) HandleFailure(ctx context.Context, task *model.Task, errorType string) error {
	next_attempt_time := retries.Interval(task.Attempts)
	next_attempt := time.Now().Add(next_attempt_time)

	query := `UPDATE tasks
				SET
				   attempts=attempts+1,
				   locked_by=NULL,
				   locked_at=NULL,
				   last_error=$3,
				   next_attempt= $2,
				   status = CASE
				   		WHEN attempts+1>=max_attempts THEN 'failed'
						ELSE 'pending'
						END
				WHERE
					id=$1`

	row := t.db.QueryRowContext(ctx, query, task.ID, next_attempt, errorType)
	if row.Err() == sql.ErrNoRows {
		return row.Err()
	}

	return nil
}
func (t *TaskStore) HandleSuccess(ctx context.Context, task *model.Task) error {
	query := `UPDATE tasks
				SET
				   locked_by=NULL,
				   locked_at=NULL,
				   status = 'completed'
				WHERE
					id=$1`

	_, err := t.db.ExecContext(ctx, query, task.ID)
	if err != nil {
		return err
	}

	return nil
}

func (t *TaskStore) HandleFailed(ctx context.Context, task *model.Task, errorType string) error {
	query := `UPDATE tasks
				SET
				   locked_by=NULL,
				   locked_at=NULL,
				   last_error=$2,
				   status='failed'
				WHERE
					id=$1`

	row := t.db.QueryRowContext(ctx, query, task.ID, errorType)
	if row.Err() == sql.ErrNoRows {
		return row.Err()
	}

	return nil
}

func (t *TaskStore) FetchEligibleTask(ctx context.Context) (*model.Task, error) {
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
			updated_at,
			last_error
		FROM tasks
		WHERE status= 'pending'
		  AND next_attempt <= NOW()
		  AND scheduled_at<=NOW()
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
		&task.LastError,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return task, nil
}

func (t *TaskStore) GetTasks(ctx context.Context, params payloads.TaskQueryPayload) ([]model.Task, error) {
	query := `
		SELECT
			id, task_name,payload, status, attempts, max_attempts,
			locked_by, locked_at, created_at, updated_at,last_error,scheduled_at
		FROM tasks
		WHERE
		($1::text IS NULL OR task_name = $1::text)
AND ($2::text IS NULL OR status = $2::text)
AND ($3::timestamptz IS NULL OR created_at <= $3::timestamptz)
AND ($4::timestamptz IS NULL OR created_at >= $4::timestamptz)
ORDER BY created_at DESC
LIMIT $5 OFFSET $6
        `
	tasks := []model.Task{}
	row, err := t.db.QueryContext(ctx, query, params.TaskName, params.Status, params.CreatedBefore, params.CreatedAfter, params.Limit, params.Skip)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tasks, nil
		} else {
			return nil, err
		}
	}

	for row.Next() {
		var task model.Task
		err := row.Scan(&task.ID,
			&task.TaskName,
			&task.Payload,
			&task.Status,
			&task.Attempts,
			&task.MaxAttempts,
			&task.LockedBy,
			&task.LockedAt,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.LastError,
			&task.ScheduledAt,
		)
		if err != nil {
			fmt.Print("error:", err.Error())
			return nil, err
		}
		tasks = append(tasks, task)
		fmt.Print("hit here")
	}
	return tasks, nil
}

func (t *TaskStore) MakeTaskEligible(ctx context.Context, uuid string) error {
	query := ` UPDATE tasks
			SET next_attempt=now(),
				scheduled_at=now(),
				status='pending'
			WHERE id=$1`
	_, err := t.db.ExecContext(ctx, query, uuid)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("No record found")
	}
	return nil
}

func (t *TaskStore) MakeTaskFailed(ctx context.Context, uuid string) error {
	query := ` UPDATE tasks
			SET status='failed'
			WHERE id=$1`
	_, err := t.db.ExecContext(ctx, query, uuid)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("No record found")
	}
	return nil
}
