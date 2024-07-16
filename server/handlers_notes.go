package server

import (
	"fmt"
	"net/http"

	"github.com/KatrinSalt/notes-service/api"
	"github.com/KatrinSalt/notes-service/notes"
	"github.com/google/uuid"
)

// http.Handler - interface with ServeHTTP method
// http.HandlerFunc - a function type that accepts same args as ServeHTTP method.
// It also implements the http.Handler interface.
func (s server) createNote() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		noteReq, err := decode[api.NoteRequest](r)
		if err != nil {
			fmt.Printf("Failed to decode the request: %s\n", err)
			statusCode, code := errorCodes(err)
			writeError(w, statusCode, code, err)
			return
		}
		if err := s.notes.CreateNote(toCreateNote(noteReq)); err != nil {
			if statusCode, code := errorCodes(err); statusCode != 0 {
				fmt.Printf("Failed to convert to a Note type: %s\n", err)
				writeError(w, statusCode, code, err)
				return
			}
			fmt.Printf("Failed to create a note: %s\n", err)
			writeServerError(w)
			return
		}

		if err := encode(w, http.StatusCreated, "Note is created"); err != nil {
			fmt.Printf("Failed to create a note: %s\n", err)
			writeServerError(w)
			return
		}
	})
}

func (s server) updateNote() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if len(id) == 0 {
			writeError(w, http.StatusBadRequest, CodeIDRequired, ErrIDRequired)
			return
		}

		noteReq, err := decode[api.NoteRequest](r)
		if err != nil {
			fmt.Printf("Failed to decode the request: %s\n", err)
			statusCode, code := errorCodes(err)
			writeError(w, statusCode, code, err)
			return
		}

		note, err := toUpdateNote(id, noteReq)
		if err != nil {
			fmt.Printf("Failed to convert to a Note type: %s\n", err)
			writeServerError(w)
			return
		}

		if err := s.notes.UpdateNote(note); err != nil {
			if statusCode, code := errorCodes(err); statusCode != 0 {
				fmt.Printf("Failed to convert to a Note type: %s\n", err)
				writeError(w, statusCode, code, err)
				return
			}
			fmt.Printf("Failed to update a note: %s\n", err)
			writeServerError(w)
			return
		}

		if err := encode(w, http.StatusOK, "Note is created"); err != nil {
			fmt.Printf("Failed to update a note: %s\n", err)
			writeServerError(w)
			return
		}
	})
}

func (s server) deleteNote() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if len(id) == 0 {
			writeError(w, http.StatusBadRequest, CodeIDRequired, ErrIDRequired)
			return
		}

		noteReq, err := decode[api.NoteRequest](r)
		if err != nil {
			fmt.Printf("Failed to decode the request: %s\n", err)
			statusCode, code := errorCodes(err)
			fmt.Printf("Status code: %d, Code: %s\n", statusCode, code)
			writeError(w, statusCode, code, err)
			return
		}

		if len(noteReq.Category) == 0 {
			writeError(w, http.StatusBadRequest, CodeCategoryRequired, ErrCategoryRequired)
			return
		}

		note, err := toDeleteNote(id, noteReq)
		if err != nil {
			fmt.Printf("Failed to convert to a Note type: %s\n", err)
			writeServerError(w)
			return
		}

		if err := s.notes.DeleteNote(note); err != nil {
			if statusCode, code := errorCodes(err); statusCode != 0 {
				fmt.Printf("Failed to delete a note: %s\n", err)
				writeError(w, statusCode, code, err)
				return
			}
			fmt.Printf("Failed to delete a note: %s\n", err)
			writeServerError(w)
			return
		}

		if err := encode(w, http.StatusOK, "Note is deleted"); err != nil {
			fmt.Printf("Failed to delete a note: %s\n", err)
			writeServerError(w)
			return
		}

	})
}

func toCreateNote(req api.NoteRequest) notes.Note {
	note := notes.Note{
		Category: req.Category,
		Note:     req.Note,
	}
	fmt.Printf("Note Type of notes.Note: %+v\n", note)
	return note
}

func toUpdateNote(id string, req api.NoteRequest) (notes.Note, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		fmt.Printf("Failed to parse the ID from the request: %s\n", err)
		return notes.Note{}, err

	}

	note := notes.Note{
		ID:       parsedID,
		Category: req.Category,
		Note:     req.Note,
	}
	fmt.Printf("Note Type of notes.Note: %+v\n", note)
	return note, nil
}

func toDeleteNote(id string, req api.NoteRequest) (notes.Note, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		fmt.Printf("Failed to parse the ID from the request: %s\n", err)
		return notes.Note{}, err

	}

	note := notes.Note{
		ID:       parsedID,
		Category: req.Category,
	}
	fmt.Printf("Note Type of notes.Note: %+v\n", note)
	return note, nil
}
