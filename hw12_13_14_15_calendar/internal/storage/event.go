package storage

import "time"

type Event struct {
	ID          string
	Title       string
	Description string
	UserID      int
	StartAt     time.Time
	FinishAt    time.Time
	NotifyAt    time.Time
}
