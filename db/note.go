package db

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID        uuid.UUID `json:"id"`
	Category  string    `json:"category"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"timestamp"`
}
