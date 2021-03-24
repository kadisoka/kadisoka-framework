package iam

import (
	"encoding/binary"
	"strings"

	azfl "github.com/alloyzeus/go-azfl/azfl"
	azer "github.com/alloyzeus/go-azfl/azfl/azer"
	errors "github.com/alloyzeus/go-azfl/azfl/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the azfl package it is being compiled against.
// A compilation error at this line likely means your copy of the
// azfl package needs to be updated.
var _ = azfl.AZCorePackageIsVersion1

// Reference imports to suppress errors if they are not otherwise used.
var _ = azer.BinDataTypeUnspecified
var _ = strings.Compare

// Adjunct-entity ApplicationAccessKey of Application.

//region IDNum

// ApplicationAccessKeyIDNum is a scoped identifier
// used to identify an instance of adjunct entity ApplicationAccessKey
// scoped within its host entity(s).
type ApplicationAccessKeyIDNum int64

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azfl.IDNum = ApplicationAccessKeyIDNumZero
var _ azfl.AdjunctEntityID = ApplicationAccessKeyIDNumZero
var _ azer.BinFieldUnmarshalable = &_ApplicationAccessKeyIDNumZeroVar

// ApplicationAccessKeyIDNumSignificantBitsMask is used to
// extract significant bits from an instance of ApplicationAccessKeyIDNum.
const ApplicationAccessKeyIDNumSignificantBitsMask uint64 = 0b11111111_11111111_11111111_11111111_11111111_11111111_11111111

// ApplicationAccessKeyIDNumZero is the zero value for ApplicationAccessKeyIDNum.
const ApplicationAccessKeyIDNumZero = ApplicationAccessKeyIDNum(0)

// _ApplicationAccessKeyIDNumZeroVar is used for testing
// pointer-based interfaces conformance.
var _ApplicationAccessKeyIDNumZeroVar = ApplicationAccessKeyIDNumZero

// ApplicationAccessKeyIDNumFromPrimitiveValue creates an instance
// of ApplicationAccessKeyIDNum from its primitive value.
func ApplicationAccessKeyIDNumFromPrimitiveValue(v int64) ApplicationAccessKeyIDNum {
	return ApplicationAccessKeyIDNum(v)
}

// ApplicationAccessKeyIDNumFromAZERBinField creates ApplicationAccessKeyIDNum from
// its azer-bin form.
func ApplicationAccessKeyIDNumFromAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (idNum ApplicationAccessKeyIDNum, readLen int, err error) {
	if typeHint != azer.BinDataTypeUnspecified && typeHint != azer.BinDataTypeInt64 {
		return ApplicationAccessKeyIDNum(0), 0,
			errors.ArgMsg("typeHint", "unsupported")
	}
	i := binary.BigEndian.Uint64(b)
	return ApplicationAccessKeyIDNum(i), 8, nil
}

// PrimitiveValue returns the value in its primitive type. Prefer to use
// this method instead of casting directly.
func (idNum ApplicationAccessKeyIDNum) PrimitiveValue() int64 {
	return int64(idNum)
}

// AZIDNum is required
// for conformance with azfl.IDNum.
func (ApplicationAccessKeyIDNum) AZIDNum() {}

// AZAdjunctEntityID is required
// for conformance with azfl.AdjunctEntityID.
func (ApplicationAccessKeyIDNum) AZAdjunctEntityID() {}

// IsZero is required as ApplicationAccessKeyIDNum is a value-object.
func (idNum ApplicationAccessKeyIDNum) IsZero() bool {
	return idNum == ApplicationAccessKeyIDNumZero
}

// IsValid returns true if this instance is valid independently
// as an ApplicationAccessKeyIDNum. It doesn't tell whether it refers to
// a valid instance of ApplicationAccessKey.
func (idNum ApplicationAccessKeyIDNum) IsValid() bool {
	return int64(idNum) > 0 &&
		(uint64(idNum)&ApplicationAccessKeyIDNumSignificantBitsMask) != 0
}

// IsNotValid returns the negation of value returned by IsValid().
func (idNum ApplicationAccessKeyIDNum) IsNotValid() bool {
	return !idNum.IsValid()
}

// AZERBinField is required for conformance
// with azfl.IDNum.
func (idNum ApplicationAccessKeyIDNum) AZERBinField() ([]byte, azer.BinDataType) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(idNum))
	return b, azer.BinDataTypeInt64
}

