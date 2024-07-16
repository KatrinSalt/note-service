package server

func (s server) routes() {
	// Create a note handler
	s.router.Handle("POST /createNote", s.createNote())
	s.router.Handle("PUT /updateNote/{id}", s.updateNote())
	s.router.Handle("DELETE /deleteNote/{id}", s.deleteNote())

}
