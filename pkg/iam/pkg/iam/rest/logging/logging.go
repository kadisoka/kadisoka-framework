package logging

import (
	"net/http"

	"github.com/tomasen/realip"

	foundationlog "github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/logging"

	iamlog "github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam/logging"
)

var (
	ResolvePkgName = foundationlog.ResolvePkgName
)

// NewPkgLogger creates a logger for use within a package. This logger
// automatically adds the name of the package where this function was called,
// not when logging.
func NewPkgLogger() Logger {
	// Call depth 1 because it's for the one that called NewPkgLogger
	return Logger{iamlog.Logger{PkgLogger: foundationlog.
		NewPkgLoggerWithProvidedName(foundationlog.ResolvePkgName(1))}}
}

// Logger wraps other logger to provide additional functionalities.
type Logger struct {
	iamlog.Logger
}

// WithRequest creates a log entry with some fields from the request.
func (logger Logger) WithRequest(
	req *http.Request,
) *foundationlog.Logger {
	// Implementation notes: don't panic

	if req == nil {
		return &logger.Logger.Logger
	}

	var urlStr string
	if req.URL != nil {
		urlStr = req.URL.String()
	}

	remoteAddr := realip.FromRequest(req)
	if remoteAddr == "" {
		remoteAddr = req.RemoteAddr
	}

	l := logger.With().
		Str("method", req.Method+" "+urlStr).
		Str("origin_addr", remoteAddr).
		Str("origin_env", req.UserAgent()).
		Logger()
	return &l
}
