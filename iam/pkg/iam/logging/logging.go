package logging

import (
	foundationlog "github.com/kadisoka/foundation/pkg/logging"

	"github.com/kadisoka/iam/pkg/iam"
)

// NewPkgLogger creates a logger for use within a package. This logger
// automatically adds the name of the package where this function was called,
// not when logging.
func NewPkgLogger() Logger {
	return Logger{PkgLogger: foundationlog.NewPkgLoggerInternal(foundationlog.CallerPkgName())}
}

// Logger is a specialized logger for logging with IAM-specific contexes.
type Logger struct {
	foundationlog.PkgLogger
}

// WithContext creates a new logger which bound to a CallContext.
//
// Call this method only at the logging points. It's not recommended to
// keep the returned logger around.
func (logger Logger) WithContext(
	ctx iam.CallContext,
) *foundationlog.Logger {
	// Implementation notes: don't panic

	if ctx == nil {
		l := logger.With().Str("class", "iam").Logger()
		return &l
	}

	logCtx := logger.With()
	hasAuth := false

	if iamCtx, ok := ctx.(iam.CallContext); ok {
		if authCtx := iamCtx.Authorization(); authCtx.IsValid() {
			logCtx = logCtx.
				Str("user", authCtx.UserID.String()).
				Str("terminal", authCtx.TerminalID().String()).
				Str("auth", authCtx.AuthorizationID.String())
			hasAuth = true
		}
	}
	if !hasAuth {
		//TODO: generalized remote IP resolver
	}
	logCtx = logCtx.
		Str("method", ctx.MethodName())

	if reqID := ctx.RequestID(); reqID != nil {
		logCtx = logCtx.
			Str("request_id", reqID.String())
	}

	l := logCtx.Logger()
	return &l
}
