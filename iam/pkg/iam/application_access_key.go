package iam

import (
	"encoding/binary"
	"strings"

	azcore "github.com/alloyzeus/go-azcore/azcore"
	azer "github.com/alloyzeus/go-azcore/azcore/azer"
	errors "github.com/alloyzeus/go-azcore/azcore/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the azcore package it is being compiled against.
// A compilation error at this line likely means your copy of the
// azcore package needs to be updated.
var _ = azcore.AZCorePackageIsVersion1

// Reference imports to suppress errors if they are not otherwise used.
var _ = azer.BinDataTypeUnspecified
var _ = strings.Compare

// Adjunct-entity ApplicationAccessKey of Application.

//region ID

// ApplicationAccessKeyID is a scoped identifier
// used to identify an instance of adjunct entity ApplicationAccessKey
// scoped within its host entity(s).
type ApplicationAccessKeyID int64

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azcore.EID = ApplicationAccessKeyIDZero
var _ azcore.AdjunctEntityID = ApplicationAccessKeyIDZero
var _ azer.BinFieldUnmarshalable = &_ApplicationAccessKeyIDZeroVar

// _ApplicationAccessKeyIDSignificantBitsMask is used to
// extract significant bits from an instance of ApplicationAccessKeyID.
const _ApplicationAccessKeyIDSignificantBitsMask uint64 = 0b11111111_11111111_11111111_11111111_11111111_11111111_11111111

// ApplicationAccessKeyIDZero is the zero value for ApplicationAccessKeyID.
const ApplicationAccessKeyIDZero = ApplicationAccessKeyID(0)

// _ApplicationAccessKeyIDZeroVar is used for testing
// pointer-based interfaces conformance.
var _ApplicationAccessKeyIDZeroVar = ApplicationAccessKeyIDZero

// ApplicationAccessKeyIDFromPrimitiveValue creates an instance
// of ApplicationAccessKeyID from its primitive value.
func ApplicationAccessKeyIDFromPrimitiveValue(v int64) ApplicationAccessKeyID {
	return ApplicationAccessKeyID(v)
}

// ApplicationAccessKeyIDFromAZERBinField creates ApplicationAccessKeyID from
// its azer-bin form.
func ApplicationAccessKeyIDFromAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (id ApplicationAccessKeyID, readLen int, err error) {
	if typeHint != azer.BinDataTypeUnspecified && typeHint != azer.BinDataTypeInt64 {
		return ApplicationAccessKeyID(0), 0,
			errors.ArgMsg("typeHint", "unsupported")
	}
	i := binary.BigEndian.Uint64(b)
	return ApplicationAccessKeyID(i), 8, nil
}

// PrimitiveValue returns the ID in its primitive type. Prefer to use
// this method instead of casting directly.
func (id ApplicationAccessKeyID) PrimitiveValue() int64 {
	return int64(id)
}

// AZEID is required
// for conformance with azcore.EID.
func (ApplicationAccessKeyID) AZEID() {}

// AZAdjunctEntityID is required
// for conformance with azcore.AdjunctEntityID.
func (ApplicationAccessKeyID) AZAdjunctEntityID() {}

// IsZero is required as ApplicationAccessKeyID is a value-object.
func (id ApplicationAccessKeyID) IsZero() bool {
	return id == ApplicationAccessKeyIDZero
}

// IsValid returns true if this instance is valid independently as an ID.
// It doesn't tell whether it refers to a valid instance of ApplicationAccessKey.
func (id ApplicationAccessKeyID) IsValid() bool {
	return int64(id) > 0 &&
		(uint64(id)&_ApplicationAccessKeyIDSignificantBitsMask) != 0
}

// AZERBinField is required for conformance
// with azcore.EID.
func (id ApplicationAccessKeyID) AZERBinField() ([]byte, azer.BinDataType) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(id))
	return b, azer.BinDataTypeInt64
}

// UnmarshalAZERBinField is required for conformance
// with azer.BinFieldUnmarshalable.
func (id *ApplicationAccessKeyID) UnmarshalAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (readLen int, err error) {
	i, readLen, err := ApplicationAccessKeyIDFromAZERBinField(b, typeHint)
	if err == nil {
		*id = i
	}
	return readLen, err
}

// Equals is required as ApplicationAccessKeyID is a value-object.
//
// Use EqualsApplicationAccessKeyID method if the other value
// has the same type.
func (id ApplicationAccessKeyID) Equals(other interface{}) bool {
	if x, ok := other.(ApplicationAccessKeyID); ok {
		return x == id
	}
	if x, _ := other.(*ApplicationAccessKeyID); x != nil {
		return *x == id
	}
	return false
}

// Equal is a wrapper for Equals method. It is required for
// compatibility with github.com/google/go-cmp
func (id ApplicationAccessKeyID) Equal(other interface{}) bool {
	return id.Equals(other)
}

// EqualsApplicationAccessKeyID determines if the other instance
// is equal to this instance.
func (id ApplicationAccessKeyID) EqualsApplicationAccessKeyID(
	other ApplicationAccessKeyID,
) bool {
	return id == other
}

//endregion

//region RefKey

// ApplicationAccessKeyRefKey is used to identify
// an instance of adjunct entity ApplicationAccessKey system-wide.
type ApplicationAccessKeyRefKey struct {
	application ApplicationRefKey
	id          ApplicationAccessKeyID
}

