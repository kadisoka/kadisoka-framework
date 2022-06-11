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
var _ azcore.AdjunctEntityIDNumMethods = ApplicationAccessKeyIDNumZero
var _ azcore.ValueObjectAssert[ApplicationAccessKeyIDNum] = ApplicationAccessKeyIDNumZero

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

// Clone returns a copy of self.
func (idNum ApplicationAccessKeyIDNum) Clone() ApplicationAccessKeyIDNum { return idNum }

// AZIDNum is required
// for conformance with azid.IDNum.
func (ApplicationAccessKeyIDNum) AZIDNum() {}

// AZAdjunctEntityIDNum is required
// for conformance with azcore.AdjunctEntityIDNum.
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

//region ID

// ApplicationAccessKeyID is used to identify
// an instance of adjunct entity ApplicationAccessKey system-wide.
type ApplicationAccessKeyID struct {
	application ApplicationID
	idNum       ApplicationAccessKeyIDNum
}

// The total number of fields in the struct.
const _ApplicationAccessKeyIDFieldCount = 1 + 1

// NewApplicationAccessKeyID returns a new instance
// of ApplicationAccessKeyID with the provided attribute values.
func NewApplicationAccessKeyID(
	application ApplicationID,
	idNum ApplicationAccessKeyIDNum,
) ApplicationAccessKeyID {
	return ApplicationAccessKeyID{
		application: application,
		idNum:       idNum,
	}
}

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azid.ID[ApplicationAccessKeyIDNum] = _ApplicationAccessKeyIDZero
var _ azid.BinFieldUnmarshalable = &_ApplicationAccessKeyIDZero
var _ azid.TextUnmarshalable = &_ApplicationAccessKeyIDZero
var _ azcore.AdjunctEntityID[ApplicationAccessKeyIDNum] = _ApplicationAccessKeyIDZero
var _ azcore.ValueObjectAssert[ApplicationAccessKeyID] = _ApplicationAccessKeyIDZero

var _ApplicationAccessKeyIDZero = ApplicationAccessKeyID{
	application: ApplicationIDZero(),
	idNum:       ApplicationAccessKeyIDNumZero,
}

// ApplicationAccessKeyIDZero returns
// a zero-valued instance of ApplicationAccessKeyID.
func ApplicationAccessKeyIDZero() ApplicationAccessKeyID {
	return _ApplicationAccessKeyIDZero
}

// Clone returns a copy of self.
func (id ApplicationAccessKeyID) Clone() ApplicationAccessKeyID { return id }

// AZID is required by azid.ID interface.
func (ApplicationAccessKeyID) AZID() {}

// AZAdjunctEntityID is required
// by azcore.AdjunctEntityID interface.
func (ApplicationAccessKeyID) AZAdjunctEntityID() {}

// IDNum returns the scoped identifier of the entity.
func (id ApplicationAccessKeyID) IDNum() ApplicationAccessKeyIDNum {
	return id.idNum
}

// IDNumPtr returns a pointer to a copy of the id-num if it's considered valid
// otherwise it returns nil.
func (id ApplicationAccessKeyID) IDNumPtr() *ApplicationAccessKeyIDNum {
	if id.IsNotStaticallyValid() {
		return nil
	}
	i := id.IDNum()
	return &i
}

// AZIDNum is required for conformance with azid.ID.
func (id ApplicationAccessKeyID) AZIDNum() ApplicationAccessKeyIDNum {
	return id.idNum
}

// IsZero is required as ApplicationAccessKeyID is a value-object.
func (id ApplicationAccessKeyID) IsZero() bool {
	return id.application.IsZero() &&
		id.idNum == ApplicationAccessKeyIDNumZero
}

// IsStaticallyValid returns true if this instance is valid as an isolated value
// of ApplicationAccessKeyID.
// It doesn't tell whether it refers to a valid instance of ApplicationAccessKey.
func (id ApplicationAccessKeyID) IsStaticallyValid() bool {
	return id.application.IsStaticallyValid() &&
		id.idNum.IsStaticallyValid()
}

// IsNotStaticallyValid returns the negation of value returned by IsStaticallyValid.
func (id ApplicationAccessKeyID) IsNotStaticallyValid() bool {
	return !id.IsStaticallyValid()
}

