package api

import (
	"bytes"
	"time"

	"github.com/alloyzeus/go-azfl/azcore"
	dataerrs "github.com/alloyzeus/go-azfl/errors/data"
	"github.com/google/uuid"
)

type OpInfo interface {
}

// A IdempotencyKey in our implementation is used as idempotency token.
//
// A good explanation of idempotency token can be viewed here:
// https://www.youtube.com/watch?v=IP-rGJKSZ3s
//
// Check the RFC https://datatracker.ietf.org/doc/html/draft-ietf-httpapi-idempotency-key-header-01
type IdempotencyKey uuid.UUID

var _ azcore.ServiceMethodIdempotencyKey = _IdempotencyKeyZero

var _IdempotencyKeyZero = IdempotencyKey(uuid.Nil)

func IdempotencyKeyZero() IdempotencyKey { return IdempotencyKey(uuid.Nil) }

func IdempotencyKeyFromString(idempotencyKeyStr string) (IdempotencyKey, error) {
	u, err := uuid.Parse(idempotencyKeyStr)
	if err != nil {
		return IdempotencyKeyZero(), dataerrs.Malformed(err)
	}
	i := IdempotencyKey(u)
	if isIdempotencyStaticallyValid(i) {
		return IdempotencyKeyZero(), dataerrs.ErrMalformed
	}
	return i, nil
}

func isIdempotencyStaticallyValid(idempotencyKey IdempotencyKey) bool {
	//TODO: more checks?
	asUUID := uuid.UUID(idempotencyKey)
	return !bytes.Equal(idempotencyKey[:], uuid.Nil[:]) &&
		asUUID.Version() == uuid.Version(4) &&
		asUUID.Variant() == uuid.RFC4122
}

func (idempotencyKey IdempotencyKey) String() string { return uuid.UUID(idempotencyKey).String() }

func (IdempotencyKey) AZServiceMethodIdempotencyKey() {}
func (idempotencyKey IdempotencyKey) Equal(other interface{}) bool {
	return idempotencyKey.Equals(other)
}
func (idempotencyKey IdempotencyKey) Equals(other interface{}) bool {
	if x, ok := other.(IdempotencyKey); ok {
		return bytes.Equal(idempotencyKey[:], x[:])
	}
	if x, _ := other.(*IdempotencyKey); x != nil {
		return bytes.Equal(idempotencyKey[:], (*x)[:])
	}
	return false
}

// OpInputContext holds information obtained from the request. This information
// are generally obtained from the request's metadata (e.g., HTTP request
// header).
//TODO: proxied context.
type OpInputContext[
	SessionIDNumT azcore.SessionIDNum, SessionRefKeyT azcore.SessionRefKey[SessionIDNumT],
	TerminalIDNumT azcore.TerminalIDNum, TerminalRefKeyT azcore.TerminalRefKey[TerminalIDNumT],
	UserIDNumT azcore.UserIDNum, UserRefKeyT azcore.UserRefKey[UserIDNumT],
	SessionSubjectT azcore.SessionSubject[
		TerminalIDNumT, TerminalRefKeyT,
		UserIDNumT, UserRefKeyT],
	SessionT azcore.Session[
		SessionIDNumT, SessionRefKeyT,
		TerminalIDNumT, TerminalRefKeyT,
		UserIDNumT, UserRefKeyT,
		SessionSubjectT],
	IdempotencyKeyT azcore.ServiceMethodIdempotencyKey,
] interface {
	azcore.ServiceMethodCallInputContext[
		SessionIDNumT, SessionRefKeyT,
		TerminalIDNumT, TerminalRefKeyT,
		UserIDNumT, UserRefKeyT, SessionSubjectT, SessionT, IdempotencyKeyT]

	// OpInputMetadata returns some details about the request.
	OpInputMetadata() OpInputMetadata
}

type OpInputMetadata struct {
	// IdempotencyKey returns the idempotency token for mutating operation.
	IdempotencyKey *IdempotencyKey

	// ReceiveTime returns the time when request was accepted by
	// the handler.
	ReceiveTime time.Time
}
