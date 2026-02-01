package store

import (
	"context"

	"github.com/samualhalder/task-queue-go/internal/model"
)

type AdiminStore struct {
	db DBTX
}

func (a *AdiminStore) Create(ctx context.Context, admin model.Admin) error {
	query := `INSERT INTO admins(email,password_hash) VALUES($1,$2)`
	_, err := a.db.ExecContext(ctx, query, admin.Email, admin.Password.PasswordHash)
	if err != nil {
		return err
	}
	return nil
}
