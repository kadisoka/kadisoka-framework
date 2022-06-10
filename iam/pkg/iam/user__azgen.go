package iam

import (
	"crypto/rand"
	"encoding/binary"
	"strings"

	azcore "github.com/alloyzeus/go-azfl/azcore"
	azid "github.com/alloyzeus/go-azfl/azid"
	errors "github.com/alloyzeus/go-azfl/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the azfl package it is being compiled against.
// A compilation error at this line likely means your copy of the
// azfl package needs to be updated.
var _ = azcore.AZCorePackageIsVersion1

// Reference imports to suppress errors if they are not otherwise used.
var _ = azid.BinDataTypeUnspecified
var _ = strings.Compare
var _ = rand.Reader

// Entity User.

//region IDNum

// UserIDNum is a scoped identifier
// used to identify an instance of entity User.
type UserIDNum int64

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azid.IDNumMethods = UserIDNumZero
var _ azid.BinFieldUnmarshalable = &_UserIDNumZeroVar
var _ azcore.EntityIDNumMethods = UserIDNumZero
var _ azcore.UserIDNumMethods = UserIDNumZero

// UserIDNumIdentifierBitsMask is used to
// extract identifier bits from an instance of UserIDNum.
const UserIDNumIdentifierBitsMask uint64 = 0b_00000000_00000000_11111111_11111111_11111111_11111111_11111111_11111111

// UserIDNumZero is the zero value
// for UserIDNum.
const UserIDNumZero = UserIDNum(0)

// _UserIDNumZeroVar is used for testing
// pointer-based interfaces conformance.
var _UserIDNumZeroVar = UserIDNumZero

// UserIDNumFromPrimitiveValue creates an instance
// of UserIDNum from its primitive value.
func UserIDNumFromPrimitiveValue(v int64) UserIDNum {
	return UserIDNum(v)
}

// UserIDNumFromAZIDBinField creates UserIDNum from
// its azid-bin-field form.
func UserIDNumFromAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (idNum UserIDNum, readLen int, err error) {
	if typeHint != azid.BinDataTypeUnspecified && typeHint != azid.BinDataTypeInt64 {
		return UserIDNum(0), 0,
			errors.ArgMsg("typeHint", "unsupported")
	}
	i := binary.BigEndian.Uint64(b)
	return UserIDNum(i), 8, nil
}

// PrimitiveValue returns the value in its primitive type. Prefer to use
// this method instead of casting directly.
func (idNum UserIDNum) PrimitiveValue() int64 {
	return int64(idNum)
}

// AZIDNum is required for conformance
// with azid.IDNum.
func (UserIDNum) AZIDNum() {}

// AZEntityIDNum is required for conformance
// with azcore.EntityIDNum.
func (UserIDNum) AZEntityIDNum() {}

// AZUserIDNum is required for conformance
// with azcore.UserIDNum.
func (UserIDNum) AZUserIDNum() {}

// IsZero is required as UserIDNum is a value-object.
func (idNum UserIDNum) IsZero() bool {
	return idNum == UserIDNumZero
}

// IsStaticallyValid returns true if this instance is valid as an isolated value
// of UserIDNum. It doesn't tell whether it refers to
// a valid instance of User.
//
// What is considered valid in this context here is that the data
// contained in this instance doesn't break any rule for an instance of
// UserIDNum. Whether the instance is valid in a certain context,
// it requires case-by-case validation which is out of the scope of this
// method.
//
// For example, age 1000 is a valid as an instance of age, but in the context
// of human living age, we can consider it as invalid.
//
// Another example, a ticket has a date of validity for today, but
// after it got checked in to the counter, it turns out that its serial number
// is not registered in the issuer's database. The ticket claims that it's
// valid, but it's considered invalid because it's a fake.
func (idNum UserIDNum) IsStaticallyValid() bool {
	return int64(idNum) > 0 &&
		(uint64(idNum)&UserIDNumIdentifierBitsMask) != 0
}

// IsNotStaticallyValid returns the negation of value returned by IsStaticallyValid.
func (idNum UserIDNum) IsNotStaticallyValid() bool {
	return !idNum.IsStaticallyValid()
}

// Equals is required as UserIDNum is a value-object.
//
// Use EqualsUserIDNum method if the other value
// has the same type.
func (idNum UserIDNum) Equals(other interface{}) bool {
	if x, ok := other.(UserIDNum); ok {
		return x == idNum
	}
	if x, _ := other.(*UserIDNum); x != nil {
		return *x == idNum
	}
	return false
}

// Equal is a wrapper for Equals method. It is required for
// compatibility with github.com/google/go-cmp
func (idNum UserIDNum) Equal(other interface{}) bool {
	return idNum.Equals(other)
}

// EqualsUserIDNum determines if the other instance is equal
// to this instance.
func (idNum UserIDNum) EqualsUserIDNum(
	other UserIDNum,
) bool {
	return idNum == other
}

// AZIDBinField is required for conformance
// with azid.IDNum.
func (idNum UserIDNum) AZIDBinField() ([]byte, azid.BinDataType) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(idNum))
	return b, azid.BinDataTypeInt64
}

