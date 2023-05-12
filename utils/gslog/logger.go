package gslog

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// minLogLevel specifies the log level from which the lines will be written.
// The order is DEBUG < INFO < WARN < ERROR.
var minLogLevel logLevel = INFO

// maxBodyLength specifies the number of characters of the request/response 
// body to be logged.
var maxBodyLength int = 1000

// excludeUrls is a set (sort of) of urls that, if found in a request, its
// request and response won't be logged.
var excludeUrls map[string] struct{}

// LogConfig contains the configuration properties used in the logging. If this
// function is not called, default values will be used.
type LoggerConfig struct {

	// Level allows the configuration of the internal var minLogLevel.
	Level         *logLevel `json:"level"`

	// MaxBodyLength allows the configuration of hte internal var maxBodyLength.
	MaxBodyLength *int      `json:"max_body_length"`

	// ExcludeUrls specifies the urls to be added to the internal set excludeUrls.
	ExcludeUrls   []string  `json:"exclude_urls"`
}

// Initialization. Only initializes the excludeUrls map.
func init() {
	excludeUrls = make(map[string]struct{})
}

// ConfigureLog enables the configuration of the logger.
func ConfigureLog(c LoggerConfig) {
	if c.Level != nil {
		minLogLevel = *c.Level
	}
	if c.MaxBodyLength != nil {
		maxBodyLength = *c.MaxBodyLength
	}
	if c.ExcludeUrls != nil {
		for _, url := range c.ExcludeUrls {
			excludeUrls[url] = struct{}{}
		}
	}
}

// logMessage writes a new log line. Receives the level, message and traceID.
func logMessage(level string, m string, traceID string) {
	l := &customLog{
		Time: timeString(),
		TraceID: traceID,
		Level: level,
		Type: MESSAGE,
		Message: m,
	}
	l.print()
}

// Debug creates a Log with the given message and Level DEBUG and prints it.
func Debug(m string, traceID string) {
	if minLogLevel <= DEBUG {
		logMessage("DEBUG", m, traceID)
	}
}

// Info creates a Log with the given message and Level INFO and prints it.
func Info(m string, traceID string) {
	if minLogLevel <= INFO {
		logMessage("INFO", m, traceID)
	}
}

// Warn creates a Log with the given message and Level WARN and prints it.
func Warn(m string, traceID string) {
	if minLogLevel <= WARN {
		logMessage("WARN", m, traceID)
	}
}

// Error creates a Log with the given message and Level ERROR and prints it.
func Error(m string, traceID string) {
	logMessage("ERROR", m, traceID)
}

// ErrorFrom creates a Log from the given Error and prints it.
func ErrorFrom(err error, traceID string) {
	Error(err.Error(), traceID)
}

// Server writes the message with the following format: [server - %time] %s.
// This function is only meant to be called when starting up a server. Any
// program watching log files should ignore the lines logged by this function.
func Server(m string) {
	log.Printf("[server/%s] %s\n", timeString(), m)
}

// Request writes a new line to the log containing all the relevant information
// of an http.Request, if its path shouldn't be excluded.
//
// LogType is received but only OUTER_REQUEST and INNER_REQUEST should be passed.
func Request(t logType, r *http.Request, trace string) {

	// Check if the recieved url should be excluded from log.
	if _, ok := excludeUrls[r.URL.Path]; ok {
		return;
	}

	// Extract request body leaving the Reader untouched.
	body := ""
	if r.Body != nil {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ErrorFrom(err, trace)
		} else {
			body = string(buf)
			r.Body = io.NopCloser(bytes.NewReader(buf))
		}
	}

	// Write the log line to output.
	l := &customLog{
		Time: timeString(),
		TraceID: trace,
		Level: "INFO",
		Type: t,
		Method: r.Method,
		Url: r.URL.Host + r.URL.Path,
		Headers: r.Header,
		Body: stringUpperBound(trim(body), maxBodyLength),
	}
	l.print()
}

// Response allows to log an http.Response. Allows trace id to be specified.
// Logs information from the http.Request that was sent to obtain the given response.
func Response(r *http.Response, trace string) {

	// Check if the recieved url should be excluded from log.
	if _, ok := excludeUrls[r.Request.URL.Path]; ok {
		return;
	}

	// Extract request body leaving the Reader untouched.
	body := ""
	if r.Body != nil {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ErrorFrom(err, trace)
		} else {
			body = string(buf)
			r.Body = io.NopCloser(bytes.NewReader(buf))
		}
	}

	// Write the log line to output.
	l := &customLog{
		Time: timeString(),
		TraceID: trace,
		Level: "INFO",
		Type: INNER_RESPONSE,
		Method: r.Request.Method,
		Url: r.Request.URL.Host + r.Request.URL.Path,
		Status: r.StatusCode,
		Headers: r.Header,
		Body: stringUpperBound(trim(body), maxBodyLength),
	}
	l.print()
}

// ResponseWriter allows to log a controller response. Allows trace id to be specified.
// Receives the http.Request that was received that ended in the given response and logs
// information from it.
func ResponseWriter(rwBody *bytes.Buffer, rwHeaders http.Header, status int, 
	r *http.Request, trace string) {
	
	// Check if the recieved url should be excluded from log.
	if _, ok := excludeUrls[r.URL.Path]; ok {
		return;
	}

	l := &customLog{
		Time: timeString(),
		TraceID: trace,
		Level: "INFO",
		Type: OUTER_RESPONSE,
		Method: r.Method,
		Url: r.URL.Host + r.URL.Path,
		Status: status,
		Headers: rwHeaders,
		Body: stringUpperBound(rwBody.String(), maxBodyLength),
	}
	l.print()
}

// trim revomes newlines (works for Windows and Linux) and whitespace.
func trim(s string) string {
	re := regexp.MustCompile(`\r?\n`)
	resp := re.ReplaceAllString(s, " ")
	resp = strings.ReplaceAll(resp, " ", "")
	return resp
}

// stringUpperBound Returns first n characters of string or s if len(s) < n.
// This is useful for adding upper bound to strings.
//
// If n <= 0, returns s.
func stringUpperBound(s string, n int) string {
	if n <= 0 {
		return s
	}
	if len(s) < n {
		return s
	}
	return s[:n]
}

// timeString returns the current time in the following format:
// yyyy-mm-ddTHH:mm:ss.SSS
func timeString() string {
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.%03d", 
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/int(time.Millisecond))
}