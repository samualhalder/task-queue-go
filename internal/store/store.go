package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samualhalder/task-queue-go/internal/dto"
	"github.com/samualhalder/task-queue-go/internal/model"
	"github.com/samualhalder/task-queue-go/internal/payloads"
)

type TaskRepository interface {
	GetById(context.Context, uuid.UUID) (*model.Task, error)
	Create(context.Context, *dto.TaskResponse) error
	FailedExecution(context.Context, uuid.UUID) (bool, error)
	SuccessfullExecution(ctx context.Context, task *model.Task) error
	ClaimTask(ctx context.Context, taskId *model.Task, workerId string) error
	FetchStuckTasks(context.Context, time.Duration) ([]uuid.UUID, error)
	GetAndClaimEligibleTask(ctx context.Context, workerID string) (*model.Task, error)
	HandleFailure(ctx context.Context, task *model.Task, errorType string) error
	HandleSuccess(ctx context.Context, task *model.Task) error
	HandleFailed(ctx context.Context, task *model.Task, errorType string) error
	FetchEligibleTask(ctx context.Context) (*model.Task, error)
	GetTasks(ctx context.Context, params payloads.TaskQueryPayload) ([]model.Task, error)
	MakeTaskEligible(ctx context.Context, uuid string) error
	MakeTaskFailed(ctx context.Context, uuid string) error
}
type DlqRepository interface {
	Push(ctx context.Context, task *model.Task, errorType string) error
}

type TaskRateLimitRepository interface {
	AllowTx(
		ctx context.Context,
		task *model.Task,
	) (bool, error)
}

type Store struct {
	db            *sql.DB
	Task          TaskRepository
	Dlq           DlqRepository
	TaskRateLimit TaskRateLimitRepository
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:            db,
		Task:          &TaskStore{db},
		Dlq:           &DlqStore{db},
		TaskRateLimit: &TaskRateLimitStore{db},
	}
}

func (s *Store) ExecTx(ctx context.Context, fn func(*Store) error) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txStore := &Store{
		db:            s.db,
		Task:          &TaskStore{db: tx},
		Dlq:           &DlqStore{DB: tx},
		TaskRateLimit: &TaskRateLimitStore{db: tx},
	}

	if err := fn(txStore); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
