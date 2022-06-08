package api

import (
	"strconv"
	"time"

	"github.com/alloyzeus/go-azfl/azcore"
	dataerrs "github.com/alloyzeus/go-azfl/errors/data"
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
] interface {
	azcore.ServiceMethodCallInputContext[
		SessionIDNumT, SessionRefKeyT,
		TerminalIDNumT, TerminalRefKeyT,
		UserIDNumT, UserRefKeyT, SessionSubjectT, SessionT]

	// OpInputMetadata returns some details about the request.
	OpInputMetadata() OpInputMetadata
}

type OpInputMetadata struct {
	// ID returns the idempotency token for mutating operation.
	ID *OpID

	// ReceiveTime returns the time when request was accepted by
	// the handler.
	ReceiveTime time.Time
}
