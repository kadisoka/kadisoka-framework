package iam

// UserInstanceService is a service which provides
// methods for manipulating entity instances.
type UserInstanceService interface {
	UserRefKeyService

	UserInstanceInfoService

	// DeleteUserInstance deletes an instance of user entity based identfied
	// by instRefToDelete.
	//TODO: returns the revision ID of the instance.
	DeleteUserInstance(
		callCtx CallContext,
		instRefToDelete UserRefKey,
		input UserInstanceDeletionInput,
	) (deleted bool, err error)
}

// UserInstanceInfoService is a service which
// provides an access to instances metadata.
type UserInstanceInfoService interface {
	// GetUserInstanceInfo checks if the provided
	// ref-key is valid and whether the instance is deleted.
	//
	// This method returns nil if the refKey is not referencing to any valid
	// instance.
	GetUserInstanceInfo(
		/*callCtx CallContext,*/ //TODO: call context
		refKey UserRefKey,
	) (*UserInstanceInfo, error)
}

// UserInstanceService holds information about
// an instance of User.
type UserInstanceInfo struct {
	RevisionNumber int32

	// Deletion holds information about the deletion of the instance. If
	// the instance has not been deleted, this field value will be nil.
	Deletion *UserInstanceDeletionInfo
}

// IsActive returns true if the instance is considered as active.
func (instInfo UserInstanceInfo) IsActive() bool {
	// Note: we will check other flags in the future, but that's said,
	// deleted instance is considered inactive.
	return !instInfo.IsDeleted()
}

// IsDeleted returns true if the instance was deleted.
func (instInfo UserInstanceInfo) IsDeleted() bool {
	return instInfo.Deletion != nil && instInfo.Deletion.Deleted
}

// UserInstanceDeletionInfo holds information about
// the deletion of an instance if the instance has been deleted.
type UserInstanceDeletionInfo struct {
	Deleted       bool
	DeletionNotes string
}

//TODO: reason and notes
type UserInstanceDeletionInput struct {
	DeletionNotes string
}

//TODO: make this struct instances connect to IAM server and manage
// synchronization of user account states.
type UserInstanceStateServiceClientCore struct {
}

func (uaStateSvcClient *UserInstanceStateServiceClientCore) GetUserInstanceInfo(
	_ UserRefKey,
) (*UserInstanceInfo, error) {
	return &UserInstanceInfo{RevisionNumber: -1, Deletion: nil}, nil
}
