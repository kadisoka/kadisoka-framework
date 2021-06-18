package iamserver

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"

	errors "github.com/alloyzeus/go-azfl/azfl/errors"
	goqu "github.com/doug-martin/goqu/v9"
	lru "github.com/hashicorp/golang-lru"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

const userDBTableName = "user_dt"

// UserServiceServerbase is the server-side
// base implementation of UserService.
type UserServiceServerBase struct {
	db *sqlx.DB

	deletionTxHook func(iam.OpInputContext, *sqlx.Tx) error

	registeredUserIDNumCache *lru.ARCCache
	deletedUserIDNumCache    *lru.ARCCache
}

// Interface conformance assertions.
var (
	_ iam.UserService                 = &UserServiceServerBase{}
	_ iam.UserRefKeyService           = &UserServiceServerBase{}
	_ iam.UserInstanceServiceInternal = &UserServiceServerBase{}
)

func (srv *UserServiceServerBase) IsUserRefKeyRegistered(refKey iam.UserRefKey) bool {
	idNum := refKey.IDNum()

	// Look up for the ID num in the cache.
	if _, idRegistered := srv.registeredUserIDNumCache.Get(idNum); idRegistered {
		return true
	}

	idRegistered, _, err := srv.
		getUserInstanceStateByIDNum(idNum)
	if err != nil {
		panic(err)
	}

	if idRegistered {
		srv.registeredUserIDNumCache.Add(idNum, nil)
	}

	return idRegistered
}

// GetUserInstanceInfo retrieves the state of an User instance.
// It includes the existence of the ID, and whether the instance
// has been deleted.
//
// If it's required only to determine the existence of the ID,
// IsUserRefKeyRegistered is generally more efficient.
func (srv *UserServiceServerBase) GetUserInstanceInfo(
	callCtx iam.OpInputContext,
	refKey iam.UserRefKey,
) (*iam.UserInstanceInfo, error) {
	//TODO: access control
	return srv.getUserInstanceInfoNoAC(callCtx, refKey)
}

