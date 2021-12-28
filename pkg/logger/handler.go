package logger

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ContextualHandler returns middleware adding a contextual logger with request IDs.
func ContextualHandler(opts ...ContextualOption) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		c := &contextualHandler{
			handler: h,
			logger:  zap.L(),
		}
		return parseContextualOptions(c, opts...)
	}
}

type contextualHandler struct {
	handler http.Handler
	logger  *zap.Logger
	headers []string
}

// ServeHTTP implements http.Handler.
func (h contextualHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	// Create request id and attach it to the context
	id := uuid.New().String()
	ctx = WithRequestID(ctx, id)

	// Create contextual logger with request id and attach it to the contet
	logger := h.logger.With(zap.String("id", id))
	ctx = ContextWithLogger(ctx, logger)

	// Log the incoming request with optional headers.
	var fields = []zap.Field{
		zap.String("method", r.Method),
		zap.String("path", r.RequestURI),
		zap.String("remote", r.RemoteAddr),
	}
	for _, header := range h.headers {
		fields = append(fields, zap.String(strings.ToLower(header), r.Header.Get(header)))
	}
	logger.Info("incoming request", fields...)

	// Call the wrapped handler attaching the contextual logger to the request.
	h.handler.ServeHTTP(w, r.WithContext(ctx))
}

// ContextualOption allows optional contextual logger configuration.
type ContextualOption func(http.Handler)

func parseContextualOptions(h http.Handler, opts ...ContextualOption) http.Handler {
	for _, option := range opts {
		option(h)
	}
	return h
}

// BaseLogger configures the logger from which to derive the contextual logger.
func WithBaseLogger(logger *zap.Logger) ContextualOption {
	return func(h http.Handler) {
		c := h.(*contextualHandler)
		c.logger = logger
	}
}

// WithHeaders configures which HTTP headers to log.
func WithHeaders(headers ...string) ContextualOption {
	return func(h http.Handler) {
		c := h.(*contextualHandler)
		c.headers = headers
	}
}
