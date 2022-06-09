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

//region RefKey

// UserRefKey is used to identify
// an instance of entity User system-wide.
type UserRefKey UserIDNum

// NewUserRefKey returns a new instance
// of UserRefKey with the provided attribute values.
func NewUserRefKey(
	idNum UserIDNum,
) UserRefKey {
	return UserRefKey(idNum)
}

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azid.RefKey[UserIDNum] = _UserRefKeyZero
var _ azid.BinUnmarshalable = &_UserRefKeyZeroVar
var _ azid.BinFieldUnmarshalable = &_UserRefKeyZeroVar
var _ azid.TextUnmarshalable = &_UserRefKeyZeroVar
var _ azcore.EntityRefKey[UserIDNum] = _UserRefKeyZero
var _ azcore.UserRefKey[UserIDNum] = _UserRefKeyZero

const _UserRefKeyZero = UserRefKey(UserIDNumZero)

var _UserRefKeyZeroVar = _UserRefKeyZero

// UserRefKeyZero returns
// a zero-valued instance of UserRefKey.
func UserRefKeyZero() UserRefKey {
	return _UserRefKeyZero
}

// AZRefKey is required for conformance with azid.RefKey.
func (UserRefKey) AZRefKey() {}

// AZEntityRefKey is required for conformance
// with azcore.EntityRefKey.
func (UserRefKey) AZEntityRefKey() {}

// IDNum returns the scoped identifier of the entity.
func (refKey UserRefKey) IDNum() UserIDNum {
	return UserIDNum(refKey)
}

// IDNumPtr returns a pointer to a copy of the id-num if it's considered valid
// otherwise it returns nil.
func (refKey UserRefKey) IDNumPtr() *UserIDNum {
	if refKey.IsNotStaticallyValid() {
		return nil
	}
	i := refKey.IDNum()
	return &i
}

// AZIDNum is required for conformance with azid.RefKey.
func (refKey UserRefKey) AZIDNum() UserIDNum {
	return UserIDNum(refKey)
}

// UserIDNum is required for conformance
// with azcore.UserRefKey.
func (refKey UserRefKey) UserIDNum() UserIDNum {
	return UserIDNum(refKey)
}

// IsZero is required as UserRefKey is a value-object.
func (refKey UserRefKey) IsZero() bool {
	return UserIDNum(refKey) == UserIDNumZero
}

// IsStaticallyValid returns true if this instance is valid as an isolated value
// of UserRefKey.
// It doesn't tell whether it refers to a valid instance of User.
func (refKey UserRefKey) IsStaticallyValid() bool {
	return UserIDNum(refKey).IsStaticallyValid()
}

// IsNotStaticallyValid returns the negation of value returned by IsStaticallyValid.
func (refKey UserRefKey) IsNotStaticallyValid() bool {
	return !refKey.IsStaticallyValid()
}

// Equals is required for conformance with azcore.EntityRefKey.
func (refKey UserRefKey) Equals(other interface{}) bool {
	if x, ok := other.(UserRefKey); ok {
		return x == refKey
	}
	if x, _ := other.(*UserRefKey); x != nil {
		return *x == refKey
	}
	return false
}

// Equal is required for conformance with azcore.EntityRefKey.
func (refKey UserRefKey) Equal(other interface{}) bool {
	return refKey.Equals(other)
}

// EqualsUserRefKey returs true
// if the other value has the same attributes as refKey.
func (refKey UserRefKey) EqualsUserRefKey(
	other UserRefKey,
) bool {
	return other == refKey
}

func (refKey UserRefKey) AZIDBin() []byte {
	b := make([]byte, 8+1)
	b[0] = azid.BinDataTypeInt64.Byte()
	binary.BigEndian.PutUint64(b[1:], uint64(refKey))
	return b
}

func UserRefKeyFromAZIDBin(b []byte) (refKey UserRefKey, readLen int, err error) {
	typ, err := azid.BinDataTypeFromByte(b[0])
	if err != nil {
		return _UserRefKeyZero, 0,
			errors.ArgWrap("", "type parsing", err)
	}
	if typ != azid.BinDataTypeInt64 {
		return _UserRefKeyZero, 0,
			errors.Arg("", errors.EntMsg("type", "unsupported"))
	}

	i, readLen, err := UserRefKeyFromAZIDBinField(b[1:], typ)
	if err != nil {
		return _UserRefKeyZero, 0,
			errors.ArgWrap("", "id-num data parsing", err)
	}

	return UserRefKey(i), 1 + readLen, nil
}

