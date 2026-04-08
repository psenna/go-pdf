package config

import (
	"time"
)

// Config holds the application configuration
type Config struct {
	// ConcurrencyLimit is the maximum number of concurrent PDF processing requests
	ConcurrencyLimit int
	// Timeout is the maximum time allowed for PDF processing
	Timeout time.Duration
	// MaxFileSize is the maximum allowed file size for PDF uploads
	MaxFileSize int64
}

// GlobalConfig is the singleton instance of Config
var GlobalConfig = &Config{
	ConcurrencyLimit: 10,
	Timeout:          30 * time.Second,
	MaxFileSize:      100 * 1024 * 1024, // 100MB default
}

// Init initializes the configuration with provided values
func Init(cfg Config) {
	GlobalConfig = &cfg
}
