package api

type NoteCreateRequest struct {
	Category string `json:"category"`
	Note     string `json:"note"`
}
