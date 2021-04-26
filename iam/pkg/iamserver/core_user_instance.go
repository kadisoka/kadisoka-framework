package iamserver

import (
	"github.com/alloyzeus/go-azfl/azfl/errors"
	goqu "github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

// Interface conformance assertion.
var _ iam.UserInstanceServiceInternal = &Core{}

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
		di, err := core.UserService.getUserInstanceInfoNoAC(callCtx, userRef)
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
		if !core.UserService.IsUserRefKeyRegistered(userRef) {
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
