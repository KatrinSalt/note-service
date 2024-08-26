package db

import (
	"time"
)

type Note struct {
	ID        string    `json:"id"`
	Category  string    `json:"category"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"timestamp"`
}
