package store

import (
	"context"
	"database/sql"
	"errors"

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

func (a *AdiminStore) GetByEmail(ctx context.Context, email string) (*model.Admin, error) {
	admin := &model.Admin{}
	query := `SELECT id,email,password_hash from admins where email=$1 AND is_active=true`
	err := a.db.QueryRowContext(ctx, query, email).Scan(&admin.Id, &admin.Email, &admin.Password.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return admin, nil
}

func (a *AdiminStore) GetById(ctx context.Context, id string) (*model.Admin, error) {
	admin := &model.Admin{}
	query := `SELECT id,email,password_hash from admins where id=$1 `
	err := a.db.QueryRowContext(ctx, query, id).Scan(&admin.Id, &admin.Email, &admin.Password.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return admin, nil
}
