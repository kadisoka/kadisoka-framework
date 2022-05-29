package api

import (
	"context"
	"strconv"
	"time"

	"golang.org/x/text/language"

	azcore "github.com/alloyzeus/go-azfl/azfl"
	dataerrs "github.com/alloyzeus/go-azfl/azfl/errors/data"
)

type OpInfo interface {
}

// A OpID in our implementation is used as idempotency token.
//
// A good explanation of idempotency token can be viewed here:
// https://www.youtube.com/watch?v=IP-rGJKSZ3s
//
// Check the RFC https://datatracker.ietf.org/doc/html/draft-ietf-httpapi-idempotency-key-header-01
type OpID int32

var _ azcore.ServiceMethodOpID = OpIDZero

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

func (OpID) AZServiceMethodOpID()              {}
func (opID OpID) Equal(other interface{}) bool { return opID.Equals(other) }
func (opID OpID) Equals(other interface{}) bool {
	if x, ok := other.(OpID); ok {
		return x == opID
	}
	if x, _ := other.(*OpID); x != nil {
		return *x == opID
	}
	return false
}

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
	// ID returns the idempotency token for mutating operation.
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

	// AcceptLanguage is analogous to HTTP Accept-Language header field. The
	// languages must be ordered by the human's preference.
	// If the languages comes as weighted, as found in HTTP Accept-Language,
	// sort the languages by their weights then drop the weights.
	AcceptLanguage []language.Tag

	// DateTime is the time of the device where this operation was initiated
	// from at the time the operation was posted.
	//
	// Analogous to HTTP Date header field.
	DateTime *time.Time
}
