package api

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/language"
)

// A RequestID in our implementation is used as idempotency token.
//
// A good explanation of idempotency token can be viewed here:
// https://www.youtube.com/watch?v=IP-rGJKSZ3s
type RequestID = uuid.UUID

// CallContext holds information obtained from the request. This information
// are generally obtained from the request's metadata (e.g., HTTP request
// header).
//TODO: proxied context.
type CallContext interface {
	context.Context

	// MethodName returns the name of the method this call is directed to.
	//
	// For HTTP, this method returns the value as "<HTTP_METHOD> <URL>", e.g.,
	// GET /users/me
	//
	MethodName() string

	// RequestInfo returns some details about the request.
	RequestInfo() CallRequestInfo

	// OriginInfo returns some details about the caller.
	OriginInfo() CallOriginInfo
}

type CallRequestInfo struct {
	// ID returns the idempotency token if provided in the call request.
	ID *RequestID

	// ReceiveTime returns the time when request was accepted by
	// the handler.
	ReceiveTime time.Time
}

type CallOriginInfo struct {
	// Address returns the IP address or hostname where this call was initiated
	// from. This field might be empty if it's not possible to resolve
	// the address (e.g., the server is behind a proxy or a load-balancer and
	// they didn't forward the the origin IP).
	Address string

	// EnvironmentString returns some details of the environment,
	// might include application's version information, where the application
	// which made the request runs on. For web app, this method usually
	// returns the browser's user-agent string.
	EnvironmentString string

	AcceptLanguage []language.Tag
}
