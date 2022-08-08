package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Logger contains the data necessary to log status and error messages in a standard way
type Logger struct {
	Environment string
	Prefix      string
	infoLog     *log.Logger
	errLog      *log.Logger
	errProvider ErrorProvider
}

// NewLogger creates a new logger from the service name and environment name
func NewLogger(service string, environment string) *Logger {
	return &Logger{
		infoLog:     log.New(os.Stdout, "", log.LstdFlags),
		errLog:      log.New(os.Stderr, "", log.LstdFlags),
		errProvider: ErrorProvider{SkipFrames: 2, PackageBase: "Woody1193"},
		Environment: environment,
		Prefix:      fmt.Sprintf("[%s][%s] ", environment, service),
	}
}

// ChangeFrame creates a new logger from an existing logger with a different
// number of frames to skip, allowing for errors to be referenced from a different
// part of the call stack than the default logger
func (logger *Logger) ChangeFrame(skipFrames int) *Logger {
	return &Logger{
		infoLog:     logger.infoLog,
		errLog:      logger.errLog,
		errProvider: ErrorProvider{SkipFrames: skipFrames, PackageBase: "xefino"},
		Environment: logger.Environment,
		Prefix:      logger.Prefix,
	}
}

// Log a message to the standard output
func (logger *Logger) Log(message string, args ...interface{}) {
	logger.infoLog.Printf(fmt.Sprintf("[Info]%s%s", logger.Prefix, fmt.Sprintf(message, args...)))
}

// Generate and log an error from the inner error and message. The
// resulting error will be returned for use by the caller
func (logger *Logger) Error(inner error, message string, args ...interface{}) *Error {
	err := logger.errProvider.GenerateError(logger.Environment, inner, message, args...)
	logger.errLog.Printf("[Error]%s%v", logger.Prefix, err)
	return err
}

// Discard is primarily used for testing. It operates by setting the
// output from all the loggers to a discard stream so that no logging
// actually appears on the screen
func (logger *Logger) Discard() {
	logger.infoLog.SetOutput(ioutil.Discard)
	logger.errLog.SetOutput(ioutil.Discard)
}
