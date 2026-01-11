package queue

import (
	"context"

	"github.com/google/uuid"
)

type Queue interface{
	Push(context.Context, uuid.UUID) error
	Pop(context.Context) (uuid.UUID,error)
}

