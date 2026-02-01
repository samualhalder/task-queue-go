package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/samualhalder/task-queue-go/internal/app"
	"github.com/samualhalder/task-queue-go/internal/model"
)

type CreateAdmin struct {
	App app.Application
}

func (c *CreateAdmin) Name() string {
	return "create-admin"
}

func (c *CreateAdmin) Run(args []string) error {
	if len(args) != 2 {
		return errors.New("usage: create-admin <email> <password>")
	}
	admin := model.Admin{
		Email: args[0],
	}
	if err := admin.Password.Set(args[1]); err != nil {
		return err
	}
	if err := c.App.Store.Admin.Create(context.Background(), admin); err != nil {
		return err
	}
	fmt.Println("admin created successfully")
	return nil
}
