package model

import "time"

type Task struct {
  ID          int64
  Type        string
  Payload     []byte
  Status      string
  Attempts    int
  MaxAttempts int
  LockedUntil *time.Time
  CreatedAt   time.Time
  UpdatedAt   time.Time
}
