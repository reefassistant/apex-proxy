package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type loggerKey int

const (
	ctxLogger loggerKey = iota
	ctxRequestID
)

// ContextFields returns a new context appending fields to the contextual logger.
func ContextWithFields(ctx context.Context, fields ...zap.Field) context.Context {
	return context.WithValue(ctx, ctxLogger, Context(ctx).With(fields...))
}

// ContextLogger returns a new context with given logger attached.
func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger, logger)
}

// Context returns the logger from given context.
func Context(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(ctxLogger).(*zap.Logger); ok {
		return logger
	}
	return zap.L() // Fallback to global logger if none is attached
}

// WithRequestID returns a context with given request ID.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxRequestID, id)
}

// RequestIDFrom returns the request ID from context.
func RequestIDFrom(ctx context.Context) (string, error) {
	if id, ok := ctx.Value(ctxRequestID).(string); ok {
		return id, nil
	}
	return "", fmt.Errorf("request id not set")
}
