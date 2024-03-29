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

	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
)

const (
	applicationDBTableName           = "application_dt"
	applicationDBTablePrimaryKeyName = applicationDBTableName + "_pkey"
)

// ApplicationServiceServerbase is the server-side
// base implementation of ApplicationService.
type ApplicationServiceServerBase struct {
	db *sqlx.DB

	deletionTxHook func(iam.CallInputContext, *sqlx.Tx) error

	registeredApplicationIDNumCache *lru.ARCCache
	deletedApplicationIDNumCache    *lru.ARCCache
}

// Interface conformance assertions.
var (
	_ iam.ApplicationService                 = &ApplicationServiceServerBase{}
	_ iam.ApplicationIDService               = &ApplicationServiceServerBase{}
	_ iam.ApplicationInstanceServiceInternal = &ApplicationServiceServerBase{}
)

func (srv *ApplicationServiceServerBase) IsApplicationIDRegistered(
	id iam.ApplicationID,
) bool {
	idNum := id.IDNum()

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

// GetApplicationInstanceInfo retrieves the state of an Application
// instance. It includes the existence of the ID, and whether the instance
// has been deleted.
//
// If it's required only to determine the existence of the ID,
// IsApplicationIDRegistered is generally more efficient.
func (srv *ApplicationServiceServerBase) GetApplicationInstanceInfo(
	inputCtx iam.CallInputContext,
	id iam.ApplicationID,
) (*iam.ApplicationInstanceInfo, error) {
	//TODO: access control
	return srv.getApplicationInstanceInfoInsecure(inputCtx, id)
}

func (srv *ApplicationServiceServerBase) getApplicationInstanceInfoInsecure(
	inputCtx iam.CallInputContext,
	id iam.ApplicationID,
) (*iam.ApplicationInstanceInfo, error) {
	idRegistered := false
	idRegisteredCacheHit := false
	instDeleted := false
	instDeletionCacheHit := false

	// Look up for the ID num in the cache.
	if _, idRegistered = srv.registeredApplicationIDNumCache.Get(id); idRegistered {
		// ID num is positively registered.
		idRegisteredCacheHit = true
	}

	// Look up in the deletion cache
	if _, instDeleted = srv.deletedApplicationIDNumCache.Get(id); instDeleted {
		// Instance is positively deleted
		instDeletionCacheHit = true
	}

	if idRegisteredCacheHit && instDeletionCacheHit {
		if !idRegistered {
			return nil, nil
		}
		var deletion *iam.ApplicationInstanceDeletionInfo
		if instDeleted {
			deletion = &iam.ApplicationInstanceDeletionInfo{
				Deleted_: true}
		}
		//TODO: populate revision number
		return &iam.ApplicationInstanceInfo{
			Deletion_: deletion,
		}, nil
	}

	var err error
	idRegistered, instDeleted, err = srv.
		getApplicationInstanceStateByIDNum(id.IDNum())
	if err != nil {
		return nil, err
	}

	if !idRegisteredCacheHit && idRegistered {
		srv.registeredApplicationIDNumCache.Add(id, nil)
	}
	if !instDeletionCacheHit && instDeleted {
		srv.deletedApplicationIDNumCache.Add(id, nil)
	}

	if !idRegistered {
		return nil, nil
	}

	var deletion *iam.ApplicationInstanceDeletionInfo
	if instDeleted {
		//TODO: deletion notes. store the notes as the value in the cache
		deletion = &iam.ApplicationInstanceDeletionInfo{
			Deleted_: true}
	}

	//TODO: populate revision number
	return &iam.ApplicationInstanceInfo{
		RevisionNumber_: -1,
		Deletion_:       deletion,
	}, nil
}

func (srv *ApplicationServiceServerBase) getApplicationInstanceStateByIDNum(
	idNum iam.ApplicationIDNum,
) (idRegistered, instanceDeleted bool, err error) {
	sqlString, _, _ := goqu.From(applicationDBTableName).
		Select(
			goqu.Case().
				When(goqu.C(applicationDBColMDDeletionTimestamp).IsNull(), false).
				Else(true).
				As("deleted"),
		).
		Where(
			goqu.C(applicationDBColIDNum).Eq(idNum.PrimitiveValue()),
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
	inputCtx iam.CallInputContext,
	input iam.ApplicationInstanceCreationInput,
) (id iam.ApplicationID, initialState iam.ApplicationInstanceInfo, err error) {
	//TODO: access control

	id, err = srv.createApplicationInstanceInsecure(inputCtx)

	//TODO: revision number
	return id, iam.ApplicationInstanceInfo{
		RevisionNumber_: -1,
	}, err
}

func (srv *ApplicationServiceServerBase) createApplicationInstanceInsecure(
	inputCtx iam.CallInputContext,
) (iam.ApplicationID, error) {
	ctxAuth := inputCtx.Authorization()

	const attemptNumMax = 5

	var err error
	var newInstanceIDNum iam.ApplicationIDNum
	cTime := inputCtx.CallInputMetadata().ReceiveTime

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
					applicationDBColIDNum:                newInstanceIDNum,
					applicationDBColMDCreationTimestamp:  cTime,
					applicationDBColMDCreationUserID:     ctxAuth.UserIDNumPtr(),
					applicationDBColMDCreationTerminalID: ctxAuth.TerminalIDNumPtr(),
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
			pqErr.Constraint == applicationDBTablePrimaryKeyName {
			if attemptNum >= attemptNumMax {
				return iam.ApplicationIDZero(), errors.Wrap("insert max attempts", err)
			}
			continue
		}

		return iam.ApplicationIDZero(), errors.Wrap("insert", err)
	}

	//TODO: update caches, emit an event

	return iam.NewApplicationID(newInstanceIDNum), nil
}

func (srv *ApplicationServiceServerBase) DeleteApplicationInstanceInternal(
	inputCtx iam.CallInputContext,
	toDelete iam.ApplicationID,
	input iam.ApplicationInstanceDeletionInput,
) (justDeleted bool, currentState iam.ApplicationInstanceInfo, err error) {
	if inputCtx == nil {
		return false, iam.ApplicationInstanceInfoZero(), nil
	}

	//TODO: access control

	return srv.deleteApplicationInstanceInsecure(inputCtx, toDelete, input)
}

func (srv *ApplicationServiceServerBase) deleteApplicationInstanceInsecure(
	inputCtx iam.CallInputContext,
	toDelete iam.ApplicationID,
	input iam.ApplicationInstanceDeletionInput,
) (justDeleted bool, currentState iam.ApplicationInstanceInfo, err error) {
	ctxAuth := inputCtx.Authorization()
	ctxTime := inputCtx.CallInputMetadata().ReceiveTime

	err = doTx(srv.db, func(dbTx *sqlx.Tx) error {
		sqlString, _, _ := goqu.
			From(applicationDBTableName).
			Where(
				goqu.C(applicationDBColIDNum).Eq(ctxAuth.UserIDNum().PrimitiveValue()),
				goqu.C(applicationDBColMDDeletionTimestamp).IsNull(),
			).
			Update().
			Set(
				goqu.Record{
					applicationDBColMDDeletionTimestamp:  ctxTime,
					applicationDBColMDDeletionTerminalID: ctxAuth.TerminalIDNum().PrimitiveValue(),
					applicationDBColMDDeletionUserID:     ctxAuth.UserIDNum().PrimitiveValue(),
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
		justDeleted = n == 1

		if srv.deletionTxHook != nil {
			return srv.deletionTxHook(inputCtx, dbTx)
		}

		return nil
	})
	if err != nil {
		return false, iam.ApplicationInstanceInfoZero(), err
	}

	var deletion *iam.ApplicationInstanceDeletionInfo
	if justDeleted {
		deletion = &iam.ApplicationInstanceDeletionInfo{
			Deleted_: true}
	} else {
		di, err := srv.getApplicationInstanceInfoInsecure(inputCtx, toDelete)
		if err != nil {
			return false, iam.ApplicationInstanceInfoZero(), err
		}

		if di != nil {
			deletion = di.Deletion()
		}
	}

	currentState = iam.ApplicationInstanceInfo{
		RevisionNumber_: -1, //TODO: get from the DB
		Deletion_:       deletion}

	//TODO: update caches, emit an event if there's any changes

	return justDeleted, currentState, nil
}

// Designed to perform the migration if required
//TODO: context: target version, current version (assert), prefix, etc.
func (srv *ApplicationServiceServerBase) initDataStoreInTx(dbTx *sqlx.Tx) error {
	_, err := dbTx.Exec(
		`CREATE TABLE ` + applicationDBTableName + ` ( ` +
			applicationDBColIDNum + `  integer, ` +
			applicationDBColMDCreationTimestamp + `  timestamp with time zone NOT NULL DEFAULT now(), ` +
			applicationDBColMDCreationTerminalID + `  bigint, ` +
			applicationDBColMDCreationUserID + `  bigint, ` +
			applicationDBColMDDeletionTimestamp + `  timestamp with time zone, ` +
			applicationDBColMDDeletionTerminalID + `  bigint, ` +
			applicationDBColMDDeletionUserID + `  bigint, ` +
			`CONSTRAINT ` + applicationDBTablePrimaryKeyName + ` PRIMARY KEY(` + applicationDBColIDNum + `), ` +
			`CHECK (` + applicationDBColIDNum + ` > 0) ` +
			`);`,
	)
	if err != nil {
		return err
	}
	return nil
}

const (
	applicationDBColMDCreationTimestamp  = "md_c_ts"
	applicationDBColMDCreationTerminalID = "md_c_tid"
	applicationDBColMDCreationUserID     = "md_c_uid"
	applicationDBColMDDeletionTimestamp  = "md_d_ts"
	applicationDBColMDDeletionTerminalID = "md_d_tid"
	applicationDBColMDDeletionUserID     = "md_d_uid"
	applicationDBColIDNum                = "id_num"
)

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
func GenerateApplicationIDNum(
	embeddedFieldBits uint32,
) (iam.ApplicationIDNum, error) {
	idBytes := make([]byte, 4)
	_, err := rand.Read(idBytes)
	if err != nil {
		return iam.ApplicationIDNumZero, errors.Wrap("random number source reading", err)
	}

	idUint := (embeddedFieldBits & iam.ApplicationIDNumEmbeddedFieldsMask) |
		(binary.BigEndian.Uint32(idBytes) & iam.ApplicationIDNumIdentifierBitsMask)
	return iam.ApplicationIDNum(idUint), nil
}
