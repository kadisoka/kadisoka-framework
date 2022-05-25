package iam

import (
	"crypto/rand"
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
var _ = rand.Reader

// Adjunct-entity ApplicationAccessKey of Application.

//region IDNum

// ApplicationAccessKeyIDNum is a scoped identifier
// used to identify an instance of adjunct entity ApplicationAccessKey
// scoped within its host entity(s).
type ApplicationAccessKeyIDNum int64

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azid.IDNumMethods = ApplicationAccessKeyIDNumZero
var _ azid.BinFieldUnmarshalable = &_ApplicationAccessKeyIDNumZeroVar
var _ azfl.AdjunctEntityIDNumMethods = ApplicationAccessKeyIDNumZero

// ApplicationAccessKeyIDNumIdentifierBitsMask is used to
// extract identifier bits from an instance of ApplicationAccessKeyIDNum.
const ApplicationAccessKeyIDNumIdentifierBitsMask uint64 = 0b_00000000_11111111_11111111_11111111_11111111_11111111_11111111_11111111

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

// ApplicationAccessKeyIDNumFromAZIDBinField creates ApplicationAccessKeyIDNum from
// its azid-bin form.
func ApplicationAccessKeyIDNumFromAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (idNum ApplicationAccessKeyIDNum, readLen int, err error) {
	if typeHint != azid.BinDataTypeUnspecified && typeHint != azid.BinDataTypeInt64 {
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
// for conformance with azid.IDNum.
func (ApplicationAccessKeyIDNum) AZIDNum() {}

// AZAdjunctEntityIDNum is required
// for conformance with azfl.AdjunctEntityIDNum.
func (ApplicationAccessKeyIDNum) AZAdjunctEntityIDNum() {}

// IsZero is required as ApplicationAccessKeyIDNum is a value-object.
func (idNum ApplicationAccessKeyIDNum) IsZero() bool {
	return idNum == ApplicationAccessKeyIDNumZero
}

// IsStaticallyValid returns true if this instance is valid as an isolated value
// of ApplicationAccessKeyIDNum. It doesn't tell whether it refers to
// a valid instance of ApplicationAccessKey.
func (idNum ApplicationAccessKeyIDNum) IsStaticallyValid() bool {
	return int64(idNum) > 0 &&
		(uint64(idNum)&ApplicationAccessKeyIDNumIdentifierBitsMask) != 0
}

// IsNotStaticallyValid returns the negation of value returned by IsStaticallyValid.
func (idNum ApplicationAccessKeyIDNum) IsNotStaticallyValid() bool {
	return !idNum.IsStaticallyValid()
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

// AZIDBinField is required for conformance
// with azid.IDNum.
func (idNum ApplicationAccessKeyIDNum) AZIDBinField() ([]byte, azid.BinDataType) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(idNum))
	return b, azid.BinDataTypeInt64
}

// UnmarshalAZIDBinField is required for conformance
// with azid.BinFieldUnmarshalable.
func (idNum *ApplicationAccessKeyIDNum) UnmarshalAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (readLen int, err error) {
	i, readLen, err := ApplicationAccessKeyIDNumFromAZIDBinField(b, typeHint)
	if err == nil {
		*idNum = i
	}
	return readLen, err
}

// Embedded fields
const (
	ApplicationAccessKeyIDNumEmbeddedFieldsMask = 0b_00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
)

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
var _ azid.RefKey[ApplicationAccessKeyIDNum] = _ApplicationAccessKeyRefKeyZero
var _ azid.BinFieldUnmarshalable = &_ApplicationAccessKeyRefKeyZero
var _ azid.TextUnmarshalable = &_ApplicationAccessKeyRefKeyZero
var _ azfl.AdjunctEntityRefKey[ApplicationAccessKeyIDNum] = _ApplicationAccessKeyRefKeyZero

var _ApplicationAccessKeyRefKeyZero = ApplicationAccessKeyRefKey{
	application: ApplicationRefKeyZero(),
	idNum:       ApplicationAccessKeyIDNumZero,
}

// ApplicationAccessKeyRefKeyZero returns
// a zero-valued instance of ApplicationAccessKeyRefKey.
func ApplicationAccessKeyRefKeyZero() ApplicationAccessKeyRefKey {
	return _ApplicationAccessKeyRefKeyZero
}

// AZRefKey is required by azid.RefKey interface.
func (ApplicationAccessKeyRefKey) AZRefKey() {}

// AZAdjunctEntityRefKey is required
// by azfl.AdjunctEntityRefKey interface.
func (ApplicationAccessKeyRefKey) AZAdjunctEntityRefKey() {}

// IDNum returns the scoped identifier of the entity.
func (refKey ApplicationAccessKeyRefKey) IDNum() ApplicationAccessKeyIDNum {
	return refKey.idNum
}

// IDNumPtr returns a pointer to a copy of the id-num if it's considered valid
// otherwise it returns nil.
func (refKey ApplicationAccessKeyRefKey) IDNumPtr() *ApplicationAccessKeyIDNum {
	if refKey.IsNotStaticallyValid() {
		return nil
	}
	i := refKey.IDNum()
	return &i
}

// AZIDNum is required for conformance with azid.RefKey.
func (refKey ApplicationAccessKeyRefKey) AZIDNum() ApplicationAccessKeyIDNum {
	return refKey.idNum
}

// IsZero is required as ApplicationAccessKeyRefKey is a value-object.
func (refKey ApplicationAccessKeyRefKey) IsZero() bool {
	return refKey.application.IsZero() &&
		refKey.idNum == ApplicationAccessKeyIDNumZero
}

// IsStaticallyValid returns true if this instance is valid as an isolated value
// of ApplicationAccessKeyRefKey.
// It doesn't tell whether it refers to a valid instance of ApplicationAccessKey.
func (refKey ApplicationAccessKeyRefKey) IsStaticallyValid() bool {
	return refKey.application.IsStaticallyValid() &&
		refKey.idNum.IsStaticallyValid()
}

// IsNotStaticallyValid returns the negation of value returned by IsStaticallyValid.
func (refKey ApplicationAccessKeyRefKey) IsNotStaticallyValid() bool {
	return !refKey.IsStaticallyValid()
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

// AZIDBin is required for conformance
// with azid.RefKey.
func (refKey ApplicationAccessKeyRefKey) AZIDBin() []byte {
	data, typ := refKey.AZIDBinField()
	out := []byte{typ.Byte()}
	return append(out, data...)
}

// ApplicationAccessKeyRefKeyFromAZIDBin creates a new instance of
// ApplicationAccessKeyRefKey from its azid-bin form.
func ApplicationAccessKeyRefKeyFromAZIDBin(
	b []byte,
) (refKey ApplicationAccessKeyRefKey, readLen int, err error) {
	typ, err := azid.BinDataTypeFromByte(b[0])
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(), 0,
			errors.ArgWrap("", "type parsing", err)
	}
	if typ != azid.BinDataTypeArray {
		return ApplicationAccessKeyRefKeyZero(), 0,
			errors.Arg("", errors.EntMsg("type", "unsupported"))
	}

	refKey, readLen, err = ApplicationAccessKeyRefKeyFromAZIDBinField(b[1:], typ)
	return refKey, readLen + 1, err
}

// AZIDBinField is required for conformance
// with azid.RefKey.
func (refKey ApplicationAccessKeyRefKey) AZIDBinField() ([]byte, azid.BinDataType) {
	var typesBytes []byte
	var dataBytes []byte
	var fieldBytes []byte
	var fieldType azid.BinDataType

	fieldBytes, fieldType = refKey.application.AZIDBinField()
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

// ApplicationAccessKeyRefKeyFromAZIDBinField creates ApplicationAccessKeyRefKey from
// its azid-bin field form.
func ApplicationAccessKeyRefKeyFromAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (refKey ApplicationAccessKeyRefKey, readLen int, err error) {
	if typeHint != azid.BinDataTypeArray {
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

	var fieldType azid.BinDataType

	fieldType, err = azid.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(), 0,
			errors.ArgWrap("", "application ref-key type parsing", err)
	}
	typeCursor++
	applicationRefKey, readLen, err := ApplicationRefKeyFromAZIDBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(), 0,
			errors.ArgWrap("", "application ref-key data parsing", err)
	}
	dataCursor += readLen

	fieldType, err = azid.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(), 0,
			errors.ArgWrap("", "id-num type parsing", err)
	}
	typeCursor++
	idNum, readLen, err := ApplicationAccessKeyIDNumFromAZIDBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(), 0,
			errors.ArgWrap("", "id-num data parsing", err)
	}
	dataCursor += readLen

	return ApplicationAccessKeyRefKey{
		application: applicationRefKey,
		idNum:       idNum,
	}, dataCursor, nil
}

