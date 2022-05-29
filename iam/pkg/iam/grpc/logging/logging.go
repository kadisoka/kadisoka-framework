package logging

import (
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api"
	foundationlog "github.com/kadisoka/kadisoka-framework/foundation/pkg/logging"
	"google.golang.org/grpc"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

// NewPkgLogger creates a logger for use within a package. This logger
// automatically adds the name of the package where this function was called,
// not when logging.
func NewPkgLogger() Logger {
	return Logger{PkgLogger: foundationlog.
		NewPkgLoggerInternal(foundationlog.CallerPkgName())}
}

// Logger wraps other logger to provide additional functionalities.
type Logger struct {
	foundationlog.PkgLogger
}

// WithContext creates a new logger which bound to a OpInputContext.
//
//TODO: don't populate the entry before the actual logging call.
func (logger Logger) WithContext(
	ctx api.OpInputContext,
) *foundationlog.Logger {
	// Implementation notes: don't panic

	if ctx == nil {
		l := logger.With().Str("class", "grpc").Logger()
		return &l
	}

	logCtx := logger.With()
	hasAuth := false

	if iamCtx, ok := ctx.(iam.OpInputContext); ok {
		if ctxAuth := iamCtx.Authorization(); ctxAuth.IsStaticallyValid() {
			logCtx = logCtx.
				Str("session", ctxAuth.Session.AZIDText()).
				Str("terminal", ctxAuth.Session.Terminal().AZIDText()).
				Str("user", ctxAuth.Session.Terminal().User().AZIDText())
			hasAuth = true
		}
	}
	if !hasAuth {
		logCtx = logCtx.
			Str("remote_addr", ctx.OpOriginInfo().Address)
	}
	if method, ok := grpc.Method(ctx); ok {
		logCtx = logCtx.
			Str("method", method)
	} else {
		logCtx = logCtx.
			Str("method", ctx.OpName())
	}

	if reqID := ctx.OpInputMetadata().ID; reqID != nil {
		logCtx = logCtx.
			Str("op_id", reqID.String())
	}

	l := logCtx.Logger()
	return &l
}
