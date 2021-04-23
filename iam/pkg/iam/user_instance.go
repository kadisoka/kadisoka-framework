package iam

type UserMetaService interface{}

//TODO: make this struct instances connect to IAM server and manage
// synchronization of user account states.
type UserInstanceInfoServiceClientCore struct {
}

func (uaStateSvcClient *UserInstanceInfoServiceClientCore) GetUserInstanceInfo(
	_ CallContext,
	_ UserRefKey,
) (*UserInstanceInfo, error) {
	return &UserInstanceInfo{RevisionNumber: -1, Deletion: nil}, nil
}
