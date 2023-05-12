package gsmiddleware

import (
	"goserver/utils/gslog"
	"net/http"
)

// This handler logs the request as the server receives it and the response
// as it is returned by the server.
// Uses a wrapper for ResponseWriter to store the response body.
func HttpLogHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// Log Request.
		gslog.Request(
			gslog.OUTER_REQUEST,
			r,
			GetTraceID(r.Context()),
		)

		ww := NewResponseWriterWrapper(w, true)

		// Defer Log Response.
		defer func() {
			gslog.ResponseWriter(
				ww.body,
				ww.Header(),
				*ww.statusCode,
				r,
				GetTraceID(r.Context()),
			)
		}()

		// Pass the wrapped ResponseWriter to next handler.
		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(fn)
}