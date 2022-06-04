package logging

import (
	foundationlog "github.com/kadisoka/kadisoka-framework/foundation/pkg/logging"

	iamlog "github.com/kadisoka/kadisoka-framework/iam/pkg/iam/logging"
)

// NewPkgLogger creates a logger for use within a package. This logger
// automatically adds the name of the package where this function was called,
// not when logging.
func NewPkgLogger() Logger {
	// Call depth 1 because it's for the one that called NewPkgLogger
	return Logger{iamlog.Logger{PkgLogger: foundationlog.
		NewPkgLoggerExplicit(foundationlog.CallerPkgName(1))}}
}

// Logger wraps other logger to provide additional functionalities.
type Logger struct {
	iamlog.Logger
}
