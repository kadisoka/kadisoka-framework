package iamserver

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"time"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	"github.com/jmoiron/sqlx"
	iampb "github.com/rez-go/crux-apis/crux/iam/v1"
	"golang.org/x/crypto/blake2b"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

const userTableName = "user_dt"
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

// GetUserAccountState retrieves the state of an user account. It includes
// the existence of the ID, and wether the account has been deleted.
//
// If it's required only to determine the existence of the ID,
// IsUserIDRegistered is generally more efficient.
func (core *Core) GetUserAccountState(
	userRef iam.UserRefKey,
) (*iam.UserAccountState, error) {
	idRegistered := false
	idRegisteredCacheHit := false
	accountDeleted := false
	accountDeletedCacheHit := false
	// Look up for an user ID in the cache.
	if _, idRegistered = core.registeredUserIDCache.Get(userRef); idRegistered {
		// User ID is positively registered.
		idRegisteredCacheHit = true
	}

	// Look up in the cache
	if _, accountDeleted := core.deletedUserAccountIDCache.Get(userRef); accountDeleted {
		// Account is positively deleted
		accountDeletedCacheHit = true
	}

	if idRegisteredCacheHit && accountDeletedCacheHit {
		if !idRegistered {
			return nil, nil
		}
		return &iam.UserAccountState{
			Deleted: accountDeleted,
		}, nil
	}

	var err error
	idRegistered, accountDeleted, err = core.
		getUserAccountState(userRef.ID())
	if err != nil {
		return nil, err
	}

	if !idRegisteredCacheHit && idRegistered {
		core.registeredUserIDCache.Add(userRef, nil)
	}
	if !accountDeletedCacheHit && accountDeleted {
		core.deletedUserAccountIDCache.Add(userRef, nil)
	}

	if !idRegistered {
		return nil, nil
	}
	return &iam.UserAccountState{
		Deleted: accountDeleted,
	}, nil
}

// IsUserIDRegistered is used to determine that a user ID has been registered.
// It's not checking if the account is active or not.
//
// This function is generally cheap if the user ID has been registered.
func (core *Core) IsUserIDRegistered(id iam.UserID) bool {
	// Look up for an user ID in the cache.
	if _, idRegistered := core.registeredUserIDCache.Get(id); idRegistered {
		return true
	}

	idRegistered, _, err := core.
		getUserAccountState(id)
	if err != nil {
		panic(err)
	}

	if idRegistered {
		core.registeredUserIDCache.Add(id, nil)
	}

	return idRegistered
}

func (core *Core) getUserAccountState(
	id iam.UserID,
) (idRegistered, accountDeleted bool, err error) {
	err = core.db.
		QueryRow(
			`SELECT CASE WHEN d_ts IS NULL THEN false ELSE true END `+
				`FROM `+userTableName+` WHERE id = $1`,
			id).
		Scan(&accountDeleted)
	if err == sql.ErrNoRows {
		return false, false, nil
	}
	if err != nil {
		return false, false, err
	}

	return true, accountDeleted, nil
}

func (core *Core) DeleteUserAccount(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
	input iam.UserAccountDeleteInput,
) (deleted bool, err error) {
	if callCtx == nil {
		return false, nil
	}
	authCtx := callCtx.Authorization()
	if !authCtx.IsUserContext() || !userRef.EqualsUserRefKey(authCtx.UserRef()) {
		return false, nil
	}

	err = doTx(core.db, func(dbTx *sqlx.Tx) error {
		xres, txErr := dbTx.Exec(
			`UPDATE `+userTableName+` `+
				"SET d_ts = $1, d_uid = $2, d_tid = $3, d_notes = $4 "+
				"WHERE id = $2 AND d_ts IS NULL",
			callCtx.RequestReceiveTime(),
			authCtx.UserID().PrimitiveValue(),
			authCtx.TerminalID().PrimitiveValue(),
			input.DeletionNotes)
		if txErr != nil {
			return txErr
		}
		n, txErr := xres.RowsAffected()
		if txErr != nil {
			return txErr
		}
		deleted = n == 1

		if txErr == nil {
			_, txErr = dbTx.Exec(
				`UPDATE `+userKeyPhoneNumberTableName+` `+
					"SET d_ts = $1, d_uid = $2, d_tid = $3 "+
					"WHERE user_id = $2 AND d_ts IS NULL",
				callCtx.RequestReceiveTime(),
				authCtx.UserID().PrimitiveValue(),
				authCtx.TerminalID().PrimitiveValue())
		}

		if txErr == nil {
			_, txErr = dbTx.Exec(
				`UPDATE `+userProfileImageKeyTableName+` `+
					"SET d_ts = $1, d_uid = $2, d_tid = $3 "+
					"WHERE user_id = $2 AND d_ts IS NULL",
				callCtx.RequestReceiveTime(),
				authCtx.UserID().PrimitiveValue(),
				authCtx.TerminalID().PrimitiveValue())
		}

		return txErr
	})
	if err != nil {
		return false, err
	}

	//TODO: update caches, emit events if there's any changes

	return deleted, nil
}

