package server

import (
	"fmt"
	"net/http"

	"github.com/KatrinSalt/notes-service/api"
	"github.com/KatrinSalt/notes-service/notes"
)

func (s server) createNote() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		note, err := decode[api.NoteCreateRequest](r)
		if err != nil {
			fmt.Printf("Failed to decode the request: %s\n", err)
			statusCode, code := errorCodes(err)
			writeError(w, statusCode, code, err)
			return
		}

		if err := s.notes.CreateNote(toCreateNote(note)); err != nil {
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

func toCreateNote(req api.NoteCreateRequest) notes.Note {
	note := notes.Note{
		ID:   req.ID,
		Note: req.Note,
	}
	fmt.Printf("Note Type of notes.Note: %+v\n", note)
	return note
	// return notes.Note{
	// 	Note: req.Note,
	// }
}
