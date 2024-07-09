package notes

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KatrinSalt/notes-service/db"
)

const (
	// defaultServiceTimeout is the default timeout for service operations.
	defaultServiceTimeout = 15 * time.Second
)

type Database interface {
	// CreateNote creates a new note.
	CreateNote(ctx context.Context, note db.Note) error
}

type Service interface {
	// CreateNote creates a new note.
	CreateNote(note Note) error
	// GetNoteByID returns a note by its ID.
	// GetNoteByID(id string) (string, error)
	// // UpdateNoteByID updates a note by its ID.
	// UpdateNoteByID(id string, note Note) error
	// // DeleteNoteByID deletes a note by its ID.
	// DeleteNoteByID(id string) error
}

type service struct {
	db      Database
	timeout time.Duration
}

// ServiceOptions contains options for the service.
type ServiceOptions struct {
	Timeout time.Duration
}

// ServiceOption is a function that sets options on the service.
type ServiceOption func(o *ServiceOptions)

func NewService(db Database, options ...ServiceOption) (*service, error) {
	if db == nil {
		return nil, errors.New("database must not be nil")
	}

	opts := ServiceOptions{
		Timeout: defaultServiceTimeout,
	}
	for _, option := range options {
		option(&opts)
	}

	return &service{
		db:      db,
		timeout: opts.Timeout,
	}, nil
}

func (s service) CreateNote(note Note) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	if err := s.db.CreateNote(ctx, toNoteDB(note)); err != nil {
		fmt.Printf("Failed to create a note in DB: %s\n", err)
		return err
	}

	return nil
}

func toNoteDB(note Note) db.Note {
	noteDB := db.Note{
		ID:        note.ID,
		Note:      note.Note,
		CreatedAt: time.Now().UTC(),
	}
	fmt.Printf("Note Type of db.Note: %+v\n", noteDB)
	return noteDB
	// return db.Note{
	// 	Note:      note.Note,
	// 	CreatedAt: time.Now().UTC(),
	// }
}
