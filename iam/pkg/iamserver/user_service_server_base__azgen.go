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

const userDBTableName = "user_dt"

// UserServiceServerbase is the server-side
// base implementation of UserService.
type UserServiceServerBase struct {
	db *sqlx.DB

	registeredUserIDNumCache *lru.ARCCache
	deletedUserIDNumCache    *lru.ARCCache
}

// Interface conformance assertions.
var _ iam.UserService = &UserServiceServerBase{}
var _ iam.UserRefKeyService = &UserServiceServerBase{}

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
	callCtx iam.CallContext,
	refKey iam.UserRefKey,
) (*iam.UserInstanceInfo, error) {
	//TODO: access control
	return srv.getUserInstanceInfoNoAC(callCtx, refKey)
}

func (srv *UserServiceServerBase) getUserInstanceInfoNoAC(
	callCtx iam.CallContext,
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
		return iam.UserIDNumZero, errors.ArgWrap("", "random source reading", err)
	}

	idUint := (embeddedFieldBits & iam.UserIDNumEmbeddedFieldsMask) |
		(binary.BigEndian.Uint64(idBytes) & iam.UserIDNumIdentifierBitsMask)
	return iam.UserIDNum(idUint), nil
}
