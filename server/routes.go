package server

func (s server) routes() {
	s.router.Handle("POST /notes/create/{category}", s.createNote())
	s.router.Handle("PUT /notes/update/{category}/{id}", s.updateNote())
	s.router.Handle("DELETE /notes/delete/{category}/{id}", s.deleteNote())
	s.router.Handle("GET /notes/categories/{category}/ids/{id}", s.getNoteByID())
	s.router.Handle("GET /notes/categories/{category}", s.getNotesByCategory())
}