// UnmarshalAZERBinField is required for conformance
// with azer.BinFieldUnmarshalable.
func (idNum *ApplicationAccessKeyIDNum) UnmarshalAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (readLen int, err error) {
	i, readLen, err := ApplicationAccessKeyIDNumFromAZERBinField(b, typeHint)
	if err == nil {
		*idNum = i
	}
	return readLen, err
}

// Equals is required as ApplicationAccessKeyIDNum is a value-object.
//
// Use EqualsApplicationAccessKeyIDNum method if the other value
// has the same type.
func (idNum ApplicationAccessKeyIDNum) Equals(other interface{}) bool {
	if x, ok := other.(ApplicationAccessKeyIDNum); ok {
		return x == idNum
	}
	if x, _ := other.(*ApplicationAccessKeyIDNum); x != nil {
		return *x == idNum
	}
	return false
}

// Equal is a wrapper for Equals method. It is required for
// compatibility with github.com/google/go-cmp
func (idNum ApplicationAccessKeyIDNum) Equal(other interface{}) bool {
	return idNum.Equals(other)
}

// EqualsApplicationAccessKeyIDNum determines if the other instance
// is equal to this instance.
func (idNum ApplicationAccessKeyIDNum) EqualsApplicationAccessKeyIDNum(
	other ApplicationAccessKeyIDNum,
) bool {
	return idNum == other
}

//endregion

//region RefKey

// ApplicationAccessKeyRefKey is used to identify
// an instance of adjunct entity ApplicationAccessKey system-wide.
type ApplicationAccessKeyRefKey struct {
	application ApplicationRefKey
	idNum       ApplicationAccessKeyIDNum
}

// The total number of fields in the struct.
const _ApplicationAccessKeyRefKeyFieldCount = 1 + 1

// NewApplicationAccessKeyRefKey returns a new instance
// of ApplicationAccessKeyRefKey with the provided attribute values.
func NewApplicationAccessKeyRefKey(
	application ApplicationRefKey,
	idNum ApplicationAccessKeyIDNum,
) ApplicationAccessKeyRefKey {
	return ApplicationAccessKeyRefKey{
		application: application,
		idNum:       idNum,
	}
}

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azfl.RefKey = _ApplicationAccessKeyRefKeyZero
var _ azfl.AdjunctEntityRefKey = _ApplicationAccessKeyRefKeyZero

var _ApplicationAccessKeyRefKeyZero = ApplicationAccessKeyRefKey{
	application: ApplicationRefKeyZero(),
	idNum:       ApplicationAccessKeyIDNumZero,
}

// ApplicationAccessKeyRefKeyZero returns
// a zero-valued instance of ApplicationAccessKeyRefKey.
func ApplicationAccessKeyRefKeyZero() ApplicationAccessKeyRefKey {
	return _ApplicationAccessKeyRefKeyZero
}

// AZRefKey is required by azfl.RefKey interface.
func (ApplicationAccessKeyRefKey) AZRefKey() {}

// AZAdjunctEntityRefKey is required
// by azfl.AdjunctEntityRefKey interface.
func (ApplicationAccessKeyRefKey) AZAdjunctEntityRefKey() {}

// IDNum returns the scoped identifier of the entity.
func (refKey ApplicationAccessKeyRefKey) IDNum() ApplicationAccessKeyIDNum {
	return refKey.idNum
}

// IDNumPtr returns a pointer to a copy of the IDNum if it's considered valid
// otherwise it returns nil.
func (refKey ApplicationAccessKeyRefKey) IDNumPtr() *ApplicationAccessKeyIDNum {
	if refKey.IsNotValid() {
		return nil
	}
	i := refKey.IDNum()
	return &i
}

// AZIDNum is required for conformance with azfl.RefKey.
func (refKey ApplicationAccessKeyRefKey) AZIDNum() azfl.IDNum {
	return refKey.idNum
}

// IsZero is required as ApplicationAccessKeyRefKey is a value-object.
func (refKey ApplicationAccessKeyRefKey) IsZero() bool {
	return refKey.application.IsZero() &&
		refKey.idNum == ApplicationAccessKeyIDNumZero
}

// IsValid returns true if this instance is valid independently as a ref-key.
// It doesn't tell whether it refers to a valid instance of ApplicationAccessKey.
func (refKey ApplicationAccessKeyRefKey) IsValid() bool {
	return refKey.application.IsValid() &&
		refKey.idNum.IsValid()
}