// Equals is required for conformance with azcore.AdjunctEntityID.
func (id ApplicationAccessKeyID) Equals(other interface{}) bool {
	if x, ok := other.(ApplicationAccessKeyID); ok {
		return id.application.EqualsApplicationID(x.application) &&
			id.idNum == x.idNum
	}
	if x, _ := other.(*ApplicationAccessKeyID); x != nil {
		return id.application.EqualsApplicationID(x.application) &&
			id.idNum == x.idNum
	}
	return false
}

// Equal is required for conformance with azcore.AdjunctEntityID.
func (id ApplicationAccessKeyID) Equal(other interface{}) bool {
	return id.Equals(other)
}

// EqualsApplicationAccessKeyID returns true
// if the other value has the same attributes as id.
func (id ApplicationAccessKeyID) EqualsApplicationAccessKeyID(
	other ApplicationAccessKeyID,
) bool {
	return id.application.EqualsApplicationID(other.application) &&
		id.idNum == other.idNum
}

// AZIDBin is required for conformance
// with azid.ID.
func (id ApplicationAccessKeyID) AZIDBin() []byte {
	data, typ := id.AZIDBinField()
	out := []byte{typ.Byte()}
	return append(out, data...)
}

// ApplicationAccessKeyIDFromAZIDBin creates a new instance of
// ApplicationAccessKeyID from its azid-bin form.
func ApplicationAccessKeyIDFromAZIDBin(
	b []byte,
) (id ApplicationAccessKeyID, readLen int, err error) {
	typ, err := azid.BinDataTypeFromByte(b[0])
	if err != nil {
		return ApplicationAccessKeyIDZero(), 0,
			errors.ArgWrap("", "type parsing", err)
	}
	if typ != azid.BinDataTypeArray {
		return ApplicationAccessKeyIDZero(), 0,
			errors.Arg("", errors.EntMsg("type", "unsupported"))
	}

	id, readLen, err = ApplicationAccessKeyIDFromAZIDBinField(b[1:], typ)
	return id, readLen + 1, err
}

// AZIDBinField is required for conformance
// with azid.ID.
func (id ApplicationAccessKeyID) AZIDBinField() ([]byte, azid.BinDataType) {
	var typesBytes []byte
	var dataBytes []byte
	var fieldBytes []byte
	var fieldType azid.BinDataType

	fieldBytes, fieldType = id.application.AZIDBinField()
	typesBytes = append(typesBytes, fieldType.Byte())
	dataBytes = append(dataBytes, fieldBytes...)

	fieldBytes, fieldType = id.idNum.AZIDBinField()
	typesBytes = append(typesBytes, fieldType.Byte())
	dataBytes = append(dataBytes, fieldBytes...)

	var out = []byte{byte(len(typesBytes))}
	out = append(out, typesBytes...)
	out = append(out, dataBytes...)
	return out, azid.BinDataTypeArray
}

// ApplicationAccessKeyIDFromAZIDBinField creates ApplicationAccessKeyID from
// its azid-bin field form.
func ApplicationAccessKeyIDFromAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (id ApplicationAccessKeyID, readLen int, err error) {
	if typeHint != azid.BinDataTypeArray {
		return ApplicationAccessKeyIDZero(), 0,
			errors.Arg("", errors.EntMsg("type", "unsupported"))
	}

	arrayLen := int(b[0])
	if arrayLen != _ApplicationAccessKeyIDFieldCount {
		return ApplicationAccessKeyIDZero(), 0,
			errors.Arg("", errors.EntMsg("field count", "mismatch"))
	}

	typeCursor := 1
	dataCursor := typeCursor + arrayLen

	var fieldType azid.BinDataType

	fieldType, err = azid.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return ApplicationAccessKeyIDZero(), 0,
			errors.ArgWrap("", "application ref-key type parsing", err)
	}
	typeCursor++
	applicationID, readLen, err := ApplicationIDFromAZIDBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return ApplicationAccessKeyIDZero(), 0,
			errors.ArgWrap("", "application ref-key data parsing", err)
	}
	dataCursor += readLen

	fieldType, err = azid.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return ApplicationAccessKeyIDZero(), 0,
			errors.ArgWrap("", "id-num type parsing", err)
	}
	typeCursor++
	idNum, readLen, err := ApplicationAccessKeyIDNumFromAZIDBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return ApplicationAccessKeyIDZero(), 0,
			errors.ArgWrap("", "id-num data parsing", err)
	}
	dataCursor += readLen

	return ApplicationAccessKeyID{
		application: applicationID,
		idNum:       idNum,
	}, dataCursor, nil
}

