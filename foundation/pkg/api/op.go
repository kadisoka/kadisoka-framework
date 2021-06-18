package api

import (
	"context"
	"strconv"
	"time"

	"golang.org/x/text/language"

	dataerrs "github.com/alloyzeus/go-azfl/azfl/errors/data"
)

// A OpID in our implementation is used as idempotency token.
//
// A good explanation of idempotency token can be viewed here:
// https://www.youtube.com/watch?v=IP-rGJKSZ3s
type OpID int32

const OpIDZero = OpID(0)

func OpIDFromString(opIDStr string) (OpID, error) {
	u, err := strconv.ParseInt(opIDStr, 10, 32)
	if err != nil {
		return OpIDZero, dataerrs.Malformed(err)
	}
	i := OpID(u)
	if isOpIDSound(i) {
		return OpIDZero, dataerrs.ErrMalformed
	}
	return i, nil
}

func isOpIDSound(opID OpID) bool {
	return opID > 0
}

func (opID OpID) String() string { return strconv.FormatInt(int64(opID), 10) }

// OpInputContext holds information obtained from the request. This information
// are generally obtained from the request's metadata (e.g., HTTP request
// header).
//TODO: proxied context.
type OpInputContext interface {
	context.Context

	// OpName returns the name of the method or the endpoint.
	//
	// For HTTP, this method returns the value as "<HTTP_METHOD> <URL>", e.g.,
	// GET /users/me
	//
	OpName() string

	// OpInputMetadata returns some details about the request.
	OpInputMetadata() OpInputMetadata

	// OpOriginInfo returns some details about the caller.
	OpOriginInfo() OpOriginInfo
}

type OpInputMetadata struct {
	// ID returns the idempotency token if provided in the call request.
	ID *OpID

	// ReceiveTime returns the time when request was accepted by
	// the handler.
	ReceiveTime time.Time
}

type OpOriginInfo struct {
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
