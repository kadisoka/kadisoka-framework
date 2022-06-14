package iam

import (
	"github.com/alloyzeus/go-azfl/azcore"
	iampb "github.com/alloyzeus/go-azgrpc/azgrpc/iam/v1"
)

type UserProfileService interface {
	GetUserInfoV1(
		inputCtx CallInputContext,
		userID UserID,
	) (*iampb.UserInfoData, error)
	GetUserBaseProfile(
		inputCtx CallInputContext,
		userID UserID,
	) (*UserBaseProfileData, error)
}

type UserBaseProfileData struct {
	InstanceInfo *UserInstanceInfo

	DisplayName     string
	ProfileImageURL string
}

var _ azcore.EntityAttributes = UserBaseProfileData{}
var _ azcore.ValueObjectAssert[UserBaseProfileData] = UserBaseProfileData{}

func (profileData UserBaseProfileData) Clone() UserBaseProfileData {
	if instInfo := profileData.InstanceInfo; instInfo != nil {
		cp := profileData
		instInfoCp := instInfo.Clone()
		cp.InstanceInfo = &instInfoCp
		return cp
	}
	return profileData
}

func (UserBaseProfileData) AZAttributes()       {}
func (UserBaseProfileData) AZEntityAttributes() {}

func (profileData UserBaseProfileData) IsDeleted() bool {
	return profileData.InstanceInfo != nil && profileData.InstanceInfo.IsDeleted()
}

// JSONV1 models

type UserJSONV1 struct {
	ID           UserID                  `json:"id"`
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

func UserJSONV1FromBaseProfile(model *UserBaseProfileData, id UserID) *UserJSONV1 {
	if model == nil {
		return nil
	}
	data := UserDataJSONV1FromBaseProfile(model)
	result := UserJSONV1{
		ID: id,
	}
	if data != nil {
		result.Data = *data
	}
	return &result
}
