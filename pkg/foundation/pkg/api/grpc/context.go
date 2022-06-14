package grpc

import (
	"github.com/alloyzeus/go-azfl/azcore"

	"github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/api"
)

type CallInputContext[
	SessionIDNumT azcore.SessionIDNum, SessionIDT azcore.SessionID[SessionIDNumT],
	TerminalIDNumT azcore.TerminalIDNum, TerminalIDT azcore.TerminalID[TerminalIDNumT],
	UserIDNumT azcore.UserIDNum, UserIDT azcore.UserID[UserIDNumT],
	SessionSubjectT azcore.SessionSubject[
		TerminalIDNumT, TerminalIDT,
		UserIDNumT, UserIDT],
	SessionT azcore.Session[
		SessionIDNumT, SessionIDT,
		TerminalIDNumT, TerminalIDT,
		UserIDNumT, UserIDT,
		SessionSubjectT],
	IdempotencyKeyT azcore.ServiceMethodIdempotencyKey,
] interface {
	api.CallInputContext[
		SessionIDNumT, SessionIDT,
		TerminalIDNumT, TerminalIDT,
		UserIDNumT, UserIDT, SessionSubjectT, SessionT, IdempotencyKeyT]
}
