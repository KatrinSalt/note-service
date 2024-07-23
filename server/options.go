package server

// WithAddress sets the address for the server.
func WithAddress(address string) Option {
	return func(s *server) {
		s.httpServer.Addr = address
	}
}

// WithLogger sets the logger for the server.
func WithLogger(log logger) Option {
	return func(s *server) {
		s.log = log
	}
}
