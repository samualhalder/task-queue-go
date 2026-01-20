package queue

import (
	"context"

	"github.com/google/uuid"
)

type Queue interface {
	Signal(context.Context) error
	Pop(context.Context) (uuid.UUID, error)
}
