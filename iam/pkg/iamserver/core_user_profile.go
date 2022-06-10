package iamserver

import (
	"database/sql"

	"github.com/alloyzeus/go-azfl/errors"
	iampb "github.com/rez-go/crux-apis/crux/iam/v1"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

// Interface conformance assertion.
var _ iam.UserProfileService = &Core{}

const userProfileDisplayNameDBTableName = "user_display_name_dt"
const userProfileImageKeyDBTableName = "user_profile_image_key_dt"

func (core *Core) GetUserBaseProfile(
	inputCtx iam.CallInputContext,
	userID iam.UserID,
) (*iam.UserBaseProfileData, error) {
	if inputCtx == nil {
		return nil, errors.ArgMsg("inputCtx", "missing")
	}
	//TODO(exa): ensure that the context user has the privilege

	return core.getUserBaseProfileInsecure(inputCtx, userID)
}

// getUserBaseProfileInsecure is the implementation of GetUserBaseProfile
// but without access-control. This method must be only used behind the
// access control; for the end-point for public-facing APIs,
// use GetUserBaseProfile.
func (core *Core) getUserBaseProfileInsecure(
	inputCtx iam.CallInputContext,
	userID iam.UserID,
) (*iam.UserBaseProfileData, error) {
	var user iam.UserBaseProfileData
	var idNum iam.UserIDNum
	var deletion iam.UserInstanceDeletionInfo
	var displayName *string
	var profileImageURL *string

	err := core.db.
		QueryRow(
			`SELECT ua.id_num, `+
				`CASE WHEN ua._md_ts IS NULL THEN false ELSE true END AS is_deleted, `+
				`udn.display_name, upiu.profile_image_key `+
				`FROM `+userDBTableName+` AS ua `+
				`LEFT JOIN `+userProfileDisplayNameDBTableName+` udn ON udn.user_id = ua.id_num `+
				`AND udn._md_ts IS NULL `+
				`LEFT JOIN `+userProfileImageKeyDBTableName+` upiu ON upiu.user_id = ua.id_num `+
				`AND upiu._md_ts IS NULL `+
				`WHERE ua.id_num = $1`,
			userID).
		Scan(&idNum, &deletion.Deleted, &displayName, &profileImageURL)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			return nil, err
		}
	}

	if deletion.Deleted {
		//TODO: populate revision number
		user.InstanceInfo = &iam.UserInstanceInfo{Deletion: &deletion}
	} else {
		if displayName != nil {
			user.DisplayName = *displayName
		}
		if profileImageURL != nil {
			user.ProfileImageURL = core.BuildUserProfileImageURL(*profileImageURL)
		}
	}

	return &user, nil
}

func (core *Core) GetUserInfoV1(
	inputCtx iam.CallInputContext,
	userID iam.UserID,
) (*iampb.UserInfoData, error) {
	//TODO: access control

	return core.getUserInfoV1Insecure(inputCtx, userID)
}

func (core *Core) getUserInfoV1Insecure(
	inputCtx iam.CallInputContext,
	userID iam.UserID,
) (*iampb.UserInfoData, error) {
	userBaseProfile, err := core.
		getUserBaseProfileInsecure(inputCtx, userID)
	if err != nil {
		panic(err)
	}
	baseProfile := &iampb.UserBaseProfileData{
		DisplayName:     userBaseProfile.DisplayName,
		ProfileImageUrl: userBaseProfile.ProfileImageURL,
	}

	var deactivation *iampb.UserAccountDeactivationData
	if userBaseProfile.IsDeleted() {
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
		getUserContactInformationInsecure(inputCtx, userID)
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
