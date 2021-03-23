package iam

import (
	"encoding/binary"
	"strings"

	azfl "github.com/alloyzeus/go-azfl/azfl"
	azer "github.com/alloyzeus/go-azfl/azfl/azer"
	"github.com/alloyzeus/go-azfl/azfl/errors"
)

//region IDNum

// UserIDNum is a scoped identifier
// used to identify an instance of entity User.
type UserIDNum int64

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azfl.IDNum = UserIDNumZero
var _ azfl.EntityID = UserIDNumZero
var _ azer.BinFieldUnmarshalable = &_UserIDNumZeroVar
var _ azfl.UserIDNum = UserIDNumZero

// UserIDNumSignificantBitsMask is used to
// extract significant bits from an instance of UserIDNum.
const UserIDNumSignificantBitsMask uint64 = 0b11111111_11111111_11111111_11111111_11111111_11111111

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

// UserIDNumFromAZERBinField creates UserIDNum from
// its azer-bin-field form.
func UserIDNumFromAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (idNum UserIDNum, readLen int, err error) {
	if typeHint != azer.BinDataTypeUnspecified && typeHint != azer.BinDataTypeInt64 {
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
// with azfl.IDNum.
func (UserIDNum) AZIDNum() {}

// AZEntityID is required for conformance
// with azfl.EntityID.
func (UserIDNum) AZEntityID() {}

// AZUserIDNum is required for conformance
// with azfl.UserIDNum.
func (UserIDNum) AZUserIDNum() {}

// IsZero is required as UserIDNum is a value-object.
func (idNum UserIDNum) IsZero() bool {
	return idNum == UserIDNumZero
}

// IsValid returns true if this instance is valid independently
// as an UserIDNum. It doesn't tell whether it refers to
// a valid instance of User.
//
// To elaborate, validity of a data depends on the perspective of the user.
// For example, age 1000 is a valid as an instance of age, but on the context
// of human living age, we can consider it as invalid.
//
// To use some analogy, a ticket has a date of validity for today, but
// after it got checked in to the counter, it turns out that its serial number
// is not registered in the issuer's database. The ticket claims that it's
// valid, but it's considered invalid because it's a fake.
//
// Similarly, what is considered valid in this context here is that the data
// contained in this instance doesn't break any rule for an instance of
// UserIDNum. Whether the instance is valid for certain context,
// it requires case-by-case validation which is out of the scope of this
// method.
func (idNum UserIDNum) IsValid() bool {
	return int64(idNum) > 0 &&
		(uint64(idNum)&UserIDNumSignificantBitsMask) != 0
}

// IsNotValid returns the negation of value returned by IsValid().
func (idNum UserIDNum) IsNotValid() bool {
	return !idNum.IsValid()
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

// AZERBinField is required for conformance
// with azfl.IDNum.
func (idNum UserIDNum) AZERBinField() ([]byte, azer.BinDataType) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(idNum))
	return b, azer.BinDataTypeInt64
}

// UnmarshalAZERBinField is required for conformance
// with azer.BinFieldUnmarshalable.
func (idNum *UserIDNum) UnmarshalAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (readLen int, err error) {
	i, readLen, err := UserIDNumFromAZERBinField(b, typeHint)
	if err == nil {
		*idNum = i
	}
	return readLen, err
}

// IsBot returns true if
// the User instance this UserIDNum is for
// is a Bot User.
//
// Bot account is ....
func (idNum UserIDNum) IsBot() bool {
	return idNum.IsValid() && idNum.HasBotBits()
}

const _UserIDNumBotMask = 0b1000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
const _UserIDNumBotBits = 0b1000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000

// HasBotBits is only checking the bits
// without validating other information.
func (idNum UserIDNum) HasBotBits() bool {
	return (uint64(idNum) &
		_UserIDNumBotMask) ==
		_UserIDNumBotBits
}

type UserIDNumError interface {
	error
	UserIDNumError()
}

func (idNum UserIDNum) IsNormalAccount() bool {
	return idNum.IsValid() && !idNum.HasBotBits()
}

func (idNum UserIDNum) IsServiceAccount() bool {
	return idNum.IsBot()
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
var _ azfl.RefKey = _UserRefKeyZero
var _ azfl.EntityRefKey = _UserRefKeyZero
var _ azfl.UserRefKey = _UserRefKeyZero

const _UserRefKeyZero = UserRefKey(UserIDNumZero)

var _UserRefKeyZeroVar = _UserRefKeyZero

// UserRefKeyZero returns
// a zero-valued instance of UserRefKey.
func UserRefKeyZero() UserRefKey {
	return _UserRefKeyZero
}

// AZRefKey is required for conformance with azfl.RefKey.
func (UserRefKey) AZRefKey() {}

// AZEntityRefKey is required for conformance
// with azfl.EntityRefKey.
func (UserRefKey) AZEntityRefKey() {}

// IDNum returns the scoped identifier of the entity.
func (refKey UserRefKey) IDNum() UserIDNum {
	return UserIDNum(refKey)
}

// IDNumPtr returns a pointer to a copy of the IDNum if it's considered valid
// otherwise it returns nil.
func (refKey UserRefKey) IDNumPtr() *UserIDNum {
	if refKey.IsNotValid() {
		return nil
	}
	i := refKey.IDNum()
	return &i
}

// AZIDNum is required for conformance with azfl.RefKey.
func (refKey UserRefKey) AZIDNum() azfl.IDNum {
	return UserIDNum(refKey)
}

// UserIDNum is required for conformance
// with azfl.UserRefKey.
func (refKey UserRefKey) UserIDNum() azfl.UserIDNum {
	return UserIDNum(refKey)
}

// IsZero is required as UserRefKey is a value-object.
func (refKey UserRefKey) IsZero() bool {
	return UserIDNum(refKey) == UserIDNumZero
}

// IsValid returns true if this instance is valid independently as a ref-key.
// It doesn't tell whether it refers to a valid instance of User.
func (refKey UserRefKey) IsValid() bool {
	return UserIDNum(refKey).IsValid()
}

// IsNotValid returns the negation of value returned by IsValid().
func (refKey UserRefKey) IsNotValid() bool {
	return !refKey.IsValid()
}

// Equals is required for conformance with azfl.EntityRefKey.
func (refKey UserRefKey) Equals(other interface{}) bool {
	if x, ok := other.(UserRefKey); ok {
		return x == refKey
	}
	if x, _ := other.(*UserRefKey); x != nil {
		return *x == refKey
	}
	return false
}

// Equal is required for conformance with azfl.EntityRefKey.
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

func (refKey UserRefKey) AZERBin() []byte {
	b := make([]byte, 8+1)
	b[0] = azer.BinDataTypeInt64.Byte()
	binary.BigEndian.PutUint64(b[1:], uint64(refKey))
	return b
}

func UserRefKeyFromAZERBin(b []byte) (refKey UserRefKey, readLen int, err error) {
	typ, err := azer.BinDataTypeFromByte(b[0])
	if err != nil {
		return _UserRefKeyZero, 0,
			errors.ArgWrap("", "type parsing", err)
	}
	if typ != azer.BinDataTypeInt64 {
		return _UserRefKeyZero, 0,
			errors.Arg("", errors.EntMsg("type", "unsupported"))
	}

	i, readLen, err := UserRefKeyFromAZERBinField(b[1:], typ)
	if err != nil {
		return _UserRefKeyZero, 0,
			errors.ArgWrap("", "idnum data parsing", err)
	}

	return UserRefKey(i), 1 + readLen, nil
}

// UnmarshalAZERBin is required for conformance
// with azfl.BinFieldUnmarshalable.
func (refKey *UserRefKey) UnmarshalAZERBin(b []byte) (readLen int, err error) {
	i, readLen, err := UserRefKeyFromAZERBin(b)
	if err == nil {
		*refKey = i
	}
	return readLen, err
}

func (refKey UserRefKey) AZERBinField() ([]byte, azer.BinDataType) {
	return UserIDNum(refKey).AZERBinField()
}

func UserRefKeyFromAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (refKey UserRefKey, readLen int, err error) {
	idNum, n, err := UserIDNumFromAZERBinField(b, typeHint)
	if err != nil {
		return _UserRefKeyZero, n, err
	}
	return UserRefKey(idNum), n, nil
}

// UnmarshalAZERBinField is required for conformance
// with azfl.BinFieldUnmarshalable.
func (refKey *UserRefKey) UnmarshalAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (readLen int, err error) {
	i, readLen, err := UserRefKeyFromAZERBinField(b, typeHint)
	if err == nil {
		*refKey = i
	}
	return readLen, err
}

const _UserRefKeyAZERTextPrefix = "KUs0"

// AZERText is required for conformance
// with azfl.RefKey.
func (refKey UserRefKey) AZERText() string {
	if !refKey.IsValid() {
		return ""
	}

	return _UserRefKeyAZERTextPrefix +
		azer.TextEncode(refKey.AZERBin())
}

// UserRefKeyFromAZERText creates a new instance of
// UserRefKey from its azer-text form.
func UserRefKeyFromAZERText(s string) (UserRefKey, error) {
	if s == "" {
		return UserRefKeyZero(), nil
	}
	if !strings.HasPrefix(s, _UserRefKeyAZERTextPrefix) {
		return UserRefKeyZero(),
			errors.Arg("", errors.EntMsg("prefix", "mismatch"))
	}
	s = strings.TrimPrefix(s, _UserRefKeyAZERTextPrefix)
	b, err := azer.TextDecode(s)
	if err != nil {
		return UserRefKeyZero(),
			errors.ArgWrap("", "data parsing", err)
	}
	refKey, _, err := UserRefKeyFromAZERBin(b)
	if err != nil {
		return UserRefKeyZero(),
			errors.ArgWrap("", "data decoding", err)
	}
	return refKey, nil
}

// UnmarshalAZERText is required for conformance
// with azer.TextUnmarshalable.
func (refKey *UserRefKey) UnmarshalAZERText(s string) error {
	r, err := UserRefKeyFromAZERText(s)
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalText is for compatibility with Go's encoding.TextMarshaler
func (refKey UserRefKey) MarshalText() ([]byte, error) {
	return []byte(refKey.AZERText()), nil
}

// UnmarshalText is for conformance with Go's encoding.TextUnmarshaler
func (refKey *UserRefKey) UnmarshalText(b []byte) error {
	r, err := UserRefKeyFromAZERText(string(b))
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalJSON makes this type JSON-marshalable.
func (refKey UserRefKey) MarshalJSON() ([]byte, error) {
	// We assume that there's no symbols in azer-text
	return []byte("\"" + refKey.AZERText() + "\""), nil
}

// UnmarshalJSON parses a JSON value.
func (refKey *UserRefKey) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		*refKey = UserRefKeyZero()
		return nil
	}
	i, err := UserRefKeyFromAZERText(s)
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