// UnmarshalAZIDBinField is required for conformance
// with azcore.BinFieldUnmarshalable.
func (id *ApplicationAccessKeyID) UnmarshalAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (readLen int, err error) {
	i, readLen, err := ApplicationAccessKeyIDFromAZIDBinField(b, typeHint)
	if err == nil {
		*id = i
	}
	return readLen, err
}

const _ApplicationAccessKeyIDAZIDTextPrefix = "KAK0"

// AZIDText is required for conformance
// with azid.ID.
func (id ApplicationAccessKeyID) AZIDText() string {
	if !id.IsStaticallyValid() {
		return ""
	}

	return _ApplicationAccessKeyIDAZIDTextPrefix +
		azid.TextEncode(id.AZIDBin())
}

// ApplicationAccessKeyIDFromAZIDText creates a new instance of
// ApplicationAccessKeyID from its azid-text form.
func ApplicationAccessKeyIDFromAZIDText(s string) (ApplicationAccessKeyID, error) {
	if s == "" {
		return ApplicationAccessKeyIDZero(), nil
	}
	if !strings.HasPrefix(s, _ApplicationAccessKeyIDAZIDTextPrefix) {
		return ApplicationAccessKeyIDZero(),
			errors.Arg("", errors.EntMsg("prefix", "mismatch"))
	}
	s = strings.TrimPrefix(s, _ApplicationAccessKeyIDAZIDTextPrefix)
	b, err := azid.TextDecode(s)
	if err != nil {
		return ApplicationAccessKeyIDZero(),
			errors.ArgWrap("", "data parsing", err)
	}
	id, _, err := ApplicationAccessKeyIDFromAZIDBin(b)
	if err != nil {
		return ApplicationAccessKeyIDZero(),
			errors.ArgWrap("", "data decoding", err)
	}
	return id, nil
}

// UnmarshalAZIDText is required for conformance
// with azid.TextUnmarshalable.
func (id *ApplicationAccessKeyID) UnmarshalAZIDText(s string) error {
	r, err := ApplicationAccessKeyIDFromAZIDText(s)
	if err == nil {
		*id = r
	}
	return err
}

// MarshalText is for compatibility with Go's encoding.TextMarshaler
func (id ApplicationAccessKeyID) MarshalText() ([]byte, error) {
	return []byte(id.AZIDText()), nil
}

// UnmarshalText is for conformance with Go's encoding.TextUnmarshaler
func (id *ApplicationAccessKeyID) UnmarshalText(b []byte) error {
	r, err := ApplicationAccessKeyIDFromAZIDText(string(b))
	if err == nil {
		*id = r
	}
	return err
}

// MarshalJSON makes this type JSON-marshalable.
func (id ApplicationAccessKeyID) MarshalJSON() ([]byte, error) {
	// We assume that there are no symbols in azid-text
	return []byte("\"" + id.AZIDText() + "\""), nil
}

// UnmarshalJSON parses a JSON value.
func (id *ApplicationAccessKeyID) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		*id = ApplicationAccessKeyIDZero()
		return nil
	}
	i, err := ApplicationAccessKeyIDFromAZIDText(s)
	if err == nil {
		*id = i
	}
	return err
}

// Application returns instance's Application value.
func (id ApplicationAccessKeyID) Application() ApplicationID {
	return id.application
}

// ApplicationPtr returns a pointer to a copy of
// ApplicationID if it's considered valid.
func (id ApplicationAccessKeyID) ApplicationPtr() *ApplicationID {
	if id.application.IsStaticallyValid() {
		rk := id.application
		return &rk
	}
	return nil
}

// WithApplication returns a copy
// of ApplicationAccessKeyID
// with its application attribute set to the provided value.
func (id ApplicationAccessKeyID) WithApplication(
	application ApplicationID,
) ApplicationAccessKeyID {
	return ApplicationAccessKeyID{
		application: application,
		idNum:       id.idNum,
	}
}

// ApplicationAccessKeyIDError defines an interface for all
// ApplicationAccessKeyID-related errors.
type ApplicationAccessKeyIDError interface {
	error
	ApplicationAccessKeyIDError()
}

//endregion
