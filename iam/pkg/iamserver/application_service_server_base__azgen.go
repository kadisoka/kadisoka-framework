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

const applicationDBTableName = "application_dt"

// ApplicationServiceServerbase is the server-side
// base implementation of ApplicationService.
type ApplicationServiceServerBase struct {
	db *sqlx.DB

	deletionTxHook func(iam.OpInputContext, *sqlx.Tx) error

	registeredApplicationIDNumCache *lru.ARCCache
	deletedApplicationIDNumCache    *lru.ARCCache
}

// Interface conformance assertions.
var (
	_ iam.ApplicationService                 = &ApplicationServiceServerBase{}
	_ iam.ApplicationRefKeyService           = &ApplicationServiceServerBase{}
	_ iam.ApplicationInstanceServiceInternal = &ApplicationServiceServerBase{}
)

func (srv *ApplicationServiceServerBase) IsApplicationRefKeyRegistered(refKey iam.ApplicationRefKey) bool {
	idNum := refKey.IDNum()

	// Look up for the ID num in the cache.
	if _, idRegistered := srv.registeredApplicationIDNumCache.Get(idNum); idRegistered {
		return true
	}

	idRegistered, _, err := srv.
		getApplicationInstanceStateByIDNum(idNum)
	if err != nil {
		panic(err)
	}

	if idRegistered {
		srv.registeredApplicationIDNumCache.Add(idNum, nil)
	}

	return idRegistered
}

// GetApplicationInstanceInfo retrieves the state of an Application instance.
// It includes the existence of the ID, and whether the instance
// has been deleted.
//
// If it's required only to determine the existence of the ID,
// IsApplicationRefKeyRegistered is generally more efficient.
func (srv *ApplicationServiceServerBase) GetApplicationInstanceInfo(
	opInputCtx iam.OpInputContext,
	refKey iam.ApplicationRefKey,
) (*iam.ApplicationInstanceInfo, error) {
	//TODO: access control
	return srv.getApplicationInstanceInfoNoAC(opInputCtx, refKey)
}

func (srv *ApplicationServiceServerBase) getApplicationInstanceInfoNoAC(
	opInputCtx iam.OpInputContext,
	refKey iam.ApplicationRefKey,
) (*iam.ApplicationInstanceInfo, error) {
	idRegistered := false
	idRegisteredCacheHit := false
	instDeleted := false
	instDeletionCacheHit := false

	// Look up for the ID num in the cache.
	if _, idRegistered = srv.registeredApplicationIDNumCache.Get(refKey); idRegistered {
		// ID num is positively registered.
		idRegisteredCacheHit = true
	}

	// Look up in the deletion cache
	if _, instDeleted = srv.deletedApplicationIDNumCache.Get(refKey); instDeleted {
		// Instance is positively deleted
		instDeletionCacheHit = true
	}

	if idRegisteredCacheHit && instDeletionCacheHit {
		if !idRegistered {
			return nil, nil
		}
		var deletion *iam.ApplicationInstanceDeletionInfo
		if instDeleted {
			deletion = &iam.ApplicationInstanceDeletionInfo{Deleted: true}
		}
		//TODO: populate revision number
		return &iam.ApplicationInstanceInfo{
			Deletion: deletion,
		}, nil
	}

	var err error
	idRegistered, instDeleted, err = srv.
		getApplicationInstanceStateByIDNum(refKey.IDNum())
	if err != nil {
		return nil, err
	}

	if !idRegisteredCacheHit && idRegistered {
		srv.registeredApplicationIDNumCache.Add(refKey, nil)
	}
	if !instDeletionCacheHit && instDeleted {
		srv.deletedApplicationIDNumCache.Add(refKey, nil)
	}

	if !idRegistered {
		return nil, nil
	}

	var deletion *iam.ApplicationInstanceDeletionInfo
	if instDeleted {
		//TODO: deletion notes. store the notes as the value in the cache
		deletion = &iam.ApplicationInstanceDeletionInfo{Deleted: true}
	}

	//TODO: populate revision number
	return &iam.ApplicationInstanceInfo{
		RevisionNumber: -1,
		Deletion:       deletion,
	}, nil
}

