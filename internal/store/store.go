package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samualhalder/task-queue-go/internal/dto"
	"github.com/samualhalder/task-queue-go/internal/model"
)

type TaskRepository interface {
	GetById(context.Context, uuid.UUID) (*model.Task, error)
	Create(context.Context, *dto.TaskResponse) error
	FailedExecution(context.Context, uuid.UUID) (bool, error)
	SuccessfullExecution(ctx context.Context, task *model.Task) error
	ClaimTask(ctx context.Context, taskId uuid.UUID, workerId string) (bool, error)
	FetchStuckTasks(context.Context, time.Duration) ([]uuid.UUID, error)
	GetAndClaimEligibleTask(ctx context.Context, workerID string) (*model.Task, error)
	HandleFailure(ctx context.Context, task *model.Task, errorType string) error
	HandleSuccess(ctx context.Context, task *model.Task) error
}

type Store struct {
	db   *sql.DB
	Task TaskRepository
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:   db,
		Task: &TaskStore{db},
	}
}

func (s *Store) ExecTx(ctx context.Context, fn func(*Store) error) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txStore := &Store{
		db:   s.db,
		Task: &TaskStore{db: tx},
	}

	if err := fn(txStore); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
