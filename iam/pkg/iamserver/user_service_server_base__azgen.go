package iamserver

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"

	errors "github.com/alloyzeus/go-azfl/errors"
	goqu "github.com/doug-martin/goqu/v9"
	lru "github.com/hashicorp/golang-lru"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

const (
	userDBTableName           = "user_dt"
	userDBTablePrimaryKeyName = userDBTableName + "_pkey"
)

// UserServiceServerbase is the server-side
// base implementation of UserService.
type UserServiceServerBase struct {
	db *sqlx.DB

	deletionTxHook func(iam.CallInputContext, *sqlx.Tx) error

	registeredUserIDNumCache *lru.ARCCache
	deletedUserIDNumCache    *lru.ARCCache
}

// Interface conformance assertions.
var (
	_ iam.UserService                 = &UserServiceServerBase{}
	_ iam.UserIDService               = &UserServiceServerBase{}
	_ iam.UserInstanceServiceInternal = &UserServiceServerBase{}
)

func (srv *UserServiceServerBase) IsUserIDRegistered(
	id iam.UserID,
) bool {
	idNum := id.IDNum()

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

// GetUserInstanceInfo retrieves the state of an User
// instance. It includes the existence of the ID, and whether the instance
// has been deleted.
//
// If it's required only to determine the existence of the ID,
// IsUserIDRegistered is generally more efficient.
func (srv *UserServiceServerBase) GetUserInstanceInfo(
	inputCtx iam.CallInputContext,
	id iam.UserID,
) (*iam.UserInstanceInfo, error) {
	//TODO: access control
	return srv.getUserInstanceInfoInsecure(inputCtx, id)
}

func (srv *UserServiceServerBase) getUserInstanceInfoInsecure(
	inputCtx iam.CallInputContext,
	id iam.UserID,
) (*iam.UserInstanceInfo, error) {
	idRegistered := false
	idRegisteredCacheHit := false
	instDeleted := false
	instDeletionCacheHit := false

	// Look up for the ID num in the cache.
	if _, idRegistered = srv.registeredUserIDNumCache.Get(id); idRegistered {
		// ID num is positively registered.
		idRegisteredCacheHit = true
	}

	// Look up in the deletion cache
	if _, instDeleted = srv.deletedUserIDNumCache.Get(id); instDeleted {
		// Instance is positively deleted
		instDeletionCacheHit = true
	}

	if idRegisteredCacheHit && instDeletionCacheHit {
		if !idRegistered {
			return nil, nil
		}
		var deletion *iam.UserInstanceDeletionInfo
		if instDeleted {
			deletion = &iam.UserInstanceDeletionInfo{
				Deleted_: true}
		}
		//TODO: populate revision number
		return &iam.UserInstanceInfo{
			Deletion_: deletion,
		}, nil
	}

	var err error
	idRegistered, instDeleted, err = srv.
		getUserInstanceStateByIDNum(id.IDNum())
	if err != nil {
		return nil, err
	}

	if !idRegisteredCacheHit && idRegistered {
		srv.registeredUserIDNumCache.Add(id, nil)
	}
	if !instDeletionCacheHit && instDeleted {
		srv.deletedUserIDNumCache.Add(id, nil)
	}

	if !idRegistered {
		return nil, nil
	}

	var deletion *iam.UserInstanceDeletionInfo
	if instDeleted {
		//TODO: deletion notes. store the notes as the value in the cache
		deletion = &iam.UserInstanceDeletionInfo{
			Deleted_: true}
	}

	//TODO: populate revision number
	return &iam.UserInstanceInfo{
		RevisionNumber_: -1,
		Deletion_:       deletion,
	}, nil
}

func (srv *UserServiceServerBase) getUserInstanceStateByIDNum(
	idNum iam.UserIDNum,
) (idRegistered, instanceDeleted bool, err error) {
	sqlString, _, _ := goqu.From(userDBTableName).
		Select(
			goqu.Case().
				When(goqu.C(userDBColMetaDeletionTimestamp).IsNull(), false).
				Else(true).
				As("deleted"),
		).
		Where(
			goqu.C(userDBColIDNum).Eq(idNum.PrimitiveValue()),
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
	inputCtx iam.CallInputContext,
	input iam.UserInstanceCreationInput,
) (id iam.UserID, initialState iam.UserInstanceInfo, err error) {
	//TODO: access control

	id, err = srv.createUserInstanceInsecure(inputCtx)

	//TODO: revision number
	return id, iam.UserInstanceInfo{
		RevisionNumber_: -1,
	}, err
}

func (srv *UserServiceServerBase) createUserInstanceInsecure(
	inputCtx iam.CallInputContext,
) (iam.UserID, error) {
	ctxAuth := inputCtx.Authorization()

	const attemptNumMax = 5

	var err error
	var newInstanceIDNum iam.UserIDNum
	cTime := inputCtx.CallInputMetadata().ReceiveTime

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
					userDBColIDNum:                  newInstanceIDNum,
					userDBColMetaCreationTimestamp:  cTime,
					userDBColMetaCreationUserID:     ctxAuth.UserIDNumPtr(),
					userDBColMetaCreationTerminalID: ctxAuth.TerminalIDNumPtr(),
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
			pqErr.Constraint == userDBTablePrimaryKeyName {
			if attemptNum >= attemptNumMax {
				return iam.UserIDZero(), errors.Wrap("insert max attempts", err)
			}
			continue
		}

		return iam.UserIDZero(), errors.Wrap("insert", err)
	}

	//TODO: update caches, emit an event

	return iam.NewUserID(newInstanceIDNum), nil
}

func (srv *UserServiceServerBase) DeleteUserInstanceInternal(
	inputCtx iam.CallInputContext,
	toDelete iam.UserID,
	input iam.UserInstanceDeletionInput,
) (instanceMutated bool, currentState iam.UserInstanceInfo, err error) {
	if inputCtx == nil {
		return false, iam.UserInstanceInfoZero(), nil
	}
	ctxAuth := inputCtx.Authorization()
	if !ctxAuth.IsUser(toDelete) {
		return false, iam.UserInstanceInfoZero(), nil //TODO: should be an error
	}

	//TODO: access control

	return srv.deleteUserInstanceInsecure(inputCtx, toDelete, input)
}

func (srv *UserServiceServerBase) deleteUserInstanceInsecure(
	inputCtx iam.CallInputContext,
	toDelete iam.UserID,
	input iam.UserInstanceDeletionInput,
) (instanceMutated bool, currentState iam.UserInstanceInfo, err error) {
	ctxAuth := inputCtx.Authorization()
	ctxTime := inputCtx.CallInputMetadata().ReceiveTime

	err = doTx(srv.db, func(dbTx *sqlx.Tx) error {
		sqlString, _, _ := goqu.
			From(userDBTableName).
			Where(
				goqu.C(userDBColIDNum).Eq(ctxAuth.UserIDNum().PrimitiveValue()),
				goqu.C(userDBColMetaDeletionTimestamp).IsNull(),
			).
			Update().
			Set(
				goqu.Record{
					userDBColMetaDeletionTimestamp:  ctxTime,
					userDBColMetaDeletionTerminalID: ctxAuth.TerminalIDNum().PrimitiveValue(),
					userDBColMetaDeletionUserID:     ctxAuth.UserIDNum().PrimitiveValue(),
					userDBColMetaDeletionNotes:      input.DeletionNotes,
				},
			).
			ToSQL()

		xres, txErr := dbTx.
			Exec(sqlString)
		if txErr != nil {
			return txErr
		}
		n, txErr := xres.RowsAffected()
		if txErr != nil {
			return txErr
		}
		instanceMutated = n == 1

		if srv.deletionTxHook != nil {
			return srv.deletionTxHook(inputCtx, dbTx)
		}

		return nil
	})
	if err != nil {
		return false, iam.UserInstanceInfoZero(), err
	}

	var deletion *iam.UserInstanceDeletionInfo
	if instanceMutated {
		deletion = &iam.UserInstanceDeletionInfo{
			Deleted_:       true,
			DeletionNotes_: input.DeletionNotes}
	} else {
		di, err := srv.getUserInstanceInfoInsecure(inputCtx, toDelete)
		if err != nil {
			return false, iam.UserInstanceInfoZero(), err
		}

		if di != nil {
			deletion = di.Deletion()
		}
	}

	currentState = iam.UserInstanceInfo{
		RevisionNumber_: -1, //TODO: get from the DB
		Deletion_:       deletion}

	//TODO: update caches, emit an event if there's any changes

	return instanceMutated, currentState, nil
}

// Designed to perform the migration if required
//TODO: context: target version, current version (assert), prefix, etc.
func (srv *UserServiceServerBase) initDataStoreInTx(dbTx *sqlx.Tx) error {
	_, err := dbTx.Exec(
		`CREATE TABLE ` + userDBTableName + ` ( ` +
			userDBColIDNum + `  bigint, ` +
			userDBColMetaCreationTimestamp + `  timestamp with time zone NOT NULL DEFAULT now(), ` +
			userDBColMetaCreationTerminalID + `  bigint, ` +
			userDBColMetaCreationUserID + `  bigint, ` +
			userDBColMetaDeletionTimestamp + `  timestamp with time zone, ` +
			userDBColMetaDeletionTerminalID + `  bigint, ` +
			userDBColMetaDeletionUserID + `  bigint, ` +
			userDBColMetaDeletionNotes + `  jsonb, ` +
			`CONSTRAINT ` + userDBTablePrimaryKeyName + ` PRIMARY KEY(` + userDBColIDNum + `), ` +
			`CHECK (` + userDBColIDNum + ` > 0) ` +
			`);`,
	)
	if err != nil {
		return err
	}
	return nil
}

const (
	userDBColMetaCreationTimestamp  = "_mc_ts"
	userDBColMetaCreationTerminalID = "_mc_tid"
	userDBColMetaCreationUserID     = "_mc_uid"
	userDBColMetaDeletionTimestamp  = "_md_ts"
	userDBColMetaDeletionTerminalID = "_md_tid"
	userDBColMetaDeletionUserID     = "_md_uid"
	userDBColMetaDeletionNotes      = "_md_notes"
	userDBColIDNum                  = "id_num"
)

// GenerateUserIDNum generates a new iam.UserIDNum.
// Note that this function does not consult any database nor registry.
// This method will not create an instance of iam.User, i.e., the
// resulting iam.UserIDNum might or might not refer to valid instance
// of iam.User. The resulting iam.UserIDNum is designed to be
// used as an argument to create a new instance of iam.User.
//
// The embeddedFieldBits argument could be constructed by combining
// iam.UserIDNum*Bits constants. If none are defined,
// use the value of 0.
func GenerateUserIDNum(
	embeddedFieldBits uint64,
) (iam.UserIDNum, error) {
	idBytes := make([]byte, 8)
	_, err := rand.Read(idBytes)
	if err != nil {
		return iam.UserIDNumZero, errors.Wrap("random number source reading", err)
	}

	idUint := (embeddedFieldBits & iam.UserIDNumEmbeddedFieldsMask) |
		(binary.BigEndian.Uint64(idBytes) & iam.UserIDNumIdentifierBitsMask)
	return iam.UserIDNum(idUint), nil
}
