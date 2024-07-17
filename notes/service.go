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
	CreateNote(ctx context.Context, note *db.Note) error
	// UpdateNote updates a note.
	UpdateNote(ctx context.Context, note *db.Note) error
	// DeleteNote deletes a note.
	DeleteNote(ctx context.Context, id, category string) error
	// GetNotesByCategory returns a list of notes stored in DB.
	GetNotesByCategory(ctx context.Context, category string) ([]db.Note, error)
	// GetNoteByID returns a notes with id <id>.
	GetNoteByID(ctx context.Context, category, id string) (db.Note, error)
}

type Service interface {
	// CreateNote creates a new note.
	CreateNote(note Note) error
	// GetNoteByID returns a note by its ID.
	// GetNoteByID(id string) (string, error)
	// UpdateNote updates a note.
	UpdateNote(note Note) error
	// DeleteNote deletes a note by its ID.
	DeleteNote(note Note) error
	// GetNotesByCategory returns a list of notes stored in DB.
	GetNotesByCategory(category string) ([]Note, error)
	// GetNoteByID returns a notes with id <id>.
	GetNoteByID(category, id string) (Note, error)
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

func (s service) UpdateNote(note Note) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	if err := s.db.UpdateNote(ctx, toNoteDB(note)); err != nil {
		fmt.Printf("Failed to update a note in DB: %s\n", err)
		return err
	}

	return nil
}

func (s service) DeleteNote(note Note) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	if err := s.db.DeleteNote(ctx, note.ID.String(), note.Category); err != nil {
		fmt.Printf("Failed to delete a note in DB: %s\n", err)
		return err
	}

	return nil
}

func (s service) GetNotesByCategory(category string) ([]Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	notesDB, err := s.db.GetNotesByCategory(ctx, category)
	if err != nil {
		fmt.Printf("Failed to get notes from DB: %s\n", err)
		return nil, err
	}

	notes := make([]Note, len(notesDB))
	for i := range notesDB {
		notes[i] = fromNoteDB(&notesDB[i])
	}

	return notes, nil
}

func (s service) GetNoteByID(category, id string) (Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	noteDB, err := s.db.GetNoteByID(ctx, category, id)
	if err != nil {
		fmt.Printf("Failed to get a note from DB: %s\n", err)
		return Note{}, err
	}

	return fromNoteDB(&noteDB), nil
}

func toNoteDB(note Note) *db.Note {
	noteDB := &db.Note{
		ID:        note.ID,
		Category:  note.Category,
		Note:      note.Note,
		CreatedAt: time.Now().UTC(),
	}
	fmt.Printf("Note Type of db.Note: %+v\n", noteDB)
	return noteDB
}

func fromNoteDB(noteDB *db.Note) Note {
	note := Note{
		ID:       noteDB.ID,
		Category: noteDB.Category,
		Note:     noteDB.Note,
	}
	fmt.Printf("Note Type of notes.Note: %+v\n", note)
	return note
}