func (srv *UserServiceServerBase) getUserInstanceInfoNoAC(
	callCtx iam.OpInputContext,
	refKey iam.UserRefKey,
) (*iam.UserInstanceInfo, error) {
	idRegistered := false
	idRegisteredCacheHit := false
	instDeleted := false
	instDeletionCacheHit := false

	// Look up for the ID num in the cache.
	if _, idRegistered = srv.registeredUserIDNumCache.Get(refKey); idRegistered {
		// ID num is positively registered.
		idRegisteredCacheHit = true
	}

	// Look up in the deletion cache
	if _, instDeleted = srv.deletedUserIDNumCache.Get(refKey); instDeleted {
		// Instance is positively deleted
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
	idRegistered, instDeleted, err = srv.
		getUserInstanceStateByIDNum(refKey.IDNum())
	if err != nil {
		return nil, err
	}

	if !idRegisteredCacheHit && idRegistered {
		srv.registeredUserIDNumCache.Add(refKey, nil)
	}
	if !instDeletionCacheHit && instDeleted {
		srv.deletedUserIDNumCache.Add(refKey, nil)
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

func (srv *UserServiceServerBase) getUserInstanceStateByIDNum(
	idNum iam.UserIDNum,
) (idRegistered, instanceDeleted bool, err error) {
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

	err = srv.db.
		QueryRow(sqlString).
		Scan(&instanceDeleted)
	if err == sql.ErrNoRows {
		return false, false, nil
	}
	if err != nil {
		return false, false, err
	}

	return true, instanceDeleted, nil
}

func (srv *UserServiceServerBase) CreateUserInstanceInternal(
	callCtx iam.OpInputContext,
	input iam.UserInstanceCreationInput,
) (refKey iam.UserRefKey, initialState iam.UserInstanceInfo, err error) {
	//TODO: access control

	refKey, err = srv.createUserInstanceNoAC(callCtx)

	//TODO: revision number
	return refKey, iam.UserInstanceInfo{RevisionNumber: -1}, err
}

func (srv *UserServiceServerBase) createUserInstanceNoAC(
	callCtx iam.OpInputContext,
) (iam.UserRefKey, error) {
	ctxAuth := callCtx.Authorization()

	const attemptNumMax = 5

	var err error
	var newInstanceIDNum iam.UserIDNum
	cTime := callCtx.OpInputMetadata().ReceiveTime

	for attemptNum := 0; ; attemptNum++ {
		//TODO: obtain embedded fields from the argument which
		// type is iam.UserInstanceCreationInput .
		newInstanceIDNum, err = GenerateUserIDNum(0)
		if err != nil {
			panic(err)
		}

		sqlString, _, _ := goqu.
			Insert(userDBTableName).
			Rows(
				goqu.Record{
					"id":    newInstanceIDNum,
					"c_ts":  cTime,
					"c_uid": ctxAuth.UserIDNumPtr(),
					"c_tid": ctxAuth.TerminalIDNumPtr(),
				},
			).
			ToSQL()

		_, err = srv.db.
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

	//TODO: update caches, emit an event

	return iam.NewUserRefKey(newInstanceIDNum), nil
}

func (srv *UserServiceServerBase) DeleteUserInstanceInternal(
	callCtx iam.OpInputContext,
	toDelete iam.UserRefKey,
	input iam.UserInstanceDeletionInput,
) (instanceMutated bool, currentState iam.UserInstanceInfo, err error) {
	if callCtx == nil {
		return false, iam.UserInstanceInfoZero(), nil
	}
	ctxAuth := callCtx.Authorization()
	if !ctxAuth.IsUser(toDelete) {
		return false, iam.UserInstanceInfoZero(), nil //TODO: should be an error
	}

	//TODO: access control

	return srv.deleteUserInstanceNoAC(callCtx, toDelete, input)
}

func (srv *UserServiceServerBase) deleteUserInstanceNoAC(
	callCtx iam.OpInputContext,
	toDelete iam.UserRefKey,
	input iam.UserInstanceDeletionInput,
) (instanceMutated bool, currentState iam.UserInstanceInfo, err error) {
	ctxAuth := callCtx.Authorization()
	err = doTx(srv.db, func(dbTx *sqlx.Tx) error {
		xres, txErr := dbTx.Exec(
			`UPDATE `+userDBTableName+` `+
				"SET d_ts = $1, d_uid = $2, d_tid = $3, d_notes = $4 "+
				"WHERE id = $2 AND d_ts IS NULL",
			callCtx.OpInputMetadata().ReceiveTime,
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

		if srv.deletionTxHook != nil {
			return srv.deletionTxHook(callCtx, dbTx)
		}

		return nil
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
		di, err := srv.getUserInstanceInfoNoAC(callCtx, toDelete)
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

	//TODO: update caches, emit an event if there's any changes

	return instanceMutated, currentState, nil
}

// GenerateUserIDNum generates a new iam.UserIDNum.
// Note that this function does not consulting any database nor registry.
// This method will not create an instance of iam.User, i.e., the
// resulting iam.UserIDNum might or might not refer to valid instance
// of iam.User. The resulting iam.UserIDNum is designed to be
// used as an argument to create a new instance of iam.User.
//
// The embeddedFieldBits argument could be constructed by combining
// iam.UserIDNum*Bits constants. If none are defined,
// use the value of 0.
func GenerateUserIDNum(embeddedFieldBits uint64) (iam.UserIDNum, error) {
	idBytes := make([]byte, 8)
	_, err := rand.Read(idBytes)
	if err != nil {
		return iam.UserIDNumZero, errors.Wrap("random number source reading", err)
	}

	idUint := (embeddedFieldBits & iam.UserIDNumEmbeddedFieldsMask) |
		(binary.BigEndian.Uint64(idBytes) & iam.UserIDNumIdentifierBitsMask)
	return iam.UserIDNum(idUint), nil
}
