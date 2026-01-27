package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/samualhalder/task-queue-go/internal/model"
)

type TaskRateLimitStore struct {
	db DBTX
}

func (t *TaskRateLimitStore) AllowTx(
	ctx context.Context,
	task *model.Task,
) (bool, error) {

	trl := &model.TaskRateLimit{}

	query := `
		SELECT
			capacity,
			refill_rate,
			tokens,
			last_refill_at
		FROM task_rate_limit
		WHERE key = $1
		FOR UPDATE;
	`

	err := t.db.QueryRowContext(ctx, query, task.TaskName).Scan(
		&trl.Capacity,
		&trl.RefillRate,
		&trl.Tokens,
		&trl.LastRefillAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return false, nil
		}
		return false, err
	}

	now := time.Now().UTC()

	// elapsed time in seconds
	elapsed := now.Sub(trl.LastRefillAt).Seconds()
	if elapsed < 0 {
		elapsed = 0
	}

	refilledTokens := trl.Tokens + elapsed*float64(trl.RefillRate)
	if refilledTokens > float64(trl.Capacity) {
		refilledTokens = float64(trl.Capacity)
	}

	if refilledTokens < 1 {
		return false, nil
	}

	newTokens := refilledTokens - 1

	updateQuery := `
		UPDATE task_rate_limit
		SET
			tokens = $1,
			last_refill_at = $2
		WHERE key = $3;
	`

	_, err = t.db.ExecContext(
		ctx,
		updateQuery,
		newTokens,
		now,
		task.TaskName,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}