// UnmarshalAZIDBinField is required for conformance
// with azid.BinFieldUnmarshalable.
func (idNum *UserIDNum) UnmarshalAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (readLen int, err error) {
	i, readLen, err := UserIDNumFromAZIDBinField(b, typeHint)
	if err == nil {
		*idNum = i
	}
	return readLen, err
}

// Embedded fields
const (
	UserIDNumEmbeddedFieldsMask = 0b_01000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000

	UserIDNumBotMask = 0b_01000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	UserIDNumBotBits = 0b_01000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
)

// IsBot returns true if
// the User instance this UserIDNum is for
// is a Bot User.
//
// Bot account is ....
func (idNum UserIDNum) IsBot() bool {
	return idNum.IsStaticallyValid() && idNum.HasBotBits()
}

// HasBotBits is only checking the bits
// without validating other information.
func (idNum UserIDNum) HasBotBits() bool {
	return (uint64(idNum) &
		UserIDNumBotMask) ==
		UserIDNumBotBits
}

type UserIDNumError interface {
	error
	UserIDNumError()
}

//endregion

//region ID

// UserID is used to identify
// an instance of entity User system-wide.
type UserID UserIDNum

// NewUserID returns a new instance
// of UserID with the provided attribute values.
func NewUserID(
	idNum UserIDNum,
) UserID {
	return UserID(idNum)
}

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azid.ID[UserIDNum] = _UserIDZero
var _ azid.BinUnmarshalable = &_UserIDZeroVar
var _ azid.BinFieldUnmarshalable = &_UserIDZeroVar
var _ azid.TextUnmarshalable = &_UserIDZeroVar
var _ azcore.EntityID[UserIDNum] = _UserIDZero
var _ azcore.UserID[UserIDNum] = _UserIDZero

const _UserIDZero = UserID(UserIDNumZero)

var _UserIDZeroVar = _UserIDZero

// UserIDZero returns
// a zero-valued instance of UserID.
func UserIDZero() UserID {
	return _UserIDZero
}

// AZID is required for conformance with azid.ID.
func (UserID) AZID() {}

// AZEntityID is required for conformance
// with azcore.EntityID.
func (UserID) AZEntityID() {}

// IDNum returns the scoped identifier of the entity.
func (id UserID) IDNum() UserIDNum {
	return UserIDNum(id)
}

// IDNumPtr returns a pointer to a copy of the id-num if it's considered valid
// otherwise it returns nil.
func (id UserID) IDNumPtr() *UserIDNum {
	if id.IsNotStaticallyValid() {
		return nil
	}
	i := id.IDNum()
	return &i
}

// AZIDNum is required for conformance with azid.ID.
func (id UserID) AZIDNum() UserIDNum {
	return UserIDNum(id)
}

// UserIDNum is required for conformance
// with azcore.UserID.
func (id UserID) UserIDNum() UserIDNum {
	return UserIDNum(id)
}

// IsZero is required as UserID is a value-object.
func (id UserID) IsZero() bool {
	return UserIDNum(id) == UserIDNumZero
}

// IsStaticallyValid returns true if this instance is valid as an isolated value
// of UserID.
// It doesn't tell whether it refers to a valid instance of User.
func (id UserID) IsStaticallyValid() bool {
	return UserIDNum(id).IsStaticallyValid()
}

// IsNotStaticallyValid returns the negation of value returned by IsStaticallyValid.
func (id UserID) IsNotStaticallyValid() bool {
	return !id.IsStaticallyValid()
}

