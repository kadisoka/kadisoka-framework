package iam

import (
	"encoding/binary"
	"strings"

	azfl "github.com/alloyzeus/go-azfl/azfl"
	azid "github.com/alloyzeus/go-azfl/azfl/azid"
	errors "github.com/alloyzeus/go-azfl/azfl/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the azfl package it is being compiled against.
// A compilation error at this line likely means your copy of the
// azfl package needs to be updated.
var _ = azfl.AZCorePackageIsVersion1

// Reference imports to suppress errors if they are not otherwise used.
var _ = azid.BinDataTypeUnspecified
var _ = strings.Compare

// Adjunct-entity Terminal of Application, User.
//
// A Terminal is an authorized instance of application. As long the the
// authorization is still valid, a terminal could be used to access the
// service within the limit of the granted authorization. An
// authorization of a Terminal will become invalid when it is expired, or
// revoked by the user, if the terminal is associated to a user, or those
// who have the permission to revoke the authorization.
//
// After a Terminal authorization is invalid, their user must re-authorize
// their instance of Application if they wish to continue using their
// instance of Application to access the service. Re-authorization of an
// instance of an Application will generate a new Terminal. A
// de-authorized Terminal is no longer usable.
//
// A sucessful authorization will generate both a new Terminal and
// an initial Session.

//region IDNum

// TerminalIDNum is a scoped identifier
// used to identify an instance of adjunct entity Terminal
// scoped within its host entity(s).
type TerminalIDNum int64

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azid.IDNum = TerminalIDNumZero
var _ azid.BinFieldUnmarshalable = &_TerminalIDNumZeroVar
var _ azfl.AdjunctEntityIDNum = TerminalIDNumZero
var _ azfl.TerminalIDNum = TerminalIDNumZero

// TerminalIDNumSignificantBitsMask is used to
// extract significant bits from an instance of TerminalIDNum.
const TerminalIDNumSignificantBitsMask uint64 = 0b11111111_11111111_11111111_11111111_11111111_11111111_11111111

// TerminalIDNumZero is the zero value for TerminalIDNum.
const TerminalIDNumZero = TerminalIDNum(0)

// _TerminalIDNumZeroVar is used for testing
// pointer-based interfaces conformance.
var _TerminalIDNumZeroVar = TerminalIDNumZero

// TerminalIDNumFromPrimitiveValue creates an instance
// of TerminalIDNum from its primitive value.
func TerminalIDNumFromPrimitiveValue(v int64) TerminalIDNum {
	return TerminalIDNum(v)
}

// TerminalIDNumFromAZIDBinField creates TerminalIDNum from
// its azid-bin form.
func TerminalIDNumFromAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (idNum TerminalIDNum, readLen int, err error) {
	if typeHint != azid.BinDataTypeUnspecified && typeHint != azid.BinDataTypeInt64 {
		return TerminalIDNum(0), 0,
			errors.ArgMsg("typeHint", "unsupported")
	}
	i := binary.BigEndian.Uint64(b)
	return TerminalIDNum(i), 8, nil
}

// PrimitiveValue returns the value in its primitive type. Prefer to use
// this method instead of casting directly.
func (idNum TerminalIDNum) PrimitiveValue() int64 {
	return int64(idNum)
}

// AZIDNum is required
// for conformance with azid.IDNum.
func (TerminalIDNum) AZIDNum() {}

// AZAdjunctEntityIDNum is required
// for conformance with azfl.AdjunctEntityIDNum.
func (TerminalIDNum) AZAdjunctEntityIDNum() {}

// AZTerminalIDNum is required for conformance
// with azfl.TerminalIDNum.
func (TerminalIDNum) AZTerminalIDNum() {}

// IsZero is required as TerminalIDNum is a value-object.
func (idNum TerminalIDNum) IsZero() bool {
	return idNum == TerminalIDNumZero
}

// IsSound returns true if this instance is valid independently
// as an TerminalIDNum. It doesn't tell whether it refers to
// a valid instance of Terminal.
func (idNum TerminalIDNum) IsSound() bool {
	return int64(idNum) > 0 &&
		(uint64(idNum)&TerminalIDNumSignificantBitsMask) != 0
}

// IsNotSound returns the negation of value returned by IsSound.
func (idNum TerminalIDNum) IsNotSound() bool {
	return !idNum.IsSound()
}

// AZIDBinField is required for conformance
// with azid.IDNum.
func (idNum TerminalIDNum) AZIDBinField() ([]byte, azid.BinDataType) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(idNum))
	return b, azid.BinDataTypeInt64
}

