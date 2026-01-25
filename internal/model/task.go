package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID       `json:"id"`
	TaskName    string          `json:"task_name"`
	Payload     json.RawMessage `json:"payload"`
	Status      State           `json:"status"`
	Priority    int             `json:"priority"`
	ScheduledAt time.Time       `json:"scheduled_at"`
	Attempts    int             `json:"attempts"`
	MaxAttempts int             `json:"max_attempts"`
	LastError   *string         `json:"last_error"`
	LockedBy    *string         `json:"locked_by"`
	LockedAt    *time.Time      `json:"locked_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	NextAttempt time.Time       `json:"next_attempt"`
}

// 1	id	uuid	NULL	NULL	NULL	NULL	NO	gen_random_uuid()
// 2	task_name	text	NULL	NULL	NULL	NULL	NO	NULL
// 3	payload	jsonb	NULL	NULL	NULL	NULL	NO	NULL
// 4	status	text	NULL	NULL	NULL	NULL	NO	'pending'::text
// 5	priority	int2	16	NULL	0	NULL	NO	0
// 6	scheduled_at	timestamptz	NULL	6	NULL	NULL	NO	now()
// 7	attempts	int4	32	NULL	0	NULL	NO	0
// 8	max_attempts	int4	32	NULL	0	NULL	NO	3
// 9	last_error	text	NULL	NULL	NULL	NULL	YES	NULL
// 10	locked_by	text	NULL	NULL	NULL	NULL	YES	NULL
// 11	locked_at	timestamptz	NULL	6	NULL	NULL	YES	NULL
// 12	created_at	timestamptz	NULL	6	NULL	NULL	NO	now()
// 13	updated_at	timestamptz	NULL	6	NULL	NULL	NO	now()

type State string

const (
	StatePending   State = "pending"
	StateRunning   State = "running"
	StateFailed    State = "failed"
	StateCompleted State = "completed"
)

var allowedStateChanges = map[State][]State{
	StatePending:   {StateRunning},
	StateRunning:   {StateCompleted, StateFailed},
	StateFailed:    {StatePending},
	StateCompleted: {},
}

// func canStateChange(from State, to State) bool {
// 	for _, s := range allowedStateChanges[from] {
// 		if s == to {
// 			return true
// 		}
// 	}
// 	return false
// }
