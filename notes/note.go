package notes

type Note struct {
	ID       string `json:"id,omitempty"`
	Category string `json:"category,omitempty"`
	Note     string `json:"note,omitempty"`
}
