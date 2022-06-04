package logging

import (
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type (
	Logger = zerolog.Logger
)

type PkgLogger struct {
	Logger
}

// NewPkgLogger creates a logger for use within a package. This logger
// automatically adds the name of the package where this function was called,
// not when logging.
func NewPkgLogger() PkgLogger {
	// Call depth 1 because it's for the one that called NewPkgLogger
	return NewPkgLoggerInternal(CallerPkgName(1))
}

// NewPkgLoggerInternal creates a package logger which field 'pkg' is
// set to the provided name.
func NewPkgLoggerInternal(name string) PkgLogger {
	//TODO: configurable prefix trimming
	name = strings.TrimPrefix(name, "github.com/kadisoka/kadisoka-framework/")
	logCtx := newLoggerByEnv().With().Str("pkg", name)
	return PkgLogger{logCtx.Logger()}
}

//TODO: implement lookup by key which constructed from the package name
// e.g., for package example.com/mypackage, we will lookup the environment
// variables prefixed with LOG_EXAMPLE_COM_MYPACKAGE_ .
const envVarsPrefix = "LOG_"

func newLogger(prettyLog bool) Logger {
	if prettyLog {
		return zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		})
	}
	return zerolog.New(os.Stderr)
}

func newLoggerByEnv() Logger {
	prettyLog := false
	//TODO: for the default, we detect the environment we are running on.
	// e.g., if it's local, it's pretty as the default.
	if v := os.Getenv(envVarsPrefix + "PRETTY"); v == "true" {
		prettyLog = true
	}

	logger := newLogger(prettyLog)

	if logLevelStr := os.Getenv(envVarsPrefix + "LEVEL"); logLevelStr != "" {
		logLevel, err := zerolog.ParseLevel(logLevelStr)
		if err != nil {
			panic(err)
		}
		logger = logger.Level(logLevel)
	}

	logCtx := logger.With()
	if os.Getenv("AWS_EXECUTION_ENV") != "" {
		//TODO: not just on AWS. If we are detecting an environment which
		// already providing timestamp, we should disable the timestamp
		// by default.
	} else {
		logCtx = logCtx.Timestamp()
	}

	return logCtx.Logger()
}

func CallerPkgName(callDepth int) string {
	// plus one because we need to skip the call to this method
	pc, _, _, ok := runtime.Caller(callDepth + 1)
	if !ok {
		return "<unknown>"
	}

	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	partsCount := len(parts)

	if parts[partsCount-2][0] == '(' {
		// Skip the function
		return strings.Join(parts[0:partsCount-2], ".")
	}

	return strings.Join(parts[0:partsCount-1], ".")
}
