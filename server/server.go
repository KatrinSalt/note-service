package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/KatrinSalt/notes-service/notes"
)

// Defaults for server configuration.
const (
	defaultHost         = "localhost"
	defaultPort         = "3000"
	defaultReadTimeout  = 15 * time.Second
	defaultWriteTimeout = 15 * time.Second
	defaultIdleTimeout  = 30 * time.Second
)

// server holds an http.Server, a router and it's configured options.
type server struct {
	httpServer *http.Server
	router     *http.ServeMux
	// TODO: to be added later
	notes   notes.Service
	stopCh  chan os.Signal
	errCh   chan error
	started bool
}

// Options holds the configuration for the server.
type Options struct {
	Router       *http.ServeMux
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Option is a function that configures the server.
type Option func(*server)

// New returns a new server.
func New(notes notes.Service, options ...Option) (*server, error) {
	if notes == nil {
		return nil, errors.New("notes service is nil")
	}

	s := &server{
		httpServer: &http.Server{
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
			IdleTimeout:  defaultIdleTimeout,
		},
		notes:  notes,
		stopCh: make(chan os.Signal),
		errCh:  make(chan error),
	}

	for _, option := range options {
		option(s)
	}

	if s.router == nil {
		s.router = http.NewServeMux()
		s.httpServer.Handler = s.router
	}

	if len(s.httpServer.Addr) == 0 {
		s.httpServer.Addr = defaultHost + ":" + defaultPort
	}

	return s, nil
}

// Start the server.
func (s *server) Start() error {
	defer func() {
		s.started = false
	}()

	s.routes()

	// Question: in which order these two functions would be called?
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.errCh <- err
		}
	}()

	go func() {
		s.stop()
	}()

	s.started = true
	fmt.Printf("Server started at %s\n", s.httpServer.Addr)
	for {
		select {
		case err := <-s.errCh:
			close(s.errCh)
			return err
		case sig := <-s.stopCh:
			fmt.Printf("Server stopped. Reason: %s\n", sig.String())
			close(s.stopCh)
			return nil
		}
	}
}

// stop the server.
func (s server) stop() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	s.httpServer.SetKeepAlivesEnabled(false)
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.errCh <- err
	}

	s.stopCh <- sig
}

// WithOptions configures the server with the given Options.
func WithOptions(options Options) Option {
	return func(s *server) {
		if options.Router != nil {
			s.router = options.Router
			s.httpServer.Handler = s.router
		}
		// NOTE: question to logic: what if options.Host is empty and options.Port is not?
		if len(options.Host) > 0 || options.Port > 0 {
			s.httpServer.Addr = options.Host + ":" + strconv.Itoa(options.Port)
		}
		if options.ReadTimeout > 0 {
			s.httpServer.ReadTimeout = options.ReadTimeout
		}
		if options.WriteTimeout > 0 {
			s.httpServer.WriteTimeout = options.WriteTimeout
		}
		if options.IdleTimeout > 0 {
			s.httpServer.IdleTimeout = options.IdleTimeout
		}
	}
}
