package payloads

import "time"

type TaskQueryPayload struct {
	TaskName      *string
	Status        *string
	CreatedBefore *time.Time
	CreatedAfter  *time.Time
	OrderBy       string
	SortOrder     string
	Limit         int
	Skip          int
}
