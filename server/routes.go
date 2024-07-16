package server

func (s server) routes() {
	// Create a note handler
	s.router.Handle("POST /createNote", s.createNote())
	s.router.Handle("POST /updateNote/{id}", s.updateNote())

}