// UnmarshalAZIDBinField is required for conformance
// with azfl.BinFieldUnmarshalable.
func (refKey *ApplicationAccessKeyRefKey) UnmarshalAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (readLen int, err error) {
	i, readLen, err := ApplicationAccessKeyRefKeyFromAZIDBinField(b, typeHint)
	if err == nil {
		*refKey = i
	}
	return readLen, err
}

const _ApplicationAccessKeyRefKeyAZIDTextPrefix = "KAK0"

// AZIDText is required for conformance
// with azid.RefKey.
func (refKey ApplicationAccessKeyRefKey) AZIDText() string {
	if !refKey.IsStaticallyValid() {
		return ""
	}

	return _ApplicationAccessKeyRefKeyAZIDTextPrefix +
		azid.TextEncode(refKey.AZIDBin())
}

// ApplicationAccessKeyRefKeyFromAZIDText creates a new instance of
// ApplicationAccessKeyRefKey from its azid-text form.
func ApplicationAccessKeyRefKeyFromAZIDText(s string) (ApplicationAccessKeyRefKey, error) {
	if s == "" {
		return ApplicationAccessKeyRefKeyZero(), nil
	}
	if !strings.HasPrefix(s, _ApplicationAccessKeyRefKeyAZIDTextPrefix) {
		return ApplicationAccessKeyRefKeyZero(),
			errors.Arg("", errors.EntMsg("prefix", "mismatch"))
	}
	s = strings.TrimPrefix(s, _ApplicationAccessKeyRefKeyAZIDTextPrefix)
	b, err := azid.TextDecode(s)
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(),
			errors.ArgWrap("", "data parsing", err)
	}
	refKey, _, err := ApplicationAccessKeyRefKeyFromAZIDBin(b)
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(),
			errors.ArgWrap("", "data decoding", err)
	}
	return refKey, nil
}

