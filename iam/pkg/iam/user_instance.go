package iam

type UserMetaService interface{}

// UserInstanceInfoServiceClientCore provides access to info about User
// instances.
//
//TODO: make this struct instances connect to IAM server and manage
// synchronization of user account states through push and pull.
type UserInstanceInfoServiceClientCore struct {
}

func (uaStateSvcClient *UserInstanceInfoServiceClientCore) GetUserInstanceInfo(
	_ CallInputContext,
	_ UserID,
) (*UserInstanceInfo, error) {
	return &UserInstanceInfo{
		RevisionNumber_: -1,
		Deletion_:       nil,
	}, nil
}
