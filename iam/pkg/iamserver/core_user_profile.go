package iamserver

import (
	"database/sql"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	iampb "github.com/rez-go/crux-apis/crux/iam/v1"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

// Interface conformance assertion.
var _ iam.UserProfileService = &Core{}

const userProfileDisplayNameTableName = "user_display_name_dt"
const userProfileImageKeyTableName = "user_profile_image_key_dt"

func (core *Core) GetUserBaseProfile(
	callCtx iam.CallContext,
	userID iam.UserRefKey,
) (*iam.UserBaseProfileData, error) {
	if callCtx == nil {
		return nil, errors.ArgMsg("callCtx", "missing")
	}
	//TODO(exa): ensure that the context user has the privilege

	var user iam.UserBaseProfileData
	var displayName *string
	var profileImageURL *string

	err := core.db.
		QueryRow(
			`SELECT ua.id, `+
				`CASE WHEN ua.d_ts IS NULL THEN false ELSE true END AS is_deleted, `+
				`udn.display_name, upiu.profile_image_key `+
				`FROM `+userTableName+` AS ua `+
				`LEFT JOIN `+userProfileDisplayNameTableName+` udn ON udn.user_id = ua.id `+
				`AND udn.d_ts IS NULL `+
				`LEFT JOIN `+userProfileImageKeyTableName+` upiu ON upiu.user_id = ua.id `+
				`AND upiu.d_ts IS NULL `+
				`WHERE ua.id = $1`,
			userID).
		Scan(&user.RefKey, &user.IsDeleted, &displayName, &profileImageURL)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			return nil, err
		}
	}

	if displayName != nil {
		user.DisplayName = *displayName
	}
	if profileImageURL != nil {
		user.ProfileImageURL = core.BuildUserProfileImageURL(*profileImageURL)
	}

	return &user, nil
}

func (core *Core) GetUserInfoV1(
	callCtx iam.CallContext,
	userID iam.UserRefKey,
) (*iampb.UserInfoData, error) {
	//TODO: access control

	userBaseProfile, err := core.
		GetUserBaseProfile(callCtx, userID)
	if err != nil {
		panic(err)
	}
	baseProfile := &iampb.UserBaseProfileData{
		DisplayName:     userBaseProfile.DisplayName,
		ProfileImageUrl: userBaseProfile.ProfileImageURL,
	}

	var deactivation *iampb.UserAccountDeactivationData
	if userBaseProfile.IsDeleted {
		deactivation = &iampb.UserAccountDeactivationData{
			Deactivated: true,
		}
	}
	accountInfo := &iampb.UserAccountInfoData{
		Verification: &iampb.UserAccountVerificationData{
			Verified: true, //TODO: actual value
		},
		Deactivation: deactivation,
	}

	contactInfo, err := core.
		GetUserContactInformation(callCtx, userID)
	if err != nil {
		panic(err)
	}

	return &iampb.UserInfoData{
		AccountInfo: accountInfo,
		BaseProfile: baseProfile,
		ContactInfo: contactInfo,
	}, nil
}

func (core *Core) isUserProfileImageURLAllowed(profileImageURL string) bool {
	//TODO(exa): limit profile image url to certain hosts or keep only the filename
	return true
}