// UnmarshalAZIDBinField is required for conformance
// with azid.BinFieldUnmarshalable.
func (idNum *TerminalIDNum) UnmarshalAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (readLen int, err error) {
	i, readLen, err := TerminalIDNumFromAZIDBinField(b, typeHint)
	if err == nil {
		*idNum = i
	}
	return readLen, err
}

// Equals is required as TerminalIDNum is a value-object.
//
// Use EqualsTerminalIDNum method if the other value
// has the same type.
func (idNum TerminalIDNum) Equals(other interface{}) bool {
	if x, ok := other.(TerminalIDNum); ok {
		return x == idNum
	}
	if x, _ := other.(*TerminalIDNum); x != nil {
		return *x == idNum
	}
	return false
}

// Equal is a wrapper for Equals method. It is required for
// compatibility with github.com/google/go-cmp
func (idNum TerminalIDNum) Equal(other interface{}) bool {
	return idNum.Equals(other)
}

// EqualsTerminalIDNum determines if the other instance
// is equal to this instance.
func (idNum TerminalIDNum) EqualsTerminalIDNum(
	other TerminalIDNum,
) bool {
	return idNum == other
}

//endregion

//region RefKey

// TerminalRefKey is used to identify
// an instance of adjunct entity Terminal system-wide.
type TerminalRefKey struct {
	application ApplicationRefKey
	user        UserRefKey
	idNum       TerminalIDNum
}

// The total number of fields in the struct.
const _TerminalRefKeyFieldCount = 2 + 1

// NewTerminalRefKey returns a new instance
// of TerminalRefKey with the provided attribute values.
func NewTerminalRefKey(
	application ApplicationRefKey,
	user UserRefKey,
	idNum TerminalIDNum,
) TerminalRefKey {
	return TerminalRefKey{
		application: application,
		user:        user,
		idNum:       idNum,
	}
}

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azid.RefKey = _TerminalRefKeyZero
var _ azid.BinFieldUnmarshalable = &_TerminalRefKeyZero
var _ azid.TextUnmarshalable = &_TerminalRefKeyZero
var _ azfl.AdjunctEntityRefKey = _TerminalRefKeyZero
var _ azfl.TerminalRefKey = _TerminalRefKeyZero

var _TerminalRefKeyZero = TerminalRefKey{
	application: ApplicationRefKeyZero(),
	user:        UserRefKeyZero(),
	idNum:       TerminalIDNumZero,
}

// TerminalRefKeyZero returns
// a zero-valued instance of TerminalRefKey.
func TerminalRefKeyZero() TerminalRefKey {
	return _TerminalRefKeyZero
}

// AZRefKey is required by azid.RefKey interface.
func (TerminalRefKey) AZRefKey() {}

// AZAdjunctEntityRefKey is required
// by azfl.AdjunctEntityRefKey interface.
func (TerminalRefKey) AZAdjunctEntityRefKey() {}

// IDNum returns the scoped identifier of the entity.
func (refKey TerminalRefKey) IDNum() TerminalIDNum {
	return refKey.idNum
}

// IDNumPtr returns a pointer to a copy of the id-num if it's considered valid
// otherwise it returns nil.
func (refKey TerminalRefKey) IDNumPtr() *TerminalIDNum {
	if refKey.IsNotSound() {
		return nil
	}
	i := refKey.IDNum()
	return &i
}

// AZIDNum is required for conformance with azid.RefKey.
func (refKey TerminalRefKey) AZIDNum() azid.IDNum {
	return refKey.idNum
}

// TerminalIDNum is required for conformance
// with azfl.TerminalRefKey.
func (refKey TerminalRefKey) TerminalIDNum() azfl.TerminalIDNum {
	return refKey.idNum
}

// IsZero is required as TerminalRefKey is a value-object.
func (refKey TerminalRefKey) IsZero() bool {
	return refKey.application.IsZero() &&
		refKey.user.IsZero() &&
		refKey.idNum == TerminalIDNumZero
}

// IsSound returns true if this instance is valid independently as a ref-key.
// It doesn't tell whether it refers to a valid instance of Terminal.
func (refKey TerminalRefKey) IsSound() bool {
	return refKey.application.IsSound() &&
		refKey.user.IsSound() &&
		refKey.idNum.IsSound()
}

// IsNotSound returns the negation of value returned by IsSound.
func (refKey TerminalRefKey) IsNotSound() bool {
	return !refKey.IsSound()
}

// Equals is required for conformance with azfl.AdjunctEntityRefKey.
func (refKey TerminalRefKey) Equals(other interface{}) bool {
	if x, ok := other.(TerminalRefKey); ok {
		return refKey.application.EqualsApplicationRefKey(x.application) &&
			refKey.user.EqualsUserRefKey(x.user) &&
			refKey.idNum == x.idNum
	}
	if x, _ := other.(*TerminalRefKey); x != nil {
		return refKey.application.EqualsApplicationRefKey(x.application) &&
			refKey.user.EqualsUserRefKey(x.user) &&
			refKey.idNum == x.idNum
	}
	return false
}

