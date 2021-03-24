package iam

import (
	"encoding/binary"
	"strings"

	azfl "github.com/alloyzeus/go-azfl/azfl"
	azer "github.com/alloyzeus/go-azfl/azfl/azer"
	"github.com/alloyzeus/go-azfl/azfl/errors"
)

//region IDNum

// TerminalIDNum is a scoped identifier
// used to identify an instance of adjunct entity Terminal
// scoped within its host entity(s).
type TerminalIDNum int64

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azfl.IDNum = TerminalIDNumZero
var _ azfl.AdjunctEntityID = TerminalIDNumZero
var _ azer.BinFieldUnmarshalable = &_TerminalIDNumZeroVar
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

// TerminalIDNumFromAZERBinField creates TerminalIDNum from
// its azer-bin form.
func TerminalIDNumFromAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (idNum TerminalIDNum, readLen int, err error) {
	if typeHint != azer.BinDataTypeUnspecified && typeHint != azer.BinDataTypeInt64 {
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
// for conformance with azfl.IDNum.
func (TerminalIDNum) AZIDNum() {}

// AZAdjunctEntityID is required
// for conformance with azfl.AdjunctEntityID.
func (TerminalIDNum) AZAdjunctEntityID() {}

// AZTerminalIDNum is required for conformance
// with azfl.TerminalIDNum.
func (TerminalIDNum) AZTerminalIDNum() {}

// IsZero is required as TerminalIDNum is a value-object.
func (idNum TerminalIDNum) IsZero() bool {
	return idNum == TerminalIDNumZero
}

// IsValid returns true if this instance is valid independently
// as an TerminalIDNum. It doesn't tell whether it refers to
// a valid instance of Terminal.
func (idNum TerminalIDNum) IsValid() bool {
	return int64(idNum) > 0 &&
		(uint64(idNum)&TerminalIDNumSignificantBitsMask) != 0
}

// IsNotValid returns the negation of value returned by IsValid().
func (idNum TerminalIDNum) IsNotValid() bool {
	return !idNum.IsValid()
}

// AZERBinField is required for conformance
// with azfl.IDNum.
func (idNum TerminalIDNum) AZERBinField() ([]byte, azer.BinDataType) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(idNum))
	return b, azer.BinDataTypeInt64
}

// UnmarshalAZERBinField is required for conformance
// with azer.BinFieldUnmarshalable.
func (idNum *TerminalIDNum) UnmarshalAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (readLen int, err error) {
	i, readLen, err := TerminalIDNumFromAZERBinField(b, typeHint)
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
var _ azfl.RefKey = _TerminalRefKeyZero
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

// AZRefKey is required by azfl.RefKey interface.
func (TerminalRefKey) AZRefKey() {}

// AZAdjunctEntityRefKey is required
// by azfl.AdjunctEntityRefKey interface.
func (TerminalRefKey) AZAdjunctEntityRefKey() {}

// IDNum returns the scoped identifier of the entity.
func (refKey TerminalRefKey) IDNum() TerminalIDNum {
	return refKey.idNum
}

// IDNumPtr returns a pointer to a copy of the IDNum if it's considered valid
// otherwise it returns nil.
func (refKey TerminalRefKey) IDNumPtr() *TerminalIDNum {
	if refKey.IsNotValid() {
		return nil
	}
	i := refKey.IDNum()
	return &i
}

// AZIDNum is required for conformance with azfl.RefKey.
func (refKey TerminalRefKey) AZIDNum() azfl.IDNum {
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

// IsValid returns true if this instance is valid independently as a ref-key.
// It doesn't tell whether it refers to a valid instance of Terminal.
func (refKey TerminalRefKey) IsValid() bool {
	return refKey.application.IsValid() &&
		refKey.user.IsValid() &&
		refKey.idNum.IsValid()
}

// IsNotValid returns the negation of value returned by IsValid().
func (refKey TerminalRefKey) IsNotValid() bool {
	return !refKey.IsValid()
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

// AZERBin is required for conformance
// with azfl.RefKey.
func (refKey TerminalRefKey) AZERBin() []byte {
	data, typ := refKey.AZERBinField()
	out := []byte{typ.Byte()}
	return append(out, data...)
}

// TerminalRefKeyFromAZERBin creates a new instance of
// TerminalRefKey from its azer-bin form.
func TerminalRefKeyFromAZERBin(
	b []byte,
) (refKey TerminalRefKey, readLen int, err error) {
	typ, err := azer.BinDataTypeFromByte(b[0])
	if err != nil {
		return TerminalRefKeyZero(), 0,
			errors.ArgWrap("", "type parsing", err)
	}
	if typ != azer.BinDataTypeArray {
		return TerminalRefKeyZero(), 0,
			errors.Arg("", errors.EntMsg("type", "unsupported"))
	}

	refKey, readLen, err = TerminalRefKeyFromAZERBinField(b[1:], typ)
	return refKey, readLen + 1, err
}

// AZERBinField is required for conformance
// with azfl.RefKey.
func (refKey TerminalRefKey) AZERBinField() ([]byte, azer.BinDataType) {
	var typesBytes []byte
	var dataBytes []byte
	var fieldBytes []byte
	var fieldType azer.BinDataType

	fieldBytes, fieldType = refKey.application.AZERBinField()
	typesBytes = append(typesBytes, fieldType.Byte())
	dataBytes = append(dataBytes, fieldBytes...)

	fieldBytes, fieldType = refKey.user.AZERBinField()
	typesBytes = append(typesBytes, fieldType.Byte())
	dataBytes = append(dataBytes, fieldBytes...)

	fieldBytes, fieldType = refKey.idNum.AZERBinField()
	typesBytes = append(typesBytes, fieldType.Byte())
	dataBytes = append(dataBytes, fieldBytes...)

	var out = []byte{byte(len(typesBytes))}
	out = append(out, typesBytes...)
	out = append(out, dataBytes...)
	return out, azer.BinDataTypeArray
}

// TerminalRefKeyFromAZERBinField creates TerminalRefKey from
// its azer-bin field form.
func TerminalRefKeyFromAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (refKey TerminalRefKey, readLen int, err error) {
	if typeHint != azer.BinDataTypeArray {
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

	var fieldType azer.BinDataType

	fieldType, err = azer.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return TerminalRefKeyZero(), 0,
			errors.ArgWrap("", "application ref-key type parsing", err)
	}
	typeCursor++
	applicationRefKey, readLen, err := ApplicationRefKeyFromAZERBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return TerminalRefKeyZero(), 0,
			errors.ArgWrap("", "application ref-key data parsing", err)
	}
	dataCursor += readLen

	fieldType, err = azer.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return TerminalRefKeyZero(), 0,
			errors.ArgWrap("", "user ref-key type parsing", err)
	}
	typeCursor++
	userRefKey, readLen, err := UserRefKeyFromAZERBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return TerminalRefKeyZero(), 0,
			errors.ArgWrap("", "user ref-key data parsing", err)
	}
	dataCursor += readLen

	fieldType, err = azer.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return TerminalRefKeyZero(), 0,
			errors.ArgWrap("", "idnum type parsing", err)
	}
	typeCursor++
	idNum, readLen, err := TerminalIDNumFromAZERBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return TerminalRefKeyZero(), 0,
			errors.ArgWrap("", "idnum data parsing", err)
	}
	dataCursor += readLen

	return TerminalRefKey{
		application: applicationRefKey,
		user:        userRefKey,
		idNum:       idNum,
	}, dataCursor, nil
}

