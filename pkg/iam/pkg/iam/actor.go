package iam

import "github.com/alloyzeus/go-azfl/azcore"

// Actor provides information about who or what performed an action.
type Actor struct {
	// terminalID is the ID of the terminal where the action was
	// initiated from.
	terminalID TerminalID
	// userID is the ID of the user who performed the action. This might
	// be empty if the action was performed by non-user-representing agent.
	userID UserID
}

func NewActor(terminalID TerminalID, userID UserID) Actor {
	return Actor{userID: userID, terminalID: terminalID}
}

var _ azcore.SessionSubject[
	TerminalIDNum, TerminalID, UserIDNum, UserID] = Actor{}

// AZSubject is required by azcore.Subject
func (Actor) AZSessionSubject() {}

// IsRepresentingAUser is required by azcore.Subject
func (actor Actor) IsRepresentingAUser() bool {
	return actor.userID.IsStaticallyValid()
}

// TerminalID is required by azcore.Subject
func (actor Actor) TerminalID() TerminalID { return actor.terminalID }

// UserID is required by azcore.Subject
func (actor Actor) UserID() UserID { return actor.userID }
