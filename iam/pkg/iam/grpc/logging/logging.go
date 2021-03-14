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

// WithContext creates a new logger which bound to a CallContext.
//
//TODO: don't populate the entry before the actual logging call.
func (logger Logger) WithContext(
	ctx api.CallContext,
) *foundationlog.Logger {
	// Implementation notes: don't panic

	if ctx == nil {
		l := logger.With().Str("class", "grpc").Logger()
		return &l
	}

	logCtx := logger.With()
	hasAuth := false

	if iamCtx, ok := ctx.(iam.CallContext); ok {
		if ctxAuth := iamCtx.Authorization(); ctxAuth.IsValid() {
			logCtx = logCtx.
				Str("session", ctxAuth.Session.AZERText()).
				Str("terminal", ctxAuth.Session.Terminal().AZERText()).
				Str("user", ctxAuth.Session.Terminal().User().AZERText())
			hasAuth = true
		}
	}
	if !hasAuth {
		logCtx = logCtx.
			Str("remote_addr", ctx.OriginInfo().Address)
	}
	if method, ok := grpc.Method(ctx); ok {
		logCtx = logCtx.
			Str("method", method)
	} else {
		logCtx = logCtx.
			Str("method", ctx.MethodName())
	}

	if reqID := ctx.RequestInfo().ID; reqID != nil {
		logCtx = logCtx.
			Str("request_id", reqID.String())
	}

	l := logCtx.Logger()
	return &l
}
