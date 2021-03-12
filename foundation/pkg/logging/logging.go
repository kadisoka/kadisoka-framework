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
	return NewPkgLoggerInternal(CallerPkgName())
}

// NewPkgLoggerInternal creates a package logger which field 'pkg' is
// set to the provided name.
func NewPkgLoggerInternal(name string) PkgLogger {
	//TODO: configurable prefix trimming
	name = strings.TrimPrefix(name, "github.com/kadisoka/")
	logCtx := newLoggerByEnv().With().Str("pkg", name)
	return PkgLogger{logCtx.Logger()}
}

//TODO: implement lookup by key which constructed from the package name
// e.g., for package example.com/mypackage, we will lookup the environment
// variables prefixed with LOG_EXAMPLE_COM_MYPACKAGE_ .
const envVarsPrefix = "LOG_"

func newLogger() Logger {
	//TODO: for the default, we detect the environment we are running on.
	// e.g., if it's local, it's pretty as the default.
	if logPretty := os.Getenv(envVarsPrefix + "PRETTY"); logPretty == "true" {
		return zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}
	return zerolog.New(os.Stderr)
}

func newLoggerByEnv() Logger {
	logger := newLogger()

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

func CallerPkgName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "<unknown>"
	}
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	partsCount := len(parts)
	pkgPath := ""
	funcName := parts[partsCount-1]
	if parts[partsCount-2][0] == '(' {
		funcName = parts[partsCount-2] + "." + funcName
		pkgPath = strings.Join(parts[0:partsCount-2], ".")
	} else {
		pkgPath = strings.Join(parts[0:partsCount-1], ".")
	}
	return pkgPath
}
