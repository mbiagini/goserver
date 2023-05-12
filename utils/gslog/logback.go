package gslog

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// defaultMaxSize specifies the maximum size of a single log file, in MB.
// If this size is reached, a new log file will be created and the last
// one will be compressed and stored.
var defaultMaxSize int = 1

// defaultMaxBackups specifies the maximum number of log files to be
// stored once the are rolled (max size reached or new day started).
var defaultMaxBackups int = 10

// defaultMaxAge specifies the maximum number of days each log file will
// be stored before being deleted.
var defaultMaxAge int = 5

// LogFileConfig contains the configuration properties of the log files.
type LogFileConfig struct {
	MaxSize    *int `json:"max_size"`
	MaxBackups *int `json:"max_backups"`
	MaxAge     *int `json:"max_age"`
}

// ConfigureLogFile enables the configuration of the log files.
func ConfigureLogFile(c LogFileConfig) {

	maxSize := defaultMaxSize
	if c.MaxSize != nil {
		maxSize = *c.MaxSize
	}

	maxBackups := defaultMaxBackups
	if c.MaxBackups != nil {
		maxBackups = *c.MaxBackups
	}

	maxAge := defaultMaxAge
	if c.MaxAge != nil {
		maxAge = *c.MaxAge
	}

	logFilename := func() string {
		t := time.Now()
		return fmt.Sprintf("./log/go-server-%d-%02d-%02d.log", t.Year(), t.Month(), t.Day())
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename: 	logFilename(),
		MaxSize: 	maxSize,
		MaxBackups: maxBackups,
		MaxAge: 	maxAge,
		Compress: 	true,
		LocalTime: 	true,
	}

	// Configure standard logger from log package. This will be used later
	// by gslog implementation of log, which has its own prefix.
	log.SetPrefix("")
	log.SetFlags(0)

	// Fork writting into two outputs so that all logged lines are seen in
	// stdout and stored in a log file as well.
	multiWriter := io.MultiWriter(os.Stdout, lumberjackLogger)
	log.SetOutput(multiWriter)

}