// Equals is required for conformance with azcore.EntityID.
func (id UserID) Equals(other interface{}) bool {
	if x, ok := other.(UserID); ok {
		return x == id
	}
	if x, _ := other.(*UserID); x != nil {
		return *x == id
	}
	return false
}

// Equal is required for conformance with azcore.EntityID.
func (id UserID) Equal(other interface{}) bool {
	return id.Equals(other)
}

// EqualsUserID returs true
// if the other value has the same attributes as id.
func (id UserID) EqualsUserID(
	other UserID,
) bool {
	return other == id
}

func (id UserID) AZIDBin() []byte {
	b := make([]byte, 8+1)
	b[0] = azid.BinDataTypeInt64.Byte()
	binary.BigEndian.PutUint64(b[1:], uint64(id))
	return b
}

func UserIDFromAZIDBin(b []byte) (id UserID, readLen int, err error) {
	typ, err := azid.BinDataTypeFromByte(b[0])
	if err != nil {
		return _UserIDZero, 0,
			errors.ArgWrap("", "type parsing", err)
	}
	if typ != azid.BinDataTypeInt64 {
		return _UserIDZero, 0,
			errors.Arg("", errors.EntMsg("type", "unsupported"))
	}

	i, readLen, err := UserIDFromAZIDBinField(b[1:], typ)
	if err != nil {
		return _UserIDZero, 0,
			errors.ArgWrap("", "id-num data parsing", err)
	}

	return UserID(i), 1 + readLen, nil
}

// UnmarshalAZIDBin is required for conformance
// with azcore.BinFieldUnmarshalable.
func (id *UserID) UnmarshalAZIDBin(b []byte) (readLen int, err error) {
	i, readLen, err := UserIDFromAZIDBin(b)
	if err == nil {
		*id = i
	}
	return readLen, err
}

func (id UserID) AZIDBinField() ([]byte, azid.BinDataType) {
	return UserIDNum(id).AZIDBinField()
}

func UserIDFromAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (id UserID, readLen int, err error) {
	idNum, n, err := UserIDNumFromAZIDBinField(b, typeHint)
	if err != nil {
		return _UserIDZero, n, err
	}
	return UserID(idNum), n, nil
}

// UnmarshalAZIDBinField is required for conformance
// with azcore.BinFieldUnmarshalable.
func (id *UserID) UnmarshalAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (readLen int, err error) {
	i, readLen, err := UserIDFromAZIDBinField(b, typeHint)
	if err == nil {
		*id = i
	}
	return readLen, err
}

const _UserIDAZIDTextPrefix = "KUs0"

// AZIDText is required for conformance
// with azid.ID.
func (id UserID) AZIDText() string {
	if !id.IsStaticallyValid() {
		return ""
	}

	return _UserIDAZIDTextPrefix +
		azid.TextEncode(id.AZIDBin())
}

// UserIDFromAZIDText creates a new instance of
// UserID from its azid-text form.
func UserIDFromAZIDText(s string) (UserID, error) {
	if s == "" {
		return UserIDZero(), nil
	}
	if !strings.HasPrefix(s, _UserIDAZIDTextPrefix) {
		return UserIDZero(),
			errors.Arg("", errors.EntMsg("prefix", "mismatch"))
	}
	s = strings.TrimPrefix(s, _UserIDAZIDTextPrefix)
	b, err := azid.TextDecode(s)
	if err != nil {
		return UserIDZero(),
			errors.ArgWrap("", "data parsing", err)
	}
	id, _, err := UserIDFromAZIDBin(b)
	if err != nil {
		return UserIDZero(),
			errors.ArgWrap("", "data decoding", err)
	}
	return id, nil
}

// UnmarshalAZIDText is required for conformance
// with azid.TextUnmarshalable.
func (id *UserID) UnmarshalAZIDText(s string) error {
	r, err := UserIDFromAZIDText(s)
	if err == nil {
		*id = r
	}
	return err
}

// MarshalText is for compatibility with Go's encoding.TextMarshaler
func (id UserID) MarshalText() ([]byte, error) {
	return []byte(id.AZIDText()), nil
}

