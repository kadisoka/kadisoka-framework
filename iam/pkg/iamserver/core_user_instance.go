package iamserver

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"time"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	goqu "github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/blake2b"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

// Interface conformance assertion.
var _ iam.UserInstanceService = &Core{}

const userDBTableName = "user_dt"

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
				`FROM `+userDBTableName+` WHERE id = $1`,
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
) (stateChanged bool, err error) {
	if callCtx == nil {
		return false, nil
	}
	ctxAuth := callCtx.Authorization()
	if !ctxAuth.IsUser(userRef) {
		return false, nil //TODO: should be an error
	}

	return core.deleteUserInstanceNoAC(callCtx, userRef, input)
}

func (core *Core) deleteUserInstanceNoAC(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
	input iam.UserInstanceDeletionInput,
) (stateChanged bool, err error) {
	ctxAuth := callCtx.Authorization()
	err = doTx(core.db, func(dbTx *sqlx.Tx) error {
		xres, txErr := dbTx.Exec(
			`UPDATE `+userDBTableName+` `+
				"SET d_ts = $1, d_uid = $2, d_tid = $3, d_notes = $4 "+
				"WHERE id = $2 AND d_ts IS NULL",
			callCtx.RequestInfo().ReceiveTime,
			ctxAuth.UserID().PrimitiveValue(),
			ctxAuth.TerminalID().PrimitiveValue(),
			input.DeletionNotes)
		if txErr != nil {
			return txErr
		}
		n, txErr := xres.RowsAffected()
		if txErr != nil {
			return txErr
		}
		stateChanged = n == 1

		if txErr == nil {
			_, txErr = dbTx.Exec(
				`UPDATE `+userKeyPhoneNumberDBTableName+` `+
					"SET d_ts = $1, d_uid = $2, d_tid = $3 "+
					"WHERE user_id = $2 AND d_ts IS NULL",
				callCtx.RequestInfo().ReceiveTime,
				ctxAuth.UserID().PrimitiveValue(),
				ctxAuth.TerminalID().PrimitiveValue())
		}

		if txErr == nil {
			_, txErr = dbTx.Exec(
				`UPDATE `+userProfileImageKeyDBTableName+` `+
					"SET d_ts = $1, d_uid = $2, d_tid = $3 "+
					"WHERE user_id = $2 AND d_ts IS NULL",
				callCtx.RequestInfo().ReceiveTime,
				ctxAuth.UserID().PrimitiveValue(),
				ctxAuth.TerminalID().PrimitiveValue())
		}

		return txErr
	})
	if err != nil {
		return false, err
	}

	//TODO: update caches, emit events if there's any changes

	return stateChanged, nil
}

func (core *Core) contextUserOrNewInstance(
	callCtx iam.CallContext,
) (userRef iam.UserRefKey, newInstance bool, err error) {
	if callCtx == nil {
		return iam.UserRefKeyZero(), false, errors.ArgMsg("callCtx", "missing")
	}
	ctxAuth := callCtx.Authorization()
	if ctxAuth.IsUserContext() {
		userRef = ctxAuth.UserRef()
		if !core.IsUserRefKeyRegistered(userRef) {
			return iam.UserRefKeyZero(), false, errors.ArgMsg("callCtx.Authorization", "invalid")
		}
		return userRef, false, nil
	}

	userRef, err = core.newUserInstance(callCtx)
	if err != nil {
		return iam.UserRefKeyZero(), false, err
	}

	return userRef, true, nil
}

func (core *Core) newUserInstance(
	callCtx iam.CallContext,
) (iam.UserRefKey, error) {
	ctxAuth := callCtx.Authorization()

	const attemptNumMax = 5

	var err error
	var newUserID iam.UserID
	cTime := callCtx.RequestInfo().ReceiveTime

	for attemptNum := 0; ; attemptNum++ {
		newUserID, err = core.generateUserID()
		if err != nil {
			panic(err)
		}

		sqlString, _, _ := goqu.
			Insert(userDBTableName).
			Rows(
				goqu.Record{
					"id":    newUserID,
					"c_ts":  cTime,
					"c_uid": ctxAuth.UserIDPtr(),
					"c_tid": ctxAuth.TerminalIDPtr(),
				},
			).
			ToSQL()

		_, err = core.db.
			Exec(sqlString)
		if err == nil {
			break
		}

		pqErr, _ := err.(*pq.Error)
		if pqErr != nil &&
			pqErr.Code == "23505" &&
			pqErr.Constraint == userDBTableName+"_pkey" {
			if attemptNum >= attemptNumMax {
				return iam.UserRefKeyZero(), errors.Wrap("insert max attempts", err)
			}
			continue
		}

		return iam.UserRefKeyZero(), errors.Wrap("insert", err)
	}

	return iam.NewUserRefKey(newUserID), nil
}

//TODO: bitfield
func (core *Core) generateUserID() (iam.UserID, error) {
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