func (srv *ApplicationServiceServerBase) getApplicationInstanceStateByIDNum(
	idNum iam.ApplicationIDNum,
) (idRegistered, instanceDeleted bool, err error) {
	sqlString, _, _ := goqu.From(applicationDBTableName).
		Select(
			goqu.Case().
				When(goqu.C("_md_ts").IsNull(), false).
				Else(true).
				As("deleted"),
		).
		Where(
			goqu.C("id_num").Eq(idNum.PrimitiveValue()),
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

func (srv *ApplicationServiceServerBase) CreateApplicationInstanceInternal(
	opInputCtx iam.OpInputContext,
	input iam.ApplicationInstanceCreationInput,
) (refKey iam.ApplicationRefKey, initialState iam.ApplicationInstanceInfo, err error) {
	//TODO: access control

	refKey, err = srv.createApplicationInstanceNoAC(opInputCtx)

	//TODO: revision number
	return refKey, iam.ApplicationInstanceInfo{RevisionNumber: -1}, err
}

func (srv *ApplicationServiceServerBase) createApplicationInstanceNoAC(
	opInputCtx iam.OpInputContext,
) (iam.ApplicationRefKey, error) {
	ctxAuth := opInputCtx.Authorization()

	const attemptNumMax = 5

	var err error
	var newInstanceIDNum iam.ApplicationIDNum
	cTime := opInputCtx.OpInputMetadata().ReceiveTime

	for attemptNum := 0; ; attemptNum++ {
		//TODO: obtain embedded fields from the argument which
		// type is iam.ApplicationInstanceCreationInput .
		newInstanceIDNum, err = GenerateApplicationIDNum(0)
		if err != nil {
			panic(err)
		}

		sqlString, _, _ := goqu.
			Insert(applicationDBTableName).
			Rows(
				goqu.Record{
					"id_num":  newInstanceIDNum,
					"_mc_ts":  cTime,
					"_mc_uid": ctxAuth.UserIDNumPtr(),
					"_mc_tid": ctxAuth.TerminalIDNumPtr(),
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
			pqErr.Constraint == applicationDBTableName+"_pkey" {
			if attemptNum >= attemptNumMax {
				return iam.ApplicationRefKeyZero(), errors.Wrap("insert max attempts", err)
			}
			continue
		}

		return iam.ApplicationRefKeyZero(), errors.Wrap("insert", err)
	}

	//TODO: update caches, emit an event

	return iam.NewApplicationRefKey(newInstanceIDNum), nil
}

func (srv *ApplicationServiceServerBase) DeleteApplicationInstanceInternal(
	opInputCtx iam.OpInputContext,
	toDelete iam.ApplicationRefKey,
	input iam.ApplicationInstanceDeletionInput,
) (instanceMutated bool, currentState iam.ApplicationInstanceInfo, err error) {
	if opInputCtx == nil {
		return false, iam.ApplicationInstanceInfoZero(), nil
	}

	//TODO: access control

	return srv.deleteApplicationInstanceNoAC(opInputCtx, toDelete, input)
}

func (srv *ApplicationServiceServerBase) deleteApplicationInstanceNoAC(
	opInputCtx iam.OpInputContext,
	toDelete iam.ApplicationRefKey,
	input iam.ApplicationInstanceDeletionInput,
) (instanceMutated bool, currentState iam.ApplicationInstanceInfo, err error) {
	ctxAuth := opInputCtx.Authorization()
	ctxTime := opInputCtx.OpInputMetadata().ReceiveTime

	err = doTx(srv.db, func(dbTx *sqlx.Tx) error {
		sqlString, _, _ := goqu.
			From(applicationDBTableName).
			Where(
				goqu.C("id_num").Eq(ctxAuth.UserIDNum().PrimitiveValue()),
				goqu.C("_md_ts").IsNull(),
			).
			Update().
			Set(
				goqu.Record{
					"_md_ts":  ctxTime,
					"_md_tid": ctxAuth.TerminalIDNum().PrimitiveValue(),
					"_md_uid": ctxAuth.UserIDNum().PrimitiveValue(),
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
			return srv.deletionTxHook(opInputCtx, dbTx)
		}

		return nil
	})
	if err != nil {
		return false, iam.ApplicationInstanceInfoZero(), err
	}

	var deletion *iam.ApplicationInstanceDeletionInfo
	if instanceMutated {
		deletion = &iam.ApplicationInstanceDeletionInfo{
			Deleted: true,
		}
	} else {
		di, err := srv.getApplicationInstanceInfoNoAC(opInputCtx, toDelete)
		if err != nil {
			return false, iam.ApplicationInstanceInfoZero(), err
		}

		if di != nil {
			deletion = di.Deletion
		}
	}

	currentState = iam.ApplicationInstanceInfo{
		RevisionNumber: -1, //TODO: get from the DB
		Deletion:       deletion,
	}

	//TODO: update caches, emit an event if there's any changes

	return instanceMutated, currentState, nil
}

// Designed to perform the migration if required
//TODO: context: target version, current version (assert), prefix, etc.
func (srv *ApplicationServiceServerBase) initDataStoreInTx(dbTx *sqlx.Tx) error {
	_, err := dbTx.Exec(
		`CREATE TABLE ` + applicationDBTableName + ` ( ` +
			`id_num     integer PRIMARY KEY, ` +
			`_mc_ts     timestamp with time zone NOT NULL DEFAULT now(), ` +
			`_mc_tid    bigint, ` +
			`_mc_uid    bigint, ` +
			`_md_ts     timestamp with time zone, ` +
			`_md_tid    bigint, ` +
			`_md_uid    bigint, ` +
			`CHECK (id_num > 0) ` +
			`);`,
	)
	if err != nil {
		return err
	}
	return nil
}

// GenerateApplicationIDNum generates a new iam.ApplicationIDNum.
// Note that this function does not consult any database nor registry.
// This method will not create an instance of iam.Application, i.e., the
// resulting iam.ApplicationIDNum might or might not refer to valid instance
// of iam.Application. The resulting iam.ApplicationIDNum is designed to be
// used as an argument to create a new instance of iam.Application.
//
// The embeddedFieldBits argument could be constructed by combining
// iam.ApplicationIDNum*Bits constants. If none are defined,
// use the value of 0.
func GenerateApplicationIDNum(embeddedFieldBits uint32) (iam.ApplicationIDNum, error) {
	idBytes := make([]byte, 4)
	_, err := rand.Read(idBytes)
	if err != nil {
		return iam.ApplicationIDNumZero, errors.Wrap("random number source reading", err)
	}

	idUint := (embeddedFieldBits & iam.ApplicationIDNumEmbeddedFieldsMask) |
		(binary.BigEndian.Uint32(idBytes) & iam.ApplicationIDNumIdentifierBitsMask)
	return iam.ApplicationIDNum(idUint), nil
}
