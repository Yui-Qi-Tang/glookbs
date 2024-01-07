// Package httpserver is a helper to start a http server with graceful shutdown
package httpserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

// TLSFile is an interface for external pkg to setup tls file path if run http with tls
type TLSFile interface {
	Key() string
	Cert() string
}

// interSignal implements os.Signal
type interSignal string

func (s interSignal) Signal() {}

func (s interSignal) String() string {
	return string(s)
}

var sigQuit interSignal = "signal quit"

// Option is an option form to make configuration with server
type Option func(*http.Server)

// Server is http server
type Server struct {
	s   *http.Server
	err error
}

// WithAddr sets address for http server
func WithAddr(addr string) Option {
	return func(s *http.Server) {
		s.Addr = addr
	}
}

// WithHandler set http handler
func WithHandler(handler http.Handler) Option {
	return func(s *http.Server) {
		s.Handler = handler
	}
}

// New returns http(s) server
func New(opts ...Option) *Server {
	srv := &http.Server{
		Addr: ":8080",
	}

	for _, opt := range opts {
		opt(srv)
	}

	return &Server{
		s: srv,
	}
}

// Run starts httpserver
func (s *Server) Run(tls TLSFile) error {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, sigQuit)

	if tls != nil {
		go func(key, cert string) {
			if err := s.s.ListenAndServeTLS(cert, key); err != nil {
				s.err = err
				quit <- sigQuit
			}
		}(tls.Key(), tls.Cert())
	} else {
		go func() {
			if err := s.s.ListenAndServe(); err != nil {
				s.err = err
				quit <- sigQuit
			}
		}()
	}

	log.Println("server is running at:", s.s.Addr)

	v := <-quit
	if v == sigQuit {
		// server run with initial dead
		return s.err
	}

	log.Println("received shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.s.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown error")
	}

	log.Println("server is exiting...")

	return nil
}
