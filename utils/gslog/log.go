package gslog

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
)

// MaxPrefixLength specifies the number of characters that will compose the
// prefix of every log line written (excluding server logs lines). This should
// not be modified but, if really necessary, it's exported.
var MaxPrefixLength int = 30

// LogType specifies the category of a log line, diferencing logs requested
// directly by the caller from those automatically added by Handlers and
// Transports for http logging.
type logType string

const (
	MESSAGE        logType = "MESSAGE"
	OUTER_REQUEST  logType = "OUTER_REQUEST"
	OUTER_RESPONSE logType = "OUTER_RESPONSE"
	INNER_REQUEST  logType = "INNER_REQUEST"
	INNER_RESPONSE logType = "INNER_RESPONSE"
)

// LogLevel specifies the level (INFO, DEBUG, etc.) with which a log line will
// be written.
type logLevel int

const (
	DEBUG logLevel = iota
	INFO
	WARN
	ERROR
)

// customLog is the main struct of this package and contains all the information 
// to be logged. 'Method', 'Url', 'Headers' and 'Body' should be used only by Handlers
// and Transports that need to log HTTP request and response.
type customLog struct {
	Time     string              `json:"time,omitempty"`
	TraceID  string 			 `json:"trace_id,omitempty"`
	Level 	 string 			 `json:"level,omitempty"`
	Type  	 logType 			 `json:"type,omitempty"`
	Message  string 			 `json:"message,omitempty"`
	Method   string              `json:"method,omitempty"`
	Url      string              `json:"url,omitempty"`
	Status   int                 `json:"status,omitempty"`
	Headers  map[string][]string `json:"headers,omitempty"`
	Body 	 string 			 `json:"body,omitempty"`
}

// toString is an internal function to get a Log instance as a JSON string.

// Note: if json.Marshal throws an error, this code will panic.
// No application should run with logging errors.
func (l *customLog) toString() string {
	logJson, err := json.Marshal(l)
	if err != nil {
		panic(err)
	}
	return string(logJson)
} 

// Prints to output writer a new log struct. 
func (l *customLog) print() {

	// Get caller function and code line.
	// A skip value of 2 is used because the first caller to this function
	// is allways another function from this package.
	_, fn, line, ok := runtime.Caller(2)

	if !ok {
		fn = "function not available"
		line = -1
	}

	// Set and truncate the log line's prefix.
	prefix := fn + ":" + fmt.Sprint(line)
	if len(prefix) > MaxPrefixLength {
		prefix = prefix[len(prefix)-MaxPrefixLength:]
	}

	log.Printf("[%s] %s\n", prefix, l.toString())
}