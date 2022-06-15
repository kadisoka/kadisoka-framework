package iam

import "github.com/alloyzeus/go-azfl/azcore"

// AuthorizationSubject provides information about who or what an authorization
// represents.
type AuthorizationSubject struct {
	// terminalID is the ID of the terminal where the action was
	// initiated from.
	terminalID TerminalID
	// userID is the ID of the user who performed the action. This might
	// be empty if the action was performed by non-user-representing agent.
	userID UserID
}

func NewAuthorizationSubject(
	terminalID TerminalID,
	userID UserID,
) AuthorizationSubject {
	return AuthorizationSubject{userID: userID, terminalID: terminalID}
}

var _ azcore.SessionSubject[
	TerminalIDNum, TerminalID, UserIDNum, UserID] = AuthorizationSubject{}

// AZSubject is required by azcore.Subject
func (AuthorizationSubject) AZSessionSubject() {}

// IsRepresentingAUser is required by azcore.Subject
func (actor AuthorizationSubject) IsRepresentingAUser() bool {
	return actor.userID.IsStaticallyValid()
}

// TerminalID is required by azcore.Subject
func (actor AuthorizationSubject) TerminalID() TerminalID { return actor.terminalID }

// UserID is required by azcore.Subject
func (actor AuthorizationSubject) UserID() UserID { return actor.userID }
