// Package httpserver implements HTTP server.
package httpserver

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	bytesConvConst = 1024

	_defaultAddr            = ":80"
	_defaultBodyLimitMb     = 20
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

// Server -.
type Server struct {
	App    *fiber.App
	notify chan error

	address         string
	prefork         bool
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration
	bodyLimit       int
}

// New -.
func New(opts ...Option) *Server {
	s := &Server{
		App:             nil,
		notify:          make(chan error, 1),
		address:         _defaultAddr,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
		bodyLimit:       _defaultBodyLimitMb * bytesConvConst * bytesConvConst,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	app := fiber.New(fiber.Config{
		Prefork:      s.prefork,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
		BodyLimit:    s.bodyLimit,
	})

	s.App = app

	return s
}

// Start -.
func (s *Server) Start() {
	go func() {
		s.notify <- s.App.Listen(s.address)

		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	return s.App.ShutdownWithTimeout(s.shutdownTimeout)
}