func (core *Core) SetUserProfileImageURL(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
	profileImageURL string,
) error {
	authCtx := callCtx.Authorization()
	// Change this if we want to allow service client to update a user's profile
	// (we'll need a better access control for service clients)
	if !authCtx.IsUserContext() {
		return iam.ErrUserContextRequired
	}
	// Don't allow changing other user's for now
	if !userRef.EqualsUserRefKey(authCtx.UserRef()) {
		return iam.ErrContextUserNotAllowedToPerformActionOnResource
	}
	if profileImageURL != "" && !core.isUserProfileImageURLAllowed(profileImageURL) {
		return errors.ArgMsg("profileImageURL", "unsupported")
	}

	//TODO: on changes, update caches, emit events only if there's any changes

	return doTx(core.db, func(dbTx *sqlx.Tx) error {
		_, txErr := dbTx.Exec(
			`UPDATE `+userKeyPhoneNumberTableName+` `+
				"SET d_ts = $1, d_uid = $2, d_tid = $3 "+
				"WHERE user_id = $2 AND d_ts IS NULL",
			callCtx.RequestReceiveTime(),
			authCtx.UserID().PrimitiveValue(),
			authCtx.TerminalID().PrimitiveValue())
		if txErr != nil {
			return errors.Wrap("mark current profile image URL as deleted", txErr)
		}
		if profileImageURL != "" {
			_, txErr = dbTx.Exec(
				`INSERT INTO `+userProfileImageKeyTableName+` `+
					"(user_id, profile_image_key, c_uid, c_tid) VALUES "+
					"($1, $2, $3, $4)",
				authCtx.UserID().PrimitiveValue(), profileImageURL,
				authCtx.UserID().PrimitiveValue(), authCtx.TerminalID().PrimitiveValue())
			if txErr != nil {
				return errors.Wrap("insert new profile image URL", txErr)
			}
		}
		return nil
	})
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

func (core *Core) GetUserContactInformation(
	callCtx iam.CallContext,
	userID iam.UserRefKey,
) (*iampb.UserContactInfoData, error) {
	//TODO: access control
	userPhoneNumber, err := core.
		GetUserKeyPhoneNumber(callCtx, userID)
	if err != nil {
		return nil, errors.Wrap("get user key phone number", err)
	}
	if userPhoneNumber == nil {
		return nil, nil
	}
	return &iampb.UserContactInfoData{
		PhoneNumber: userPhoneNumber.String(),
	}, nil
}

func (core *Core) isUserProfileImageURLAllowed(profileImageURL string) bool {
	//TODO(exa): limit profile image url to certain hosts or keep only the filename
	return true
}

func (core *Core) ensureOrNewUserRef(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
) (iam.UserRefKey, error) {
	if callCtx == nil {
		return iam.UserRefKeyZero(), errors.ArgMsg("callCtx", "missing")
	}
	if userRef.IsValid() {
		if !core.IsUserIDRegistered(userRef.ID()) {
			return iam.UserRefKeyZero(), nil
		}
		return userRef, nil
	}

	var err error
	userRef, err = core.CreateUserAccount(callCtx)
	if err != nil {
		return iam.UserRefKeyZero(), err
	}

	return userRef, nil
}

func (core *Core) CreateUserAccount(
	callCtx iam.CallContext,
) (iam.UserRefKey, error) {
	newUserID, err := core.generateUserID()
	if err != nil {
		panic(err)
	}

	//TODO: if id conflict, generate another id and retry
	_, err = core.db.
		Exec(
			`INSERT INTO `+userTableName+` (`+
				`id, c_ts, c_uid, c_tid`+
				`) VALUES (`+
				`$1, $2, $3, $4`+
				`)`,
			newUserID,
			callCtx.RequestReceiveTime(),
			callCtx.Authorization().UserIDPtr(),
			callCtx.Authorization().TerminalIDPtr())
	if err != nil {
		return iam.UserRefKeyZero(), err
	}

	return iam.NewUserRefKey(newUserID), nil
}

func (core *Core) generateUserID() (iam.UserID, error) {
	var userID iam.UserID
	var err error
	for i := 0; i < 5; i++ {
		userID, err = core.generateUserIDImpl()
		if err == nil && userID.IsValid() {
			return userID, nil
		}
	}
	if err == nil {
		err = errors.Msg("user ID generation failed")
	}
	return iam.UserIDZero, err
}

func (core *Core) generateUserIDImpl() (iam.UserID, error) {
	tNow := time.Now().UTC()
	tbin, err := tNow.MarshalBinary()
	if err != nil {
		panic(err)
	}
	hasher, err := blake2b.New(4, nil)
	if err != nil {
		panic(err)
	}
	hasher.Write(tbin)
	hashPart := hasher.Sum(nil)
	idBytes := make([]byte, 8)
	_, err = rand.Read(idBytes[2:4])
	if err != nil {
		panic(err)
	}
	copy(idBytes[4:], hashPart)
	idUint := binary.BigEndian.Uint64(idBytes) & iam.UserIDSignificantBitsMask
	return iam.UserID(idUint), nil
}