// UnmarshalAZIDText is required for conformance
// with azid.TextUnmarshalable.
func (refKey *ApplicationAccessKeyRefKey) UnmarshalAZIDText(s string) error {
	r, err := ApplicationAccessKeyRefKeyFromAZIDText(s)
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalText is for compatibility with Go's encoding.TextMarshaler
func (refKey ApplicationAccessKeyRefKey) MarshalText() ([]byte, error) {
	return []byte(refKey.AZIDText()), nil
}

// UnmarshalText is for conformance with Go's encoding.TextUnmarshaler
func (refKey *ApplicationAccessKeyRefKey) UnmarshalText(b []byte) error {
	r, err := ApplicationAccessKeyRefKeyFromAZIDText(string(b))
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalJSON makes this type JSON-marshalable.
func (refKey ApplicationAccessKeyRefKey) MarshalJSON() ([]byte, error) {
	// We assume that there are no symbols in azid-text
	return []byte("\"" + refKey.AZIDText() + "\""), nil
}

// UnmarshalJSON parses a JSON value.
func (refKey *ApplicationAccessKeyRefKey) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		*refKey = ApplicationAccessKeyRefKeyZero()
		return nil
	}
	i, err := ApplicationAccessKeyRefKeyFromAZIDText(s)
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
	if refKey.application.IsStaticallyValid() {
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
