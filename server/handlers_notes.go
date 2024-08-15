package server

import (
	"net/http"

	"github.com/KatrinSalt/notes-service/api"
	"github.com/KatrinSalt/notes-service/notes"
)

// http.Handler - interface with ServeHTTP method
// http.HandlerFunc - a function type that accepts same args as ServeHTTP method.
// It also implements the http.Handler interface.
func (s server) createNote() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		noteReq, err := decode[api.NoteRequest](r)
		if err != nil {
			statusCode, code := errorCodes(err)
			writeError(w, statusCode, code, err)
			return
		}
		data, err := s.notes.CreateNote(toCreateNote(noteReq))
		if err != nil {
			s.log.Error("Failed to create a note.", logError(err, "createNote")...)
			if statusCode, code := errorCodes(err); statusCode != 0 {
				writeError(w, statusCode, code, err)
				return
			}
			writeServerError(w)
			return
		}

		response := api.NoteResponse{
			Message: "Note is created",
			Data:    toNoteAPI(data),
		}

		if err := encode(w, http.StatusCreated, response); err != nil {
			s.log.Error("Failed to create a note.", logError(err, "createNote")...)
			writeServerError(w)
			return
		}
		s.log.Info("Note is created.", "type", "service", "name", "noteService", "method", "Create", "noteCategory", data.Category, "noteID", data.ID)
	})
}

// TODO: can't update the category of the note!
func (s server) updateNote() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: add proper validation
		// it is assumed that the id is provided in the path
		id := r.PathValue("id")

		noteReq, err := decode[api.NoteRequest](r)
		if err != nil {
			statusCode, code := errorCodes(err)
			writeError(w, statusCode, code, err)
			return
		}

		note := toUpdateNote(id, noteReq)

		data, err := s.notes.UpdateNote(note)
		if err != nil {
			s.log.Error("Failed to update the note.", logError(err, "updateNote")...)
			if statusCode, code := errorCodes(err); statusCode != 0 {
				writeError(w, statusCode, code, err)
				return
			}
			writeServerError(w)
			return
		}

		response := api.NoteResponse{
			Message: "Note is updated",
			Data:    toNoteAPI(data),
		}

		if err := encode(w, http.StatusOK, response); err != nil {
			s.log.Error("Failed to update a note.", logError(err, "updateNote")...)
			writeServerError(w)
			return
		}
		s.log.Info("Note is updated.", "type", "service", "name", "noteService", "method", "Update", "noteCategory", data.Category, "noteID", data.ID)
	})
}

func (s server) deleteNote() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: add proper validation
		// it is assumed that the id is provided in the path
		id := r.PathValue("id")

		noteReq, err := decode[api.NoteRequest](r)
		if err != nil {
			statusCode, code := errorCodes(err)
			writeError(w, statusCode, code, err)
			return
		}

		// Note: leaving it temporary as I am gonna change the logic later (remove category as the partition key)
		if len(noteReq.Category) == 0 {
			writeError(w, http.StatusBadRequest, CodeCategoryRequired, ErrCategoryRequired)
			return
		}

		note := toDeleteNote(id, noteReq)

		if err := s.notes.DeleteNote(note); err != nil {
			s.log.Error("Failed to delete a note.", logError(err, "deleteNote")...)
			if statusCode, code := errorCodes(err); statusCode != 0 {
				writeError(w, statusCode, code, err)
				return
			}
			writeServerError(w)
			return
		}

		if err := encode(w, http.StatusOK, "Note is deleted"); err != nil {
			s.log.Error("Failed to delete a note.", logError(err, "deleteNote")...)
			writeServerError(w)
			return
		}

		s.log.Info("Note is deleted.", "type", "service", "name", "noteService", "method", "Delete", "noteCategory", note.Category, "noteID", note.ID)
	})
}

func (s server) getNotesByCategory() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: add proper validation
		// it is assumed that the category is provided in the path
		category := r.PathValue("category")

		notes, err := s.notes.GetNotesByCategory(category)
		if err != nil {
			s.log.Error("Failed to list notes in the category.", logError(err, "getNotesByCategory")...)
			if statusCode, code := errorCodes(err); statusCode != 0 {
				writeError(w, statusCode, code, err)
				return
			}
			writeServerError(w)
			return
		}

		if err := encode(w, http.StatusOK, toNotesAPI(notes)); err != nil {
			s.log.Error("Failed to list notes in the category.", logError(err, "getNotesByCategory")...)
			writeServerError(w)
			return
		}
		s.log.Info("Notes are listed.", "type", "service", "name", "noteService", "method", "getNotesByCategory", "notesCategory", category)
	})
}

func (s server) getNoteByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: add proper validation
		// it is assumed that the category and id are provided in the path
		category := r.PathValue("category")
		id := r.PathValue("id")

		note, err := s.notes.GetNoteByID(category, id)
		if err != nil {
			s.log.Error("Failed to get a note.", logError(err, "getNoteByID")...)
			if statusCode, code := errorCodes(err); statusCode != 0 {
				writeError(w, statusCode, code, err)
				return
			}
			writeServerError(w)
			return
		}

		if err := encode(w, http.StatusOK, toNoteAPI(note)); err != nil {
			s.log.Error("Failed to get a note.", logError(err, "getNoteByID")...)
			writeServerError(w)
			return
		}
		s.log.Info("Note is found.", "type", "service", "name", "noteService", "method", "getNoteByID", "noteID", id)
	})
}

func toCreateNote(req api.NoteRequest) notes.Note {
	note := notes.Note{
		Category: req.Category,
		Note:     req.Note,
	}
	return note
}

func toUpdateNote(id string, req api.NoteRequest) notes.Note {
	note := notes.Note{
		ID:       id,
		Category: req.Category,
		Note:     req.Note,
	}
	return note
}

func toDeleteNote(id string, req api.NoteRequest) notes.Note {
	note := notes.Note{
		ID:       id,
		Category: req.Category,
	}
	return note
}

func toNoteAPI(note notes.Note) api.Note {
	return api.Note{
		ID:       note.ID,
		Category: note.Category,
		Note:     note.Note,
	}
}

func toNotesAPI(notes []notes.Note) []api.Note {
	notesAPI := make([]api.Note, len(notes))
	for i := range notes {
		notesAPI[i] = toNoteAPI(notes[i])
	}
	return notesAPI
}
