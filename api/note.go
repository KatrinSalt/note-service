package api

type NoteCreateRequest struct {
	ID   string `json:"id"`
	Note string `json:"note"`
}
