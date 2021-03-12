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

type userBaseProfile struct {
	RefKey UserRefKey
	UserBaseProfileData
}

type UserBaseProfileData struct {
	RefKey          UserRefKey
	DisplayName     string
	ProfileImageURL string
	IsDeleted       bool
}

// JSONV1 models

type UserJSONV1 struct {
	ID              string `json:"id"`
	DisplayName     string `json:"display_name"`
	ProfileImageURL string `json:"profile_image_url"`
	PhoneNumber     string `json:"phone_number,omitempty"`
	EmailAddress    string `json:"email_address,omitempty"`
}

func UserJSONV1FromBaseProfile(model *UserBaseProfileData) *UserJSONV1 {
	if model == nil {
		return nil
	}
	return &UserJSONV1{
		ID:              model.RefKey.AZERText(),
		DisplayName:     model.DisplayName,
		ProfileImageURL: model.ProfileImageURL,
	}
}