// The total number of fields in the struct.
const _ApplicationAccessKeyRefKeyFieldCount = 1 + 1

// NewApplicationAccessKeyRefKey returns a new instance
// of ApplicationAccessKeyRefKey with the provided attribute values.
func NewApplicationAccessKeyRefKey(
	application ApplicationRefKey,
	id ApplicationAccessKeyID,
) ApplicationAccessKeyRefKey {
	return ApplicationAccessKeyRefKey{
		application: application,
		id:          id,
	}
}

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azcore.RefKey = _ApplicationAccessKeyRefKeyZero
var _ azcore.AdjunctEntityRefKey = _ApplicationAccessKeyRefKeyZero

var _ApplicationAccessKeyRefKeyZero = ApplicationAccessKeyRefKey{
	application: ApplicationRefKeyZero(),
	id:          ApplicationAccessKeyIDZero,
}

// ApplicationAccessKeyRefKeyZero returns
// a zero-valued instance of ApplicationAccessKeyRefKey.
func ApplicationAccessKeyRefKeyZero() ApplicationAccessKeyRefKey {
	return _ApplicationAccessKeyRefKeyZero
}

// AZRefKey is required by azcore.RefKey interface.
func (ApplicationAccessKeyRefKey) AZRefKey() {}

// AZAdjunctEntityRefKey is required
// by azcore.AdjunctEntityRefKey interface.
func (ApplicationAccessKeyRefKey) AZAdjunctEntityRefKey() {}

// ID is required for conformance with azcore.RefKey.
func (refKey ApplicationAccessKeyRefKey) ID() azcore.EID {
	return refKey.id
}

// IsZero is required as ApplicationAccessKeyRefKey is a value-object.
func (refKey ApplicationAccessKeyRefKey) IsZero() bool {
	return refKey.application.IsZero() &&
		refKey.id == ApplicationAccessKeyIDZero
}

// IsValid returns true if this instance is valid independently as a ref-key.
// It doesn't tell whether it refers to a valid instance of ApplicationAccessKey.
func (refKey ApplicationAccessKeyRefKey) IsValid() bool {
	return refKey.application.IsValid() &&
		refKey.id.IsValid()
}

// Equals is required for conformance with azcore.AdjunctEntityRefKey.
func (refKey ApplicationAccessKeyRefKey) Equals(other interface{}) bool {
	if x, ok := other.(ApplicationAccessKeyRefKey); ok {
		return refKey.application.EqualsApplicationRefKey(x.application) &&
			refKey.id == x.id
	}
	if x, _ := other.(*ApplicationAccessKeyRefKey); x != nil {
		return refKey.application.EqualsApplicationRefKey(x.application) &&
			refKey.id == x.id
	}
	return false
}

// Equal is required for conformance with azcore.AdjunctEntityRefKey.
func (refKey ApplicationAccessKeyRefKey) Equal(other interface{}) bool {
	return refKey.Equals(other)
}

// EqualsApplicationAccessKeyRefKey returns true
// if the other value has the same attributes as refKey.
func (refKey ApplicationAccessKeyRefKey) EqualsApplicationAccessKeyRefKey(
	other ApplicationAccessKeyRefKey,
) bool {
	return refKey.application.EqualsApplicationRefKey(other.application) &&
		refKey.id == other.id
}

// AZERBin is required for conformance
// with azcore.RefKey.
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
// with azcore.RefKey.
func (refKey ApplicationAccessKeyRefKey) AZERBinField() ([]byte, azer.BinDataType) {
	var typesBytes []byte
	var dataBytes []byte
	var fieldBytes []byte
	var fieldType azer.BinDataType

	fieldBytes, fieldType = refKey.application.AZERBinField()
	typesBytes = append(typesBytes, fieldType.Byte())
	dataBytes = append(dataBytes, fieldBytes...)

	fieldBytes, fieldType = refKey.id.AZERBinField()
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
			errors.ArgWrap("", "id type parsing", err)
	}
	typeCursor++
	id, readLen, err := ApplicationAccessKeyIDFromAZERBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(), 0,
			errors.ArgWrap("", "id data parsing", err)
	}
	dataCursor += readLen

	return ApplicationAccessKeyRefKey{
		application: applicationRefKey,
		id:          id,
	}, dataCursor, nil
}

// UnmarshalAZERBinField is required for conformance
// with azcore.BinFieldUnmarshalable.
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
// with azcore.RefKey.
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

// MarshalJSON makes this type JSON-marshalable.
func (refKey ApplicationAccessKeyRefKey) MarshalJSON() ([]byte, error) {
	// We assume that there's no symbols in AZRS
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

// WithApplication returns a copy
// of ApplicationAccessKeyRefKey
// with its application attribute set to the provided value.
func (refKey ApplicationAccessKeyRefKey) WithApplication(
	application ApplicationRefKey,
) ApplicationAccessKeyRefKey {
	return ApplicationAccessKeyRefKey{
		application: application,
	}
}

// ApplicationAccessKeyRefKeyError defines an interface for all
// ApplicationAccessKeyRefKey-related errors.
type ApplicationAccessKeyRefKeyError interface {
	error
	ApplicationAccessKeyRefKeyError()
}

//endregion
