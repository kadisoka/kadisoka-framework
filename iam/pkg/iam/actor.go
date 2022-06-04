package iam

import azcore "github.com/alloyzeus/go-azfl/azfl"

// Actor provides information about who or what performed an action.
//
//TODO: assuming actor
type Actor struct {
	// UserRef is the RefKey of the user who performed the action. This might
	// be empty if the action was performed by non-user-representing agent.
	UserRef UserRefKey
	// TerminalRef is the RefKey of the terminal where the action was
	// initiated from.
	TerminalRef TerminalRefKey
}

var _ azcore.SessionSubject[
	TerminalIDNum, TerminalRefKey, UserIDNum, UserRefKey] = Actor{}

// AZSubject is required by azcore.Subject
func (Actor) AZSessionSubject() {}

// IsRepresentingAUser is required by azcore.Subject
func (actor Actor) IsRepresentingAUser() bool {
	return actor.UserRef.IsStaticallyValid()
}

// TerminalRefKey is required by azcore.Subject
func (actor Actor) TerminalRefKey() TerminalRefKey { return actor.TerminalRef }

// UserRefKey is required by azcore.Subject
func (actor Actor) UserRefKey() UserRefKey { return actor.UserRef }
