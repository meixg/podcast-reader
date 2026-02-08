package server

import (
	"log"
	"os"
)

// Logger provides structured logging for the server
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	verbose     bool
}

// NewLogger creates a new logger instance
func NewLogger(verbose bool) *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime),
		debugLogger: log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile),
		verbose:     verbose,
	}
}

// Info logs informational messages
func (l *Logger) Info(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

// Error logs error messages
func (l *Logger) Error(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

// Debug logs debug messages (only when verbose is enabled)
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.verbose {
		l.debugLogger.Printf(format, v...)
	}
}
