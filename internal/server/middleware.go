package server

import (
	"log"
	"net/http"
	"runtime/debug"
)

// LoggingMiddleware logs all HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware recovers from panics and logs errors
func RecoveryMiddleware(logger *Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered: %v\nStack: %s", err, debug.Stack())
				WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", "")
			}
		}()
		http.DefaultServeMux.ServeHTTP(w, r)
	})
}

// ChainMiddleware chains multiple middleware functions
func ChainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
