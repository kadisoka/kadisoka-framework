package iam

// Actor provides information about who or what performed an action.
//
//TODO: assuming actor
type Actor struct {
	// UserRef is the ID of the user who performed the action. This might be
	// empty if the action was performed by non-user-representing agent.
	UserRef UserRefKey
	// TerminalRef is the ID of the terminal where the action was initiated
	// from.
	TerminalRef TerminalRefKey
}
