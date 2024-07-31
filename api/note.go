package api

type NoteRequest struct {
	Category string `json:"category,omitempty"`
	Note     string `json:"note,omitempty"`
}

type Note struct {
	ID       string `json:"id,omitempty"`
	Category string `json:"category,omitempty"`
	Note     string `json:"note,omitempty"`
}

type NoteResponse struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
