package server

// WithAddress sets the address for the server.
func WithAddress(address string) Option {
	return func(s *server) {
		s.httpServer.Addr = address
	}
}
