package api

type NoteRequest struct {
	Category string `json:"category"`
	Note     string `json:"note"`
}
