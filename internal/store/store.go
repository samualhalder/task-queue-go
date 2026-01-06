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

type Store struct{
	db *sql.DB
	Task TaskRepository
}

func NewStore(db *sql.DB) *Store{
	return &Store{
		db: db,
		Task: &TaskStore{db},
	}
}

func (s *Store) ExecTx(ctx context.Context, fn func(*Store) error) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txStore := &Store{
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