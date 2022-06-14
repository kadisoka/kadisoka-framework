package logging

import (
	foundationlog "github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/logging"

	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
)

var (
	ResolvePkgName = foundationlog.ResolvePkgName
)

// NewPkgLogger creates a logger for use within a package. This logger
// automatically adds the name of the package where this function was called,
// not when logging.
func NewPkgLogger() Logger {
	// Call depth 1 because it's for the one that called NewPkgLogger
	return Logger{PkgLogger: foundationlog.
		NewPkgLoggerWithProvidedName(foundationlog.ResolvePkgName(1))}
}

// Logger is a specialized logger for logging with IAM-specific contexes.
type Logger struct {
	foundationlog.PkgLogger
}

// WithContext creates a new logger which bound to a CallInputContext.
//
// Call this method only at the logging points. It's not recommended to
// keep the returned logger around.
func (logger Logger) WithContext(
	ctx iam.CallInputContext,
) *foundationlog.Logger {
	// Implementation notes: don't panic

	if ctx == nil {
		l := logger.With().Str("module", "iam").Logger()
		return &l
	}

	logCtx := logger.With()
	logCtx = logCtx.
		Str("method", ctx.MethodName()).
		Str("res_id", ctx.ResourceID())

	hasAuth := false

	if ctxAuth := ctx.Authorization(); ctxAuth.IsStaticallyValid() {
		logCtx = logCtx.
			Str("session", ctxAuth.Session.AZIDText()).
			Str("terminal", ctxAuth.Session.Terminal().AZIDText()).
			Str("user", ctxAuth.Session.Terminal().User().AZIDText())
		hasAuth = true
	}

	if !hasAuth {
		originInfo := ctx.OriginInfo()
		logCtx = logCtx.
			Str("origin_addr", originInfo.Address)
		if originEnv := originInfo.EnvironmentString; originEnv != "" {
			logCtx = logCtx.Str("origin_env", originEnv)
		}
	}

	if idempotencyKey := ctx.CallInputMetadata().IdempotencyKey; idempotencyKey != nil {
		logCtx = logCtx.
			Str("idempotency_key", idempotencyKey.String())
	}

	l := logCtx.Logger()
	return &l
}