// UnmarshalAZERBinField is required for conformance
// with azfl.BinFieldUnmarshalable.
func (refKey *TerminalRefKey) UnmarshalAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (readLen int, err error) {
	i, readLen, err := TerminalRefKeyFromAZERBinField(b, typeHint)
	if err == nil {
		*refKey = i
	}
	return readLen, err
}

const _TerminalRefKeyAZERTextPrefix = "KTx0"

// AZERText is required for conformance
// with azfl.RefKey.
func (refKey TerminalRefKey) AZERText() string {
	if !refKey.IsValid() {
		return ""
	}

	return _TerminalRefKeyAZERTextPrefix +
		azer.TextEncode(refKey.AZERBin())
}

// TerminalRefKeyFromAZERText creates a new instance of
// TerminalRefKey from its azer-text form.
func TerminalRefKeyFromAZERText(s string) (TerminalRefKey, error) {
	if s == "" {
		return TerminalRefKeyZero(), nil
	}
	if !strings.HasPrefix(s, _TerminalRefKeyAZERTextPrefix) {
		return TerminalRefKeyZero(),
			errors.Arg("", errors.EntMsg("prefix", "mismatch"))
	}
	s = strings.TrimPrefix(s, _TerminalRefKeyAZERTextPrefix)
	b, err := azer.TextDecode(s)
	if err != nil {
		return TerminalRefKeyZero(),
			errors.ArgWrap("", "data parsing", err)
	}
	refKey, _, err := TerminalRefKeyFromAZERBin(b)
	if err != nil {
		return TerminalRefKeyZero(),
			errors.ArgWrap("", "data decoding", err)
	}
	return refKey, nil
}

// UnmarshalAZERText is required for conformance
// with azer.TextUnmarshalable.
func (refKey *TerminalRefKey) UnmarshalAZERText(s string) error {
	r, err := TerminalRefKeyFromAZERText(s)
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalText is for compatibility with Go's encoding.TextMarshaler
func (refKey TerminalRefKey) MarshalText() ([]byte, error) {
	return []byte(refKey.AZERText()), nil
}

// UnmarshalText is for conformance with Go's encoding.TextUnmarshaler
func (refKey *TerminalRefKey) UnmarshalText(b []byte) error {
	r, err := TerminalRefKeyFromAZERText(string(b))
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalJSON makes this type JSON-marshalable.
func (refKey TerminalRefKey) MarshalJSON() ([]byte, error) {
	// We assume that there's no symbols in azer-text
	return []byte("\"" + refKey.AZERText() + "\""), nil
}

// UnmarshalJSON parses a JSON value.
func (refKey *TerminalRefKey) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		*refKey = TerminalRefKeyZero()
		return nil
	}
	i, err := TerminalRefKeyFromAZERText(s)
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
	if refKey.application.IsValid() {
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
	if refKey.user.IsValid() {
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
