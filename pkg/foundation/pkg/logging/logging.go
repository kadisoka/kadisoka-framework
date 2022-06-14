package logging

import (
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

//TODO: implement lookup by key which constructed from the package name
// e.g., for package example.com/mypackage, we will lookup the environment
// variables prefixed with LOG_EXAMPLE_COM_MYPACKAGE_ .
const envVarsPrefix = "LOG_"

func init() {
	if v := os.Getenv(envVarsPrefix + "TRIM_PKG_PREFIXES"); v != "" {
		parts := strings.Split(v, ",")
		ls := make([]string, 0, len(parts))
		for _, p := range parts {
			ls = append(ls, strings.TrimSpace(p))
		}
		addPackagePrefixToTrim(ls...)
	}
}

type (
	Logger = zerolog.Logger
)

// PkgLogger is a logger for a specific package. It includes the field 'pkg'
// which value is the identifier of the package.
type PkgLogger struct {
	Logger
}

// NewPkgLogger creates a logger for use within a package. This logger
// automatically adds the name of the package where this function was called,
// not when logging.
func NewPkgLogger() PkgLogger {
	// Call depth 1 because it's for the one that called NewPkgLogger
	return NewPkgLoggerWithProvidedName(ResolvePkgName(1))
}

// Packages with this prefix will be left without the prefix. This is to
// reduce noise.
var trimPackagePrefixList = []string{}
var trimPackagePrefixListMutex = sync.RWMutex{}

func addPackagePrefixToTrim(pkgPrefix ...string) int {
	if len(pkgPrefix) == 0 {
		return 0
	}

	trimPackagePrefixListMutex.Lock()
	trimPackagePrefixList = append(trimPackagePrefixList, pkgPrefix...)
	trimPackagePrefixListMutex.Unlock()

	return len(pkgPrefix)
}

// NewPkgLoggerWithProvidedName creates a package logger which field 'pkg' is
// set to the provided pkgName.
func NewPkgLoggerWithProvidedName(pkgName string) PkgLogger {
	//TODO: configurable prefix trimming
	trimPackagePrefixListMutex.RLock()
	for _, pfx := range trimPackagePrefixList {
		if strings.HasPrefix(pkgName, pfx) {
			pkgName = pkgName[len(pfx):]
			// Only the first
			break
		}
	}
	trimPackagePrefixListMutex.RUnlock()

	logCtx := newLoggerByEnv().With().Str("pkg", pkgName)
	return PkgLogger{logCtx.Logger()}
}

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
	if includeTimestampField() {
		logCtx = logCtx.Timestamp()
	}

	return logCtx.Logger()
}

func includeTimestampField() bool {
	//TODO: not just on AWS. If we are detecting an environment which
	// already providing timestamp, we should disable the timestamp
	// by default.
	if os.Getenv("AWS_EXECUTION_ENV") != "" {
		return false
	}
	return true
}

func ResolvePkgName(callDepth int) string {
	// plus one because we need to skip the call to this method
	pc, _, _, ok := runtime.Caller(callDepth + 1)
	if !ok {
		return "<unknown>"
	}

	absFuncName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(absFuncName, ".")
	partsCount := len(parts)

	if parts[partsCount-2][0] == '(' {
		// Skip the function
		return strings.Join(parts[0:partsCount-2], ".")
	}

	return strings.Join(parts[0:partsCount-1], ".")
}
