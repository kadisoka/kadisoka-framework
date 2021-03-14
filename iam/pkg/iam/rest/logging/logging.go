package logging

import (
	"net/http"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api/rest"
	foundationlog "github.com/kadisoka/kadisoka-framework/foundation/pkg/logging"
	"github.com/tomasen/realip"

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

// WithContext creates a new logger which bound to a RequestContext.
//
//TODO: don't populate the entry before the actual logging call.
func (logger Logger) WithContext(
	ctx rest.RequestContext,
) *foundationlog.Logger {
	// Implementation notes: don't panic

	if ctx == nil {
		l := logger.With().Str("class", "rest").Logger()
		return &l
	}

	logCtx := logger.With()
	hasAuth := false

	if iamCtx, _ := ctx.(iam.CallContext); iamCtx != nil {
		if ctxAuth := iamCtx.Authorization(); ctxAuth.IsValid() {
			logCtx = logCtx.
				Str("session", ctxAuth.Session.AZERText()).
				Str("terminal", ctxAuth.Session.Terminal().AZERText()).
				Str("user", ctxAuth.Session.Terminal().User().AZERText())
		}
	}

	if req := ctx.HTTPRequest(); req != nil {
		var urlStr string
		if req.URL != nil {
			urlStr = req.URL.String()
		}
		logCtx = logCtx.
			Str("method", req.Method).
			Str("url", urlStr)
		if !hasAuth {
			logCtx = logCtx.
				Str("remote_addr", ctx.OriginInfo().Address).
				Str("user_agent", req.UserAgent())
		}
	}

	if reqID := ctx.RequestInfo().ID; reqID != nil {
		logCtx = logCtx.
			Str("request_id", reqID.String())
	}

	l := logCtx.Logger()
	return &l
}

// WithRequest creates a log entry with some fields from the request.
func (logger Logger) WithRequest(
	req *http.Request,
) *foundationlog.Logger {
	// Implementation notes: don't panic

	if req == nil {
		return &logger.Logger
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
		Str("method", req.Method).
		Str("url", urlStr).
		Str("remote_addr", remoteAddr).
		Str("user_agent", req.UserAgent()).
		Logger()
	return &l
}
