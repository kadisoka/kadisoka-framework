package iam

import (
	"encoding/binary"
	"strings"

	azcore "github.com/alloyzeus/go-azcore/azcore"
	azer "github.com/alloyzeus/go-azcore/azcore/azer"
	"github.com/alloyzeus/go-azcore/azcore/errors"
)

var (
	ErrUserIDStringInvalid        = errors.Ent("user ID string", nil)
	ErrServiceUserIDStringInvalid = errors.Ent("service user ID string", nil)
)

//region ID

// UserID is a scoped identifier
// used to identify an instance of entity User.
type UserID int64

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azcore.EID = UserIDZero
var _ azcore.EntityID = UserIDZero
var _ azer.BinFieldUnmarshalable = &_UserIDZeroVar
var _ azcore.UserID = UserIDZero

// _UserIDSignificantBitsMask is used to
// extract significant bits from an instance of UserID.
const _UserIDSignificantBitsMask uint64 = 0b11111111_11111111_11111111_11111111_11111111_11111111

// UserIDZero is the zero value
// for UserID.
const UserIDZero = UserID(0)

// _UserIDZeroVar is used for testing
// pointer-based interfaces conformance.
var _UserIDZeroVar = UserIDZero

// UserIDFromPrimitiveValue creates an instance
// of UserID from its primitive value.
func UserIDFromPrimitiveValue(v int64) UserID {
	return UserID(v)
}

// UserIDFromAZERBinField creates UserID from
// its azer-bin-field form.
func UserIDFromAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (id UserID, readLen int, err error) {
	if typeHint != azer.BinDataTypeUnspecified && typeHint != azer.BinDataTypeInt64 {
		return UserID(0), 0,
			errors.ArgMsg("typeHint", "unsupported")
	}
	i := binary.BigEndian.Uint64(b)
	return UserID(i), 8, nil
}

// PrimitiveValue returns the ID in its primitive type. Prefer to use
// this method instead of casting directly.
func (id UserID) PrimitiveValue() int64 {
	return int64(id)
}

// AZEID is required for conformance
// with azcore.EID.
func (UserID) AZEID() {}

// AZEntityID is required for conformance
// with azcore.EntityID.
func (UserID) AZEntityID() {}

// AZUserID is required for conformance
// with azcore.UserID.
func (UserID) AZUserID() {}

// IsZero is required as UserID is a value-object.
func (id UserID) IsZero() bool {
	return id == UserIDZero
}

// IsValid returns true if this instance is valid independently as an ID.
// It doesn't tell whether it refers to a valid instance of User.
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
// UserID. Whether the instance is valid for certain context,
// it requires case-by-case validation which is out of the scope of this
// method.
func (id UserID) IsValid() bool {
	return int64(id) > 0 &&
		(uint64(id)&_UserIDSignificantBitsMask) != 0
}

// IsNotValid returns the negation of value returned by IsValid().
func (id UserID) IsNotValid() bool {
	return !id.IsValid()
}

// Equals is required as UserID is a value-object.
//
// Use EqualsUserID method if the other value
// has the same type.
func (id UserID) Equals(other interface{}) bool {
	if x, ok := other.(UserID); ok {
		return x == id
	}
	if x, _ := other.(*UserID); x != nil {
		return *x == id
	}
	return false
}

// Equal is a wrapper for Equals method. It is required for
// compatibility with github.com/google/go-cmp
func (id UserID) Equal(other interface{}) bool {
	return id.Equals(other)
}

// EqualsUserID determines if the other instance is equal
// to this instance.
func (id UserID) EqualsUserID(
	other UserID,
) bool {
	return id == other
}

// AZERBinField is required for conformance
// with azcore.EID.
func (id UserID) AZERBinField() ([]byte, azer.BinDataType) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(id))
	return b, azer.BinDataTypeInt64
}

// UnmarshalAZERBinField is required for conformance
// with azer.BinFieldUnmarshalable.
func (id *UserID) UnmarshalAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (readLen int, err error) {
	i, readLen, err := UserIDFromAZERBinField(b, typeHint)
	if err == nil {
		*id = i
	}
	return readLen, err
}

// IsBot returns true if the User instance
// this ID is for is a Bot User.
//
// Bot account is ....
func (id UserID) IsBot() bool {
	return id.IsValid() && id.HasBotBits()
}

// HasBotBits is only checking the bits
// without validating other information contained in the ID.
func (id UserID) HasBotBits() bool {
	return (uint64(id) &
		0b1000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000) ==
		0b1000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
}

type UserIDError interface {
	error
	UserIDError()
}

func (id UserID) IsNormalAccount() bool {
	return id.IsValid() && !id.HasBotBits()
}

func (id UserID) IsServiceAccount() bool {
	return id.IsBot()
}

//endregion

//region RefKey

// UserRefKey is used to identify
// an instance of entity User system-wide.
type UserRefKey UserID

// NewUserRefKey returns a new instance
// of UserRefKey with the provided attribute values.
func NewUserRefKey(
	id UserID,
) UserRefKey {
	return UserRefKey(id)
}

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azcore.RefKey = _UserRefKeyZero
var _ azcore.EntityRefKey = _UserRefKeyZero
var _ azcore.UserRefKey = _UserRefKeyZero

const _UserRefKeyZero = UserRefKey(UserIDZero)

var _UserRefKeyZeroVar = _UserRefKeyZero

// UserRefKeyZero returns
// a zero-valued instance of UserRefKey.
func UserRefKeyZero() UserRefKey {
	return _UserRefKeyZero
}

// AZRefKey is required for conformance with azcore.RefKey.
func (UserRefKey) AZRefKey() {}

// AZEntityRefKey is required for conformance
// with azcore.EntityRefKey.
func (UserRefKey) AZEntityRefKey() {}

// ID returns the scoped identifier of the entity.
func (refKey UserRefKey) ID() UserID {
	return UserID(refKey)
}

// IDPtr returns a pointer to a copy of the ID if it's considered valid.
func (refKey UserRefKey) IDPtr() *UserID {
	if refKey.IsNotValid() {
		return nil
	}
	i := refKey.ID()
	return &i
}

// ID is required for conformance with azcore.RefKey.
func (refKey UserRefKey) EID() azcore.EID {
	return UserID(refKey)
}

// UserID is required for conformance
// with azcore.UserRefKey.
func (refKey UserRefKey) UserID() azcore.UserID {
	return UserID(refKey)
}

// IsZero is required as UserRefKey is a value-object.
func (refKey UserRefKey) IsZero() bool {
	return UserID(refKey) == UserIDZero
}

// IsValid returns true if this instance is valid independently as a ref-key.
// It doesn't tell whether it refers to a valid instance of User.
func (refKey UserRefKey) IsValid() bool {
	return UserID(refKey).IsValid()
}

// IsNotValid returns the negation of value returned by IsValid().
func (refKey UserRefKey) IsNotValid() bool {
	return !refKey.IsValid()
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
			errors.ArgWrap("", "id data parsing", err)
	}

	return UserRefKey(i), 1 + readLen, nil
}

// UnmarshalAZERBin is required for conformance
// with azcore.BinFieldUnmarshalable.
func (refKey *UserRefKey) UnmarshalAZERBin(b []byte) (readLen int, err error) {
	i, readLen, err := UserRefKeyFromAZERBin(b)
	if err == nil {
		*refKey = i
	}
	return readLen, err
}

func (refKey UserRefKey) AZERBinField() ([]byte, azer.BinDataType) {
	return UserID(refKey).AZERBinField()
}

func UserRefKeyFromAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (refKey UserRefKey, readLen int, err error) {
	id, n, err := UserIDFromAZERBinField(b, typeHint)
	if err != nil {
		return _UserRefKeyZero, n, err
	}
	return UserRefKey(id), n, nil
}

// UnmarshalAZERBinField is required for conformance
// with azcore.BinFieldUnmarshalable.
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
// with azcore.RefKey.
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

// UserRefKeyError defines an interface for all
// UserRefKey-related errors.
type UserRefKeyError interface {
	error
	UserRefKeyError()
}

//endregion