// Equal is required for conformance with azfl.AdjunctEntityRefKey.
func (refKey TerminalRefKey) Equal(other interface{}) bool {
	return refKey.Equals(other)
}

// EqualsTerminalRefKey returns true
// if the other value has the same attributes as refKey.
func (refKey TerminalRefKey) EqualsTerminalRefKey(
	other TerminalRefKey,
) bool {
	return refKey.application.EqualsApplicationRefKey(other.application) &&
		refKey.user.EqualsUserRefKey(other.user) &&
		refKey.idNum == other.idNum
}

// AZIDBin is required for conformance
// with azid.RefKey.
func (refKey TerminalRefKey) AZIDBin() []byte {
	data, typ := refKey.AZIDBinField()
	out := []byte{typ.Byte()}
	return append(out, data...)
}

// TerminalRefKeyFromAZIDBin creates a new instance of
// TerminalRefKey from its azid-bin form.
func TerminalRefKeyFromAZIDBin(
	b []byte,
) (refKey TerminalRefKey, readLen int, err error) {
	typ, err := azid.BinDataTypeFromByte(b[0])
	if err != nil {
		return TerminalRefKeyZero(), 0,
			errors.ArgWrap("", "type parsing", err)
	}
	if typ != azid.BinDataTypeArray {
		return TerminalRefKeyZero(), 0,
			errors.Arg("", errors.EntMsg("type", "unsupported"))
	}

	refKey, readLen, err = TerminalRefKeyFromAZIDBinField(b[1:], typ)
	return refKey, readLen + 1, err
}

// AZIDBinField is required for conformance
// with azid.RefKey.
func (refKey TerminalRefKey) AZIDBinField() ([]byte, azid.BinDataType) {
	var typesBytes []byte
	var dataBytes []byte
	var fieldBytes []byte
	var fieldType azid.BinDataType

	fieldBytes, fieldType = refKey.application.AZIDBinField()
	typesBytes = append(typesBytes, fieldType.Byte())
	dataBytes = append(dataBytes, fieldBytes...)

	fieldBytes, fieldType = refKey.user.AZIDBinField()
	typesBytes = append(typesBytes, fieldType.Byte())
	dataBytes = append(dataBytes, fieldBytes...)

	fieldBytes, fieldType = refKey.idNum.AZIDBinField()
	typesBytes = append(typesBytes, fieldType.Byte())
	dataBytes = append(dataBytes, fieldBytes...)

	var out = []byte{byte(len(typesBytes))}
	out = append(out, typesBytes...)
	out = append(out, dataBytes...)
	return out, azid.BinDataTypeArray
}

// TerminalRefKeyFromAZIDBinField creates TerminalRefKey from
// its azid-bin field form.
func TerminalRefKeyFromAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (refKey TerminalRefKey, readLen int, err error) {
	if typeHint != azid.BinDataTypeArray {
		return TerminalRefKeyZero(), 0,
			errors.Arg("", errors.EntMsg("type", "unsupported"))
	}

	arrayLen := int(b[0])
	if arrayLen != _TerminalRefKeyFieldCount {
		return TerminalRefKeyZero(), 0,
			errors.Arg("", errors.EntMsg("field count", "mismatch"))
	}

	typeCursor := 1
	dataCursor := typeCursor + arrayLen

	var fieldType azid.BinDataType

	fieldType, err = azid.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return TerminalRefKeyZero(), 0,
			errors.ArgWrap("", "application ref-key type parsing", err)
	}
	typeCursor++
	applicationRefKey, readLen, err := ApplicationRefKeyFromAZIDBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return TerminalRefKeyZero(), 0,
			errors.ArgWrap("", "application ref-key data parsing", err)
	}
	dataCursor += readLen

	fieldType, err = azid.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return TerminalRefKeyZero(), 0,
			errors.ArgWrap("", "user ref-key type parsing", err)
	}
	typeCursor++
	userRefKey, readLen, err := UserRefKeyFromAZIDBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return TerminalRefKeyZero(), 0,
			errors.ArgWrap("", "user ref-key data parsing", err)
	}
	dataCursor += readLen

	fieldType, err = azid.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return TerminalRefKeyZero(), 0,
			errors.ArgWrap("", "id-num type parsing", err)
	}
	typeCursor++
	idNum, readLen, err := TerminalIDNumFromAZIDBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return TerminalRefKeyZero(), 0,
			errors.ArgWrap("", "id-num data parsing", err)
	}
	dataCursor += readLen

	return TerminalRefKey{
		application: applicationRefKey,
		user:        userRefKey,
		idNum:       idNum,
	}, dataCursor, nil
}

