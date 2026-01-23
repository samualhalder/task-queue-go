package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Dlq struct {
	Id        uuid.UUID       `json:"id"`
	TaskId    uuid.UUID       `json:"task_id"`
	Payload   json.RawMessage `json:"payload"`
	Attempts  int             `json:"attempts"`
	ErrorText string          `json:"error_text"`
	ErrorType string          `json:"error_type"`
	FailedAt  time.Time       `json:"failed_at"`
}
