package logging

import (
	foundationlog "github.com/kadisoka/kadisoka-framework/foundation/pkg/logging"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

// NewPkgLogger creates a logger for use within a package. This logger
// automatically adds the name of the package where this function was called,
// not when logging.
func NewPkgLogger() Logger {
	return Logger{PkgLogger: foundationlog.
		NewPkgLoggerInternal(foundationlog.CallerPkgName())}
}

// Logger is a specialized logger for logging with IAM-specific contexes.
type Logger struct {
	foundationlog.PkgLogger
}

// WithContext creates a new logger which bound to a OpInputContext.
//
// Call this method only at the logging points. It's not recommended to
// keep the returned logger around.
func (logger Logger) WithContext(
	ctx iam.OpInputContext,
) *foundationlog.Logger {
	// Implementation notes: don't panic

	if ctx == nil {
		l := logger.With().Str("class", "iam").Logger()
		return &l
	}

	logCtx := logger.With()
	hasAuth := false

	if ctxAuth := ctx.Authorization(); ctxAuth.IsValid() {
		logCtx = logCtx.
			Str("session", ctxAuth.Session.AZIDText()).
			Str("terminal", ctxAuth.Session.Terminal().AZIDText()).
			Str("user", ctxAuth.Session.Terminal().User().AZIDText())
		hasAuth = true
	}
	if !hasAuth {
		//TODO: generalized remote IP resolver
	}
	logCtx = logCtx.
		Str("method", ctx.OpName())

	if reqID := ctx.OpInputMetadata().ID; reqID != nil {
		logCtx = logCtx.
			Str("op_id", reqID.String())
	}

	l := logCtx.Logger()
	return &l
}
