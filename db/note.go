package db

import "time"

type Note struct {
	ID        string    `json:"id"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"timestamp"`
}
