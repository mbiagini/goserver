package gsmiddleware

import (
	"fmt"
	"goserver/utils/gslog"
	"goserver/utils/gsrender"
	"net/http"
	"runtime/debug"
)

// Recoverer is an http.Handler that recovers from panics, logs the panic and returns
// an http error to the client.

// Receives the default error to be written to the http.ResponseWriter.
// Uses mblog as logger (not an implementation of log.Logger).
func Recoverer(v interface{}, ts gsrender.TextStandard) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					gslog.Error(fmt.Sprintf("Handling panic: %s", string(debug.Stack())), GetTraceID(r.Context()))
					gsrender.Write(w, http.StatusInternalServerError, ts, v)
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}