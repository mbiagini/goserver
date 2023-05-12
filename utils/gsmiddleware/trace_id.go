package gsmiddleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// Key type to use when setting the trace ID.
type ctxKeyTraceId int

// TraceIDKey is the key that holds the unique trace ID in a request context.
const TraceIDKey ctxKeyTraceId = 0

// TraceIDHeader is the name of the HTTP header which contains the trace ID.
// Exported so that it can be changed if needed.
var TraceIDHeader = "x-trace-id"

// TracingTransport wraps an http.RoundTripper adding the trace-id if found in
// the request context. For this to work, the request should be generated using
// a context that already has the trace-id set.
type TracingTransport struct {
	Base http.RoundTripper
}

// Implements interface http.RoundTripper so it can be used as a Transport.
// For this transport to work with http requests on external APIs, the
// http.Request's Context must contain a trace id. This can be achieved by
// creating it with http.NewRequestWithContext function and passing a context
// that already contains the id.
func (tt *TracingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	traceID := GetTraceID(r.Context())
	if traceID != "" {
		r.Header.Set(TraceIDHeader, traceID)
	}
	return tt.Base.RoundTrip(r)
}

// TraceID is a middleware that injects a trace ID into the context of each
// request. UUID is used as trace ID.
// Once the ID is set in the context, it can be used in any step of the flow
// by anyone who receives the context.
func TraceID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// ID to be used across request flow.
		traceID := new()

		// Add the id as request and response headers for better trace.
		r.Header.Set(TraceIDHeader, traceID)
		w.Header().Add(TraceIDHeader, traceID)

		// Add traceId to request context.
		ctx := context.WithValue(r.Context(), TraceIDKey, traceID)

		next.ServeHTTP(w, r.Clone(ctx))
	}
	return http.HandlerFunc(fn)
}

// GetTraceID returns a request trace ID from the given context if one is present.
// Returns the empty string if no traceID can be found.
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// new returns a new unique trace-id using a random uuid string.
func new() string {
	return uuid.New().String()
}