// UnmarshalAZIDBin is required for conformance
// with azcore.BinFieldUnmarshalable.
func (refKey *UserRefKey) UnmarshalAZIDBin(b []byte) (readLen int, err error) {
	i, readLen, err := UserRefKeyFromAZIDBin(b)
	if err == nil {
		*refKey = i
	}
	return readLen, err
}

func (refKey UserRefKey) AZIDBinField() ([]byte, azid.BinDataType) {
	return UserIDNum(refKey).AZIDBinField()
}

func UserRefKeyFromAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (refKey UserRefKey, readLen int, err error) {
	idNum, n, err := UserIDNumFromAZIDBinField(b, typeHint)
	if err != nil {
		return _UserRefKeyZero, n, err
	}
	return UserRefKey(idNum), n, nil
}

// UnmarshalAZIDBinField is required for conformance
// with azcore.BinFieldUnmarshalable.
func (refKey *UserRefKey) UnmarshalAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (readLen int, err error) {
	i, readLen, err := UserRefKeyFromAZIDBinField(b, typeHint)
	if err == nil {
		*refKey = i
	}
	return readLen, err
}

const _UserRefKeyAZIDTextPrefix = "KUs0"

// AZIDText is required for conformance
// with azid.RefKey.
func (refKey UserRefKey) AZIDText() string {
	if !refKey.IsStaticallyValid() {
		return ""
	}

	return _UserRefKeyAZIDTextPrefix +
		azid.TextEncode(refKey.AZIDBin())
}

// UserRefKeyFromAZIDText creates a new instance of
// UserRefKey from its azid-text form.
func UserRefKeyFromAZIDText(s string) (UserRefKey, error) {
	if s == "" {
		return UserRefKeyZero(), nil
	}
	if !strings.HasPrefix(s, _UserRefKeyAZIDTextPrefix) {
		return UserRefKeyZero(),
			errors.Arg("", errors.EntMsg("prefix", "mismatch"))
	}
	s = strings.TrimPrefix(s, _UserRefKeyAZIDTextPrefix)
	b, err := azid.TextDecode(s)
	if err != nil {
		return UserRefKeyZero(),
			errors.ArgWrap("", "data parsing", err)
	}
	refKey, _, err := UserRefKeyFromAZIDBin(b)
	if err != nil {
		return UserRefKeyZero(),
			errors.ArgWrap("", "data decoding", err)
	}
	return refKey, nil
}

// UnmarshalAZIDText is required for conformance
// with azid.TextUnmarshalable.
func (refKey *UserRefKey) UnmarshalAZIDText(s string) error {
	r, err := UserRefKeyFromAZIDText(s)
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalText is for compatibility with Go's encoding.TextMarshaler
func (refKey UserRefKey) MarshalText() ([]byte, error) {
	return []byte(refKey.AZIDText()), nil
}

// UnmarshalText is for conformance with Go's encoding.TextUnmarshaler
func (refKey *UserRefKey) UnmarshalText(b []byte) error {
	r, err := UserRefKeyFromAZIDText(string(b))
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalJSON makes this type JSON-marshalable.
func (refKey UserRefKey) MarshalJSON() ([]byte, error) {
	// We assume that there are no symbols in azid-text
	return []byte("\"" + refKey.AZIDText() + "\""), nil
}

// UnmarshalJSON parses a JSON value.
func (refKey *UserRefKey) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		*refKey = UserRefKeyZero()
		return nil
	}
	i, err := UserRefKeyFromAZIDText(s)
	if err == nil {
		*refKey = i
	}
	return err
}

// UserRefKeyService abstracts
// UserRefKey-related services.
type UserRefKeyService interface {
	// IsUserRefKey is to check if the ref-key is
	// trully registered to system. It does not check whether the instance
	// is active or not.
	IsUserRefKeyRegistered(refKey UserRefKey) bool
}

// UserRefKeyError defines an interface for all
// UserRefKey-related errors.
type UserRefKeyError interface {
	error
	UserRefKeyError()
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
	// This method returns nil if the refKey is not referencing to any valid
	// instance.
	GetUserInstanceInfo(
		opInputCtx OpInputContext,
		refKey UserRefKey,
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
		opInputCtx OpInputContext,
		input UserInstanceCreationInput,
	) (refKey UserRefKey, initialState UserInstanceInfo, err error)

	// DeleteUserInstanceInternal deletes an instance of
	// User entity based identfied by refOfInstToDel.
	// The returned instanceMutated will have the value of
	// true if this particular call resulted the deletion of the instance and
	// it will have the value of false of subsequent calls to this method.
	DeleteUserInstanceInternal(
		opInputCtx OpInputContext,
		refOfInstToDel UserRefKey,
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
