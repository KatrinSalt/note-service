package server

func (s server) routes() {
	// Create a note handler
	s.router.Handle("POST /create", s.createNote())
	s.router.Handle("PUT /update/{id}", s.updateNote())
	s.router.Handle("DELETE /delete/{id}", s.deleteNote())
	s.router.Handle("GET /notes/category/{category}/id/{id}", s.getNoteByID())
	s.router.Handle("GET /notes/category/{category}", s.getNotesByCategory())
}
