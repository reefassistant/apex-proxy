package logger

import (
	"go.uber.org/zap/zapcore"
)

type config struct {
	development       bool
	encoding          string
	logLevel          string
	logLevelCore      zapcore.Level
	logFile           string
	logFileMaxSize    int
	logFileMaxBackups int
	logFileMaxAge     int
	logFileCompress   bool
	outputPaths       []string
}

// LoggerOption allows configuring the logger.
type LoggerOption func(*config)

// WithLogLevel configures the log level.
func WithLogLevel(level string) LoggerOption {
	return func(c *config) {
		c.logLevel = level
	}
}

// WithLogFile configures the the file to store log output.
func WithLogFile(file string) LoggerOption {
	return func(c *config) {
		c.logFile = file
	}
}

// WithLogFileMaxSize configures the maximum log file size in MB.
func WithLogFileMaxSize(size int) LoggerOption {
	return func(c *config) {
		c.logFileMaxSize = size
	}
}

// WithLogFileMaxBackups configures the maximum amount of backup files to maintain.
func WithLogFileMaxBackups(backups int) LoggerOption {
	return func(c *config) {
		c.logFileMaxBackups = backups
	}
}

// WithLogFileMaxAge configures the age in days before a log file is removed.
func WithLogFileMaxAge(age int) LoggerOption {
	return func(c *config) {
		c.logFileMaxAge = age
	}
}

// WithLogFileCompress enables compression on the rotated backup files.
func WithLogFileCompress() LoggerOption {
	return func(c *config) {
		c.logFileCompress = true
	}
}

// WithDevelopment enables more liberal stack traces.
func WithDevelopment() LoggerOption {
	return func(c *config) {
		c.development = true
		c.encoding = "console"
	}
}
