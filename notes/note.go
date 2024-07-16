package notes

import "github.com/google/uuid"

type Note struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Category string    `json:"category,omitempty"`
	Note     string    `json:"note,omitempty"`
}
