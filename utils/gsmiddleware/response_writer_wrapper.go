package gsmiddleware

import (
	"bytes"
	"net/http"
)

// ResponseWriterWrapper is used to wrap an http.ResponseWriter implementation
// an be able to read from it without modifying the actual writer.
type ResponseWriterWrapper struct {
	w          *http.ResponseWriter
	body       *bytes.Buffer
	statusCode *int
	copyBody   bool
}

func NewResponseWriterWrapper(w http.ResponseWriter, copyBody bool) *ResponseWriterWrapper {
	var buf bytes.Buffer
	var statusCode int = 200
	return &ResponseWriterWrapper{
		w: 			&w,
		body: 		&buf,
		statusCode: &statusCode,
		copyBody:   copyBody,
	}
}

func (rww *ResponseWriterWrapper) Write(buf []byte) (int, error) {
	if rww.copyBody {
		rww.body.Write(buf)
	}
	return (*rww.w).Write(buf)
}

func (rww *ResponseWriterWrapper) Header() http.Header {
	return (*rww.w).Header()
}

func (rww *ResponseWriterWrapper) WriteHeader(statusCode int) {
	(*rww.statusCode) = statusCode
	(*rww.w).WriteHeader(statusCode)
}