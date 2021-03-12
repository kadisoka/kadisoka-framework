package iam

type UserAccountService interface {
	UserIDService

	UserAccountStateService

	// DeleteUserAccount deletes an user account based identfied by userIDToDelete.
	//TODO: returns the revision ID of the account.
	DeleteUserAccount(
		callCtx CallContext,
		userRefToDelete UserRefKey,
		input UserAccountDeleteInput,
	) (deleted bool, err error)
}

type UserAccountStateService interface {
	// GetUserAccountState checks if the provided user ID is valid and whether
	// the account is deleted.
	//
	// This method returns nil if the userRef is not referencing to any valid
	// user account.
	GetUserAccountState(
		/*callCtx CallContext,*/ //TODO: call context
		userRef UserRefKey,
	) (*UserAccountState, error)
}

type UserAccountState struct {
	Deleted bool
}

func (uaState UserAccountState) IsAccountActive() bool {
	return !uaState.Deleted
}

//TODO: make this struct instances connect to IAM server and manage
// synchronization of user account states.
type UserAccountStateServiceClientCore struct {
}

func (uaStateSvcClient *UserAccountStateServiceClientCore) GetUserAccountState(
	_ UserRefKey,
) (*UserAccountState, error) {
	return &UserAccountState{false}, nil
}

//TODO: reason and comment
type UserAccountDeleteInput struct {
	DeletionNotes string
}