// UnmarshalText is for conformance with Go's encoding.TextUnmarshaler
func (id *UserID) UnmarshalText(b []byte) error {
	r, err := UserIDFromAZIDText(string(b))
	if err == nil {
		*id = r
	}
	return err
}

// MarshalJSON makes this type JSON-marshalable.
func (id UserID) MarshalJSON() ([]byte, error) {
	// We assume that there are no symbols in azid-text
	return []byte("\"" + id.AZIDText() + "\""), nil
}

// UnmarshalJSON parses a JSON value.
func (id *UserID) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		*id = UserIDZero()
		return nil
	}
	i, err := UserIDFromAZIDText(s)
	if err == nil {
		*id = i
	}
	return err
}

// UserIDService abstracts
// UserID-related services.
type UserIDService interface {
	// IsUserID is to check if the ref-key is
	// trully registered to system. It does not check whether the instance
	// is active or not.
	IsUserIDRegistered(id UserID) bool
}

// UserIDError defines an interface for all
// UserID-related errors.
type UserIDError interface {
	error
	UserIDError()
}

//endregion

//region Instance

// UserInstanceService is a service which
// provides methods to manipulate an instance of User.
type UserInstanceService interface {
	UserInstanceInfoService
}

// UserInstanceInfoService is a service which
// provides access to instances metadata.
type UserInstanceInfoService interface {
	// GetUserInstanceInfo checks if the provided
	// ref-key is valid and whether the instance is deleted.
	//
	// This method returns nil if the id is not referencing to any valid
	// instance.
	GetUserInstanceInfo(
		inputCtx CallInputContext,
		id UserID,
	) (*UserInstanceInfo, error)
}

// UserInstanceInfo holds information about
// an instance of User.
type UserInstanceInfo struct {
	RevisionNumber int32

	// Deletion holds information about the deletion of the instance. If
	// the instance has not been deleted, this field value will be nil.
	Deletion *UserInstanceDeletionInfo
}

// UserInstanceInfoZero returns an instance of
// UserInstanceInfo with attributes set their respective zero
// value.
func UserInstanceInfoZero() UserInstanceInfo {
	return UserInstanceInfo{}
}

// IsActive returns true if the instance is considered as active.
func (instInfo UserInstanceInfo) IsActive() bool {
	// Note: we will check other flags in the future, but that's said,
	// deleted instance is considered inactive.
	return !instInfo.IsDeleted()
}

// IsDeleted returns true if the instance was deleted.
func (instInfo UserInstanceInfo) IsDeleted() bool {
	return instInfo.Deletion != nil && instInfo.Deletion.Deleted
}

//----

// UserInstanceDeletionInfo holds information about
// the deletion of an instance if the instance has been deleted.
type UserInstanceDeletionInfo struct {
	Deleted       bool
	DeletionNotes string
}

//----

// UserInstanceServiceInternal is a service which provides
// methods for manipulating instances of User. Declared for
// internal use within a process, this interface contains methods that
// available to be called from another part of a process.
type UserInstanceServiceInternal interface {
	CreateUserInstanceInternal(
		inputCtx CallInputContext,
		input UserInstanceCreationInput,
	) (id UserID, initialState UserInstanceInfo, err error)

	// DeleteUserInstanceInternal deletes an instance of
	// User entity based identfied by refOfInstToDel.
	// The returned instanceMutated will have the value of
	// true if this particular call resulted the deletion of the instance and
	// it will have the value of false of subsequent calls to this method.
	DeleteUserInstanceInternal(
		inputCtx CallInputContext,
		refOfInstToDel UserID,
		input UserInstanceDeletionInput,
	) (instanceMutated bool, currentState UserInstanceInfo, err error)
}

// UserInstanceCreationInput contains data to be passed
// as an argument when invoking the method CreateUserInstanceInternal
// of UserInstanceServiceInternal.
type UserInstanceCreationInput struct {
}

// UserInstanceDeletionInput contains data to be passed
// as an argument when invoking the method DeleteUserInstanceInternal
// of UserInstanceServiceInternal.
type UserInstanceDeletionInput struct {
	DeletionNotes string
}

//endregion

//region Service

// UserService provides a contract
// for methods related to entity User.
type UserService interface {
	// AZxEntityService

	UserInstanceService
}

// UserServiceClient is the interface for
// clients of UserService.
type UserServiceClient interface {
	UserService
}

//endregion
