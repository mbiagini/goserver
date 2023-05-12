package gsmiddleware

import (
	"goserver/utils/gslog"
	"net/http"
)

// HttpLogTransport wraps an http.Transport adding the logging of request and
// response. For this, uses mblog (not an implementation of log.Logger).
type HttpLogTransport struct {
	Base http.Transport
}

// Implements interface http.RoundTripper so it can be used as a Transport.
func (lt *HttpLogTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	
	gslog.Request(
		gslog.INNER_REQUEST, 
		r, 
		GetTraceID(r.Context()),
	)

	// Send the request, get the response (or the error).
	rs, err := lt.Base.RoundTrip(r)

	// Handle the result.
	if (err != nil) {
		return nil, err
	}

	gslog.Response(
		rs,
		GetTraceID(r.Context()),
	)

	return rs, nil 
}