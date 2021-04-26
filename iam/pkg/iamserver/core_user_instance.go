package iamserver

import (
	"database/sql"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	goqu "github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

// Interface conformance assertion.
var _ iam.UserInstanceServiceInternal = &Core{}

const userDBTableName = "user_dt"

// GetUserInstanceInfo retrieves the state of an user account. It includes
// the existence of the ID, and wether the account has been deleted.
//
// If it's required only to determine the existence of the ID,
// IsUserRefKeyRegistered is generally more efficient.
func (core *Core) GetUserInstanceInfo(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
) (*iam.UserInstanceInfo, error) {
	//TODO: ACCESS CONTROL
	return core.getUserInstanceInfoNoAC(callCtx, userRef)
}

func (core *Core) getUserInstanceInfoNoAC(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
) (*iam.UserInstanceInfo, error) {
	idRegistered := false
	idRegisteredCacheHit := false
	instDeleted := false
	instDeletionCacheHit := false

	// Look up for an user ID in the cache.
	if _, idRegistered = core.registeredUserIDNumCache.Get(userRef); idRegistered {
		// User ID is positively registered.
		idRegisteredCacheHit = true
	}

	// Look up in the cache
	if _, instDeleted = core.deletedUserIDNumCache.Get(userRef); instDeleted {
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
		getUserInstanceStateByIDNum(userRef.IDNum())
	if err != nil {
		return nil, err
	}

	if !idRegisteredCacheHit && idRegistered {
		core.registeredUserIDNumCache.Add(userRef, nil)
	}
	if !instDeletionCacheHit && instDeleted {
		core.deletedUserIDNumCache.Add(userRef, nil)
	}

	if !idRegistered {
		return nil, nil
	}

	var deletion *iam.UserInstanceDeletionInfo
	if instDeleted {
		//TODO: deletion notes. store the notes as the value in the cache
		deletion = &iam.UserInstanceDeletionInfo{Deleted: true}
	}

	//TODO: populate revision number
	return &iam.UserInstanceInfo{
		RevisionNumber: -1,
		Deletion:       deletion,
	}, nil
}

func (core *Core) getUserInstanceStateByIDNum(
	idNum iam.UserIDNum,
) (idRegistered, accountDeleted bool, err error) {
	sqlString, _, _ := goqu.From(userDBTableName).
		Select(
			goqu.Case().
				When(goqu.C("d_ts").IsNull(), false).
				Else(true).
				As("deleted"),
		).
		Where(
			goqu.C("id").Eq(idNum.PrimitiveValue()),
		).
		ToSQL()

	err = core.db.
		QueryRow(sqlString).
		Scan(&accountDeleted)
	if err == sql.ErrNoRows {
		return false, false, nil
	}
	if err != nil {
		return false, false, err
	}

	return true, accountDeleted, nil
}

func (core *Core) CreateUserInstanceInternal(
	callCtx iam.CallContext,
	input iam.UserInstanceCreationInput,
) (refKey iam.UserRefKey, initialState iam.UserInstanceInfo, err error) {
	//TODO: access control

	refKey, err = core.createUserInstanceNoAC(callCtx)

	//TODO: revision number
	return refKey, iam.UserInstanceInfo{RevisionNumber: -1}, err
}

func (core *Core) createUserInstanceNoAC(
	callCtx iam.CallContext,
) (iam.UserRefKey, error) {
	ctxAuth := callCtx.Authorization()

	const attemptNumMax = 5

	var err error
	var newUserIDNum iam.UserIDNum
	cTime := callCtx.RequestInfo().ReceiveTime

	for attemptNum := 0; ; attemptNum++ {
		//TODO: obtain embedded fields from the argument which
		// type is iam.UserInstanceCreationInput .
		newUserIDNum, err = GenerateUserIDNum(0)
		if err != nil {
			panic(err)
		}

		sqlString, _, _ := goqu.
			Insert(userDBTableName).
			Rows(
				goqu.Record{
					"id":    newUserIDNum,
					"c_ts":  cTime,
					"c_uid": ctxAuth.UserIDNumPtr(),
					"c_tid": ctxAuth.TerminalIDNumPtr(),
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

	return iam.NewUserRefKey(newUserIDNum), nil
}

func (core *Core) DeleteUserInstanceInternal(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
	input iam.UserInstanceDeletionInput,
) (instanceMutated bool, currentState iam.UserInstanceInfo, err error) {
	if callCtx == nil {
		return false, iam.UserInstanceInfoZero(), nil
	}

	ctxAuth := callCtx.Authorization()
	if !ctxAuth.IsUser(userRef) {
		return false, iam.UserInstanceInfoZero(), nil //TODO: should be an error
	}

	return core.deleteUserInstanceNoAC(callCtx, userRef, input)
}

func (core *Core) deleteUserInstanceNoAC(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
	input iam.UserInstanceDeletionInput,
) (instanceMutated bool, currentState iam.UserInstanceInfo, err error) {
	ctxAuth := callCtx.Authorization()
	err = doTx(core.db, func(dbTx *sqlx.Tx) error {
		xres, txErr := dbTx.Exec(
			`UPDATE `+userDBTableName+` `+
				"SET d_ts = $1, d_uid = $2, d_tid = $3, d_notes = $4 "+
				"WHERE id = $2 AND d_ts IS NULL",
			callCtx.RequestInfo().ReceiveTime,
			ctxAuth.UserIDNum().PrimitiveValue(),
			ctxAuth.TerminalIDNum().PrimitiveValue(),
			input.DeletionNotes)
		if txErr != nil {
			return txErr
		}
		n, txErr := xres.RowsAffected()
		if txErr != nil {
			return txErr
		}
		instanceMutated = n == 1

		//TODO: move out. we don't know about key phone number here.
		if txErr == nil {
			_, txErr = dbTx.Exec(
				`UPDATE `+userKeyPhoneNumberDBTableName+` `+
					"SET d_ts = $1, d_uid = $2, d_tid = $3 "+
					"WHERE user_id = $2 AND d_ts IS NULL",
				callCtx.RequestInfo().ReceiveTime,
				ctxAuth.UserIDNum().PrimitiveValue(),
				ctxAuth.TerminalIDNum().PrimitiveValue())
		}

		//TODO: move out. we don't know about profile image here.
		if txErr == nil {
			_, txErr = dbTx.Exec(
				`UPDATE `+userProfileImageKeyDBTableName+` `+
					"SET d_ts = $1, d_uid = $2, d_tid = $3 "+
					"WHERE user_id = $2 AND d_ts IS NULL",
				callCtx.RequestInfo().ReceiveTime,
				ctxAuth.UserIDNum().PrimitiveValue(),
				ctxAuth.TerminalIDNum().PrimitiveValue())
		}

		return txErr
	})
	if err != nil {
		return false, iam.UserInstanceInfoZero(), err
	}

	var deletion *iam.UserInstanceDeletionInfo
	if instanceMutated {
		deletion = &iam.UserInstanceDeletionInfo{
			Deleted:       true,
			DeletionNotes: input.DeletionNotes,
		}
	} else {
		di, err := core.getUserInstanceInfoNoAC(callCtx, userRef)
		if err != nil {
			return false, iam.UserInstanceInfoZero(), err
		}

		if di != nil {
			deletion = di.Deletion
		}
	}

	currentState = iam.UserInstanceInfo{
		RevisionNumber: -1, //TODO: get from the DB
		Deletion:       deletion,
	}

	//TODO: update caches, emit events if there's any changes

	return instanceMutated, currentState, nil
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

	userRef, err = core.createUserInstanceNoAC(callCtx)
	if err != nil {
		return iam.UserRefKeyZero(), false, err
	}

	return userRef, true, nil
}
