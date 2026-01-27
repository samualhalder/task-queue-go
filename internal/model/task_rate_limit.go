package model

import "time"

type TaskRateLimit struct {
	Key          string
	Capacity     int
	RefillRate   float64
	Tokens       float64
	LastRefillAt time.Time
}
