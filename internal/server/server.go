package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server represents the HTTP server
type Server struct {
	server       *http.Server
	logger       *Logger
	host         string
	port         int
	downloadsDir string
	manager      interface{} // Will be *taskmanager.Manager, avoiding circular import
}

// Config holds server configuration
type Config struct {
	Host         string
	Port         int
	DownloadsDir string
	Verbose      bool
}

// NewServer creates a new server instance
func NewServer(cfg Config) (*Server, error) {
	logger := NewLogger(cfg.Verbose)

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		// T083: Add connection timeouts and keep-alive settings
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	return &Server{
		server:       srv,
		logger:       logger,
		host:         cfg.Host,
		port:         cfg.Port,
		downloadsDir: cfg.DownloadsDir,
		manager:      nil, // Will be set after creation
	}, nil
}

// ListenAndServe starts the HTTP server
func (s *Server) ListenAndServe() error {
	s.logger.Info("Starting server on %s", s.server.Addr)

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for interrupt signal
	select {
	case <-sigChan:
		s.logger.Info("Shutdown signal received")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
			return err
		}
		s.logger.Info("Server stopped")
		return nil
	case err := <-errChan:
		return err
	}
}

// GetListener creates a network listener for the server
func (s *Server) GetListener() (net.Listener, error) {
	return net.Listen("tcp", s.server.Addr)
}

// Serve serves HTTP requests using the provided listener
func (s *Server) Serve(l net.Listener) error {
	s.logger.Info("Server listening on %s", l.Addr())
	return s.server.Serve(l)
}

// RegisterHandler registers a handler pattern
func (s *Server) RegisterHandler(pattern string, handler http.HandlerFunc) {
	s.logger.Debug("Registering handler: %s", pattern)
	// Access the underlying ServeMux via Handler
	if mux, ok := s.server.Handler.(*http.ServeMux); ok {
		mux.HandleFunc(pattern, handler)
	}
}

// SetTaskManager sets the task manager (called during initialization)
func (s *Server) SetTaskManager(manager interface{}) {
	s.manager = manager
}

// GetTaskManager returns the task manager
func (s *Server) GetTaskManager() interface{} {
	return s.manager
}
