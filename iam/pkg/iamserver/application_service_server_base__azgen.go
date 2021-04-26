package iamserver

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"

	errors "github.com/alloyzeus/go-azfl/azfl/errors"
	"github.com/doug-martin/goqu/v9"
	lru "github.com/hashicorp/golang-lru"
	"github.com/jmoiron/sqlx"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

const applicationDBTableName = "application_dt"

// ApplicationServiceServerbase is the server-side
// base implementation of ApplicationService.
type ApplicationServiceServerBase struct {
	db *sqlx.DB

	registeredApplicationIDNumCache *lru.ARCCache
	deletedApplicationIDNumCache    *lru.ARCCache
}

// Interface conformance assertions.
var _ iam.ApplicationService = &ApplicationServiceServerBase{}
var _ iam.ApplicationRefKeyService = &ApplicationServiceServerBase{}

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
	callCtx iam.CallContext,
	refKey iam.ApplicationRefKey,
) (*iam.ApplicationInstanceInfo, error) {
	//TODO: access control
	return srv.getApplicationInstanceInfoNoAC(callCtx, refKey)
}

func (srv *ApplicationServiceServerBase) getApplicationInstanceInfoNoAC(
	callCtx iam.CallContext,
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

// GenerateApplicationIDNum generates a new iam.ApplicationIDNum.
// Note that this function does not consulting any database nor registry.
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
		return iam.ApplicationIDNumZero, errors.ArgWrap("", "random source reading", err)
	}

	idUint := (embeddedFieldBits & iam.ApplicationIDNumEmbeddedFieldsMask) |
		(binary.BigEndian.Uint32(idBytes) & iam.ApplicationIDNumIdentifierBitsMask)
	return iam.ApplicationIDNum(idUint), nil
}
