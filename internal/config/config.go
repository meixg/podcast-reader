package config

import (
	"fmt"
	"time"
)

// Config holds the application configuration.
type Config struct {
	// OutputDirectory is where downloaded files are saved
	OutputDirectory string

	// OverwriteExisting controls whether to overwrite existing files
	OverwriteExisting bool

	// Timeout is the HTTP request timeout
	Timeout time.Duration

	// MaxRetries is the maximum number of retry attempts for failed downloads
	MaxRetries int

	// RetryDelay is the base delay between retries (exponential backoff)
	RetryDelay time.Duration

	// ShowProgress controls whether to display download progress
	ShowProgress bool

	// ValidateFiles controls whether to validate downloaded audio files
	ValidateFiles bool

	// ServerHost is the host address for the HTTP server
	ServerHost string

	// ServerPort is the port number for the HTTP server
	ServerPort int

	// Verbose enables verbose logging
	Verbose bool
}

// DefaultConfig returns a configuration with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		OutputDirectory:   "./downloads",
		OverwriteExisting: false,
		Timeout:           30 * time.Second,
		MaxRetries:        3,
		RetryDelay:        1 * time.Second,
		ShowProgress:      true,
		ValidateFiles:     true,
		ServerHost:        "localhost",
		ServerPort:        8080,
		Verbose:           false,
	}
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	// Check OutputDirectory is valid path (simplified check)
	if c.OutputDirectory == "" {
		return fmt.Errorf("output directory cannot be empty")
	}

	// Check Timeout is positive
	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	// Check MaxRetries is non-negative
	if c.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative")
	}

	return nil
}
