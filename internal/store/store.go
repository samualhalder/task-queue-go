package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/samualhalder/task-queue-go/internal/model"
)

type TaskRepository interface{
	GetById(context.Context,int64) (*model.Task,error)
}

type store struct{
	db *sql.DB
	Task TaskRepository
}

func NewStore(db *sql.DB) *store{
	return &store{
		db: db,
		Task: &TaskStore{db},
	}
}

func (s *store) ExecTx(ctx context.Context, fn func(*store) error) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txStore := &store{
		db:    s.db, 
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