// UnmarshalAZIDBinField is required for conformance
// with azfl.BinFieldUnmarshalable.
func (refKey *TerminalRefKey) UnmarshalAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (readLen int, err error) {
	i, readLen, err := TerminalRefKeyFromAZIDBinField(b, typeHint)
	if err == nil {
		*refKey = i
	}
	return readLen, err
}

const _TerminalRefKeyAZIDTextPrefix = "KTx0"

// AZIDText is required for conformance
// with azid.RefKey.
func (refKey TerminalRefKey) AZIDText() string {
	if !refKey.IsSound() {
		return ""
	}

	return _TerminalRefKeyAZIDTextPrefix +
		azid.TextEncode(refKey.AZIDBin())
}

// TerminalRefKeyFromAZIDText creates a new instance of
// TerminalRefKey from its azid-text form.
func TerminalRefKeyFromAZIDText(s string) (TerminalRefKey, error) {
	if s == "" {
		return TerminalRefKeyZero(), nil
	}
	if !strings.HasPrefix(s, _TerminalRefKeyAZIDTextPrefix) {
		return TerminalRefKeyZero(),
			errors.Arg("", errors.EntMsg("prefix", "mismatch"))
	}
	s = strings.TrimPrefix(s, _TerminalRefKeyAZIDTextPrefix)
	b, err := azid.TextDecode(s)
	if err != nil {
		return TerminalRefKeyZero(),
			errors.ArgWrap("", "data parsing", err)
	}
	refKey, _, err := TerminalRefKeyFromAZIDBin(b)
	if err != nil {
		return TerminalRefKeyZero(),
			errors.ArgWrap("", "data decoding", err)
	}
	return refKey, nil
}

// UnmarshalAZIDText is required for conformance
// with azid.TextUnmarshalable.
func (refKey *TerminalRefKey) UnmarshalAZIDText(s string) error {
	r, err := TerminalRefKeyFromAZIDText(s)
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalText is for compatibility with Go's encoding.TextMarshaler
func (refKey TerminalRefKey) MarshalText() ([]byte, error) {
	return []byte(refKey.AZIDText()), nil
}

// UnmarshalText is for conformance with Go's encoding.TextUnmarshaler
func (refKey *TerminalRefKey) UnmarshalText(b []byte) error {
	r, err := TerminalRefKeyFromAZIDText(string(b))
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalJSON makes this type JSON-marshalable.
func (refKey TerminalRefKey) MarshalJSON() ([]byte, error) {
	// We assume that there's no symbols in azid-text
	return []byte("\"" + refKey.AZIDText() + "\""), nil
}

// UnmarshalJSON parses a JSON value.
func (refKey *TerminalRefKey) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		*refKey = TerminalRefKeyZero()
		return nil
	}
	i, err := TerminalRefKeyFromAZIDText(s)
	if err == nil {
		*refKey = i
	}
	return err
}

// Application returns instance's Application value.
func (refKey TerminalRefKey) Application() ApplicationRefKey {
	return refKey.application
}

// ApplicationPtr returns a pointer to a copy of
// ApplicationRefKey if it's considered valid.
func (refKey TerminalRefKey) ApplicationPtr() *ApplicationRefKey {
	if refKey.application.IsSound() {
		rk := refKey.application
		return &rk
	}
	return nil
}

// WithApplication returns a copy
// of TerminalRefKey
// with its application attribute set to the provided value.
func (refKey TerminalRefKey) WithApplication(
	application ApplicationRefKey,
) TerminalRefKey {
	return TerminalRefKey{
		application: application,
		user:        refKey.user,
		idNum:       refKey.idNum,
	}
}

// User returns instance's User value.
func (refKey TerminalRefKey) User() UserRefKey {
	return refKey.user
}

// UserPtr returns a pointer to a copy of
// UserRefKey if it's considered valid.
func (refKey TerminalRefKey) UserPtr() *UserRefKey {
	if refKey.user.IsSound() {
		rk := refKey.user
		return &rk
	}
	return nil
}

// WithUser returns a copy
// of TerminalRefKey
// with its user attribute set to the provided value.
func (refKey TerminalRefKey) WithUser(
	user UserRefKey,
) TerminalRefKey {
	return TerminalRefKey{
		application: refKey.application,
		user:        user,
		idNum:       refKey.idNum,
	}
}

// TerminalRefKeyError defines an interface for all
// TerminalRefKey-related errors.
type TerminalRefKeyError interface {
	error
	TerminalRefKeyError()
}

//endregion
