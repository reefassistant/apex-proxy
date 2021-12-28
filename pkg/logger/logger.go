package logger

import (
	"go.uber.org/zap"
)

const (
	LEVEL_ERROR   = "error"
	LEVEL_WARNING = "warning"
	LEVEL_INFO    = "info"
	LEVEL_DEBUG   = "debug"
)

// Initialize the global logger with given optional configuration parameters.
// This should be called only once during application bootstrapping.
func Initialize(opts ...LoggerOption) error {

	var (
		logger *zap.Logger
		cfg    *config
		err    error
	)

	// Create default configuration
	cfg = &config{
		encoding:          "json",
		logLevel:          LEVEL_INFO,
		logFileMaxSize:    10,
		logFileMaxBackups: 30,
		logFileMaxAge:     0,
		logFileCompress:   true,
		outputPaths:       []string{"stdout"},
	}

	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}

	// Determine log level
	cfg.logLevelCore, err = logLevel(cfg.logLevel)
	if err != nil {
		return err
	}

	// Handle file logging
	registerFileLogger(cfg)

	// Create the new logger and replace the global one
	logger, err = transformLoggerConfig(cfg).Build()
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(logger)
	return nil
}
