package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DTO
type TaskResponse struct {
	ID          uuid.UUID       `json:"id"`
	TaskName    string          `json:"task_name"`
	Payload     json.RawMessage `json:"payload"`
	Status      string          `json:"status"`
	Priority    int             `json:"priority"`
	ScheduledAt *time.Time       `json:"scheduled_at"`
	CreatedAt   time.Time       `json:"created_at"`
}