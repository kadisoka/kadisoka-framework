package iam

import (
	iampb "github.com/rez-go/crux-apis/crux/iam/v1"
)

type UserProfileService interface {
	GetUserInfoV1(
		callCtx CallContext,
		userRef UserRefKey,
	) (*iampb.UserInfoData, error)
	GetUserBaseProfile(
		callCtx CallContext,
		userRef UserRefKey,
	) (*UserBaseProfileData, error)
}

type UserBaseProfileData struct {
	InstanceInfo *UserInstanceInfo

	DisplayName     string
	ProfileImageURL string
}

func (aggregateData UserBaseProfileData) IsDeleted() bool {
	return aggregateData.InstanceInfo != nil && aggregateData.InstanceInfo.IsDeleted()
}

// JSONV1 models

type UserJSONV1 struct {
	ID           UserRefKey              `json:"id"`
	InstanceInfo *UserInstanceInfoJSONV1 `json:"instance_info"`
	Data         UserDataJSONV1          `json:"data"`
}

type UserInstanceInfoJSONV1 struct {
}

type UserDataJSONV1 struct {
	DisplayName     string `json:"display_name"`
	ProfileImageURL string `json:"profile_image_url"`
	PhoneNumber     string `json:"phone_number,omitempty"`
	EmailAddress    string `json:"email_address,omitempty"`
}

func UserDataJSONV1FromBaseProfile(model *UserBaseProfileData) *UserDataJSONV1 {
	if model == nil {
		return nil
	}
	return &UserDataJSONV1{
		DisplayName:     model.DisplayName,
		ProfileImageURL: model.ProfileImageURL,
	}
}

func UserJSONV1FromBaseProfile(model *UserBaseProfileData, refKey UserRefKey) *UserJSONV1 {
	if model == nil {
		return nil
	}
	data := UserDataJSONV1FromBaseProfile(model)
	result := UserJSONV1{
		ID: refKey,
	}
	if data != nil {
		result.Data = *data
	}
	return &result
}
