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

// Interface conformance assertion.
var _ iam.UserInstanceService = &Core{}

const userTableName = "user_dt"

// GetUserInstanceInfo retrieves the state of an user account. It includes
// the existence of the ID, and wether the account has been deleted.
//
// If it's required only to determine the existence of the ID,
// IsUserRefKeyRegistered is generally more efficient.
func (core *Core) GetUserInstanceInfo(
	userRef iam.UserRefKey,
) (*iam.UserInstanceInfo, error) {
	idRegistered := false
	idRegisteredCacheHit := false
	instDeleted := false
	instDeletionCacheHit := false

	// Look up for an user ID in the cache.
	if _, idRegistered = core.registeredUserInstanceIDCache.Get(userRef); idRegistered {
		// User ID is positively registered.
		idRegisteredCacheHit = true
	}

	// Look up in the cache
	if _, instDeleted = core.deletedUserInstanceIDCache.Get(userRef); instDeleted {
		// Account is positively deleted
		instDeletionCacheHit = true
	}

	if idRegisteredCacheHit && instDeletionCacheHit {
		if !idRegistered {
			return nil, nil
		}
		var deletion *iam.UserInstanceDeletionInfo
		if instDeleted {
			deletion = &iam.UserInstanceDeletionInfo{Deleted: true}
		}
		//TODO: populate revision number
		return &iam.UserInstanceInfo{
			Deletion: deletion,
		}, nil
	}

	var err error
	idRegistered, instDeleted, err = core.
		getUserInstanceStateByID(userRef.ID())
	if err != nil {
		return nil, err
	}

	if !idRegisteredCacheHit && idRegistered {
		core.registeredUserInstanceIDCache.Add(userRef, nil)
	}
	if !instDeletionCacheHit && instDeleted {
		core.deletedUserInstanceIDCache.Add(userRef, nil)
	}

	if !idRegistered {
		return nil, nil
	}

	var deletion *iam.UserInstanceDeletionInfo
	if instDeleted {
		deletion = &iam.UserInstanceDeletionInfo{Deleted: true}
	}
	//TODO: populate revision number
	return &iam.UserInstanceInfo{
		Deletion: deletion,
	}, nil
}

func (core *Core) getUserInstanceStateByID(
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

func (core *Core) DeleteUserInstance(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
	input iam.UserInstanceDeletionInput,
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
			callCtx.RequestInfo().ReceiveTime,
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
				callCtx.RequestInfo().ReceiveTime,
				authCtx.UserID().PrimitiveValue(),
				authCtx.TerminalID().PrimitiveValue())
		}

		if txErr == nil {
			_, txErr = dbTx.Exec(
				`UPDATE `+userProfileImageKeyTableName+` `+
					"SET d_ts = $1, d_uid = $2, d_tid = $3 "+
					"WHERE user_id = $2 AND d_ts IS NULL",
				callCtx.RequestInfo().ReceiveTime,
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

func (core *Core) ensureOrNewUserRef(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
) (iam.UserRefKey, error) {
	if callCtx == nil {
		return iam.UserRefKeyZero(), errors.ArgMsg("callCtx", "missing")
	}
	if userRef.IsValid() {
		if !core.IsUserRefKeyRegistered(userRef) {
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
			callCtx.RequestInfo().ReceiveTime,
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