// IsNotValid returns the negation of value returned by IsValid().
func (refKey ApplicationAccessKeyRefKey) IsNotValid() bool {
	return !refKey.IsValid()
}

// Equals is required for conformance with azfl.AdjunctEntityRefKey.
func (refKey ApplicationAccessKeyRefKey) Equals(other interface{}) bool {
	if x, ok := other.(ApplicationAccessKeyRefKey); ok {
		return refKey.application.EqualsApplicationRefKey(x.application) &&
			refKey.idNum == x.idNum
	}
	if x, _ := other.(*ApplicationAccessKeyRefKey); x != nil {
		return refKey.application.EqualsApplicationRefKey(x.application) &&
			refKey.idNum == x.idNum
	}
	return false
}

// Equal is required for conformance with azfl.AdjunctEntityRefKey.
func (refKey ApplicationAccessKeyRefKey) Equal(other interface{}) bool {
	return refKey.Equals(other)
}

// EqualsApplicationAccessKeyRefKey returns true
// if the other value has the same attributes as refKey.
func (refKey ApplicationAccessKeyRefKey) EqualsApplicationAccessKeyRefKey(
	other ApplicationAccessKeyRefKey,
) bool {
	return refKey.application.EqualsApplicationRefKey(other.application) &&
		refKey.idNum == other.idNum
}

// AZERBin is required for conformance
// with azfl.RefKey.
func (refKey ApplicationAccessKeyRefKey) AZERBin() []byte {
	data, typ := refKey.AZERBinField()
	out := []byte{typ.Byte()}
	return append(out, data...)
}

// ApplicationAccessKeyRefKeyFromAZERBin creates a new instance of
// ApplicationAccessKeyRefKey from its azer-bin form.
func ApplicationAccessKeyRefKeyFromAZERBin(
	b []byte,
) (refKey ApplicationAccessKeyRefKey, readLen int, err error) {
	typ, err := azer.BinDataTypeFromByte(b[0])
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(), 0,
			errors.ArgWrap("", "type parsing", err)
	}
	if typ != azer.BinDataTypeArray {
		return ApplicationAccessKeyRefKeyZero(), 0,
			errors.Arg("", errors.EntMsg("type", "unsupported"))
	}

	refKey, readLen, err = ApplicationAccessKeyRefKeyFromAZERBinField(b[1:], typ)
	return refKey, readLen + 1, err
}

// AZERBinField is required for conformance
// with azfl.RefKey.
func (refKey ApplicationAccessKeyRefKey) AZERBinField() ([]byte, azer.BinDataType) {
	var typesBytes []byte
	var dataBytes []byte
	var fieldBytes []byte
	var fieldType azer.BinDataType

	fieldBytes, fieldType = refKey.application.AZERBinField()
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

// ApplicationAccessKeyRefKeyFromAZERBinField creates ApplicationAccessKeyRefKey from
// its azer-bin field form.
func ApplicationAccessKeyRefKeyFromAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (refKey ApplicationAccessKeyRefKey, readLen int, err error) {
	if typeHint != azer.BinDataTypeArray {
		return ApplicationAccessKeyRefKeyZero(), 0,
			errors.Arg("", errors.EntMsg("type", "unsupported"))
	}

	arrayLen := int(b[0])
	if arrayLen != _ApplicationAccessKeyRefKeyFieldCount {
		return ApplicationAccessKeyRefKeyZero(), 0,
			errors.Arg("", errors.EntMsg("field count", "mismatch"))
	}

	typeCursor := 1
	dataCursor := typeCursor + arrayLen

	var fieldType azer.BinDataType

	fieldType, err = azer.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(), 0,
			errors.ArgWrap("", "application ref-key type parsing", err)
	}
	typeCursor++
	applicationRefKey, readLen, err := ApplicationRefKeyFromAZERBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(), 0,
			errors.ArgWrap("", "application ref-key data parsing", err)
	}
	dataCursor += readLen

	fieldType, err = azer.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(), 0,
			errors.ArgWrap("", "idnum type parsing", err)
	}
	typeCursor++
	idNum, readLen, err := ApplicationAccessKeyIDNumFromAZERBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(), 0,
			errors.ArgWrap("", "idnum data parsing", err)
	}
	dataCursor += readLen

	return ApplicationAccessKeyRefKey{
		application: applicationRefKey,
		idNum:       idNum,
	}, dataCursor, nil
}

