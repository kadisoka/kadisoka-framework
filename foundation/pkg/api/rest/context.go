package rest

import (
	"net/http"

	"github.com/alloyzeus/go-azfl/azcore"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api"
)

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
	api.OpInputContext[
		SessionIDNumT, SessionRefKeyT,
		TerminalIDNumT, TerminalRefKeyT,
		UserIDNumT, UserRefKeyT, SessionSubjectT, SessionT]

	HTTPRequest() *http.Request
}
