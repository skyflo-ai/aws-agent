package logger

import "log"

// Info logs an informational message.
func Info(format string, args ...interface{}) {
	log.Printf("INFO: "+format, args...)
}

// Error logs an error message.
func Error(format string, args ...interface{}) {
	log.Printf("ERROR: "+format, args...)
}
