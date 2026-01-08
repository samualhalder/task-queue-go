package model

import "time"

type Task struct {
  ID          int64
  Type        string
  Payload     []byte
  Status      State
  Attempts    int
  MaxAttempts int
  LockedUntil *time.Time
  CreatedAt   time.Time
  UpdatedAt   time.Time
}

type State string 

const (
  StatePending State="pending"
  StateRunning State="running"
  StateFailed State="failed"
  StateCompleted State="completed"
)

var allowedStateChanges=map[State][]State{
  StatePending:{StateRunning},
  StateRunning:{StateCompleted,StateFailed},
  StateFailed:{StatePending},
  StateCompleted:{},
}

func canStateChange(from State,to State) bool{
  for _,s := range allowedStateChanges[from]{
      if s==to{
        return true
      }
  }
  return false
}