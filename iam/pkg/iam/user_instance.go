package iam

type UserInstanceService interface {
	UserRefKeyService

	UserInstanceStateService

	// DeleteUserInstance deletes an instance of user entity based identfied
	// by instRefToDelete.
	//TODO: returns the revision ID of the instance.
	DeleteUserInstance(
		callCtx CallContext,
		instRefToDelete UserRefKey,
		input UserInstanceDeletionInput,
	) (deleted bool, err error)
}

type UserInstanceStateService interface {
	// GetUserInstanceState checks if the provided user ID is valid and whether
	// the instance is deleted.
	//
	// This method returns nil if the userRef is not referencing to any valid
	// user instance.
	GetUserInstanceState(
		/*callCtx CallContext,*/ //TODO: call context
		userRef UserRefKey,
	) (*UserInstanceStateData, error)
}

type UserInstanceStateData struct {
	Deletion *UserInstanceDeletionData
}

func (instState UserInstanceStateData) IsInstanceActive() bool {
	return instState.Deletion == nil && !instState.Deletion.Deleted
}

//TODO: make this struct instances connect to IAM server and manage
// synchronization of user account states.
type UserInstanceStateServiceClientCore struct {
}

func (uaStateSvcClient *UserInstanceStateServiceClientCore) GetUserInstanceState(
	_ UserRefKey,
) (*UserInstanceStateData, error) {
	return &UserInstanceStateData{Deletion: nil}, nil
}

//TODO: reason and notes
type UserInstanceDeletionInput struct {
	DeletionNotes string
}

type UserInstanceDeletionData struct {
	Deleted       bool
	DeletionNotes string
}