// UnmarshalAZERBinField is required for conformance
// with azfl.BinFieldUnmarshalable.
func (refKey *ApplicationAccessKeyRefKey) UnmarshalAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (readLen int, err error) {
	i, readLen, err := ApplicationAccessKeyRefKeyFromAZERBinField(b, typeHint)
	if err == nil {
		*refKey = i
	}
	return readLen, err
}

const _ApplicationAccessKeyRefKeyAZERTextPrefix = "KAK0"

// AZERText is required for conformance
// with azfl.RefKey.
func (refKey ApplicationAccessKeyRefKey) AZERText() string {
	if !refKey.IsValid() {
		return ""
	}

	return _ApplicationAccessKeyRefKeyAZERTextPrefix +
		azer.TextEncode(refKey.AZERBin())
}

// ApplicationAccessKeyRefKeyFromAZERText creates a new instance of
// ApplicationAccessKeyRefKey from its azer-text form.
func ApplicationAccessKeyRefKeyFromAZERText(s string) (ApplicationAccessKeyRefKey, error) {
	if s == "" {
		return ApplicationAccessKeyRefKeyZero(), nil
	}
	if !strings.HasPrefix(s, _ApplicationAccessKeyRefKeyAZERTextPrefix) {
		return ApplicationAccessKeyRefKeyZero(),
			errors.Arg("", errors.EntMsg("prefix", "mismatch"))
	}
	s = strings.TrimPrefix(s, _ApplicationAccessKeyRefKeyAZERTextPrefix)
	b, err := azer.TextDecode(s)
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(),
			errors.ArgWrap("", "data parsing", err)
	}
	refKey, _, err := ApplicationAccessKeyRefKeyFromAZERBin(b)
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(),
			errors.ArgWrap("", "data decoding", err)
	}
	return refKey, nil
}

// UnmarshalAZERText is required for conformance
// with azer.TextUnmarshalable.
func (refKey *ApplicationAccessKeyRefKey) UnmarshalAZERText(s string) error {
	r, err := ApplicationAccessKeyRefKeyFromAZERText(s)
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalText is for compatibility with Go's encoding.TextMarshaler
func (refKey ApplicationAccessKeyRefKey) MarshalText() ([]byte, error) {
	return []byte(refKey.AZERText()), nil
}

// UnmarshalText is for conformance with Go's encoding.TextUnmarshaler
func (refKey *ApplicationAccessKeyRefKey) UnmarshalText(b []byte) error {
	r, err := ApplicationAccessKeyRefKeyFromAZERText(string(b))
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalJSON makes this type JSON-marshalable.
func (refKey ApplicationAccessKeyRefKey) MarshalJSON() ([]byte, error) {
	// We assume that there's no symbols in azer-text
	return []byte("\"" + refKey.AZERText() + "\""), nil
}

// UnmarshalJSON parses a JSON value.
func (refKey *ApplicationAccessKeyRefKey) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		*refKey = ApplicationAccessKeyRefKeyZero()
		return nil
	}
	i, err := ApplicationAccessKeyRefKeyFromAZERText(s)
	if err == nil {
		*refKey = i
	}
	return err
}

// Application returns instance's Application value.
func (refKey ApplicationAccessKeyRefKey) Application() ApplicationRefKey {
	return refKey.application
}

// ApplicationPtr returns a pointer to a copy of
// ApplicationRefKey if it's considered valid.
func (refKey ApplicationAccessKeyRefKey) ApplicationPtr() *ApplicationRefKey {
	if refKey.application.IsValid() {
		rk := refKey.application
		return &rk
	}
	return nil
}

// WithApplication returns a copy
// of ApplicationAccessKeyRefKey
// with its application attribute set to the provided value.
func (refKey ApplicationAccessKeyRefKey) WithApplication(
	application ApplicationRefKey,
) ApplicationAccessKeyRefKey {
	return ApplicationAccessKeyRefKey{
		application: application,
		idNum:       refKey.idNum,
	}
}

// ApplicationAccessKeyRefKeyError defines an interface for all
// ApplicationAccessKeyRefKey-related errors.
type ApplicationAccessKeyRefKeyError interface {
	error
	ApplicationAccessKeyRefKeyError()
}

//endregion
