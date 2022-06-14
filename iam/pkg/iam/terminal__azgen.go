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
var _ = errors.ErrUnimplemented
var _ = binary.MaxVarintLen16
var _ = rand.Reader
var _ = strings.Compare

// Adjunct-entity Terminal of Application, User.
//
// A Terminal is an authorized instance of application. A new instance of
// Terminal will be created upon sucessful authorization-authentication.
//
// As long the the
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
// A sucessful authorization will generate both a new instance of
// Terminal and an instance of initial Session.

//region ID

// TerminalID is used to identify
// an instance of adjunct entity Terminal system-wide.
type TerminalID struct {
	application ApplicationID
	user        UserID
	idNum       TerminalIDNum
}

// The total number of fields in the struct.
const _TerminalIDFieldCount = 2 + 1

// NewTerminalID returns a new instance
// of TerminalID with the provided attribute values.
func NewTerminalID(
	application ApplicationID,
	user UserID,
	idNum TerminalIDNum,
) TerminalID {
	return TerminalID{
		application: application,
		user:        user,
		idNum:       idNum,
	}
}

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azid.ID[TerminalIDNum] = _TerminalIDZero
var _ azid.BinFieldUnmarshalable = &_TerminalIDZero
var _ azid.TextUnmarshalable = &_TerminalIDZero
var _ azcore.AdjunctEntityID[TerminalIDNum] = _TerminalIDZero
var _ azcore.ValueObjectAssert[TerminalID] = _TerminalIDZero
var _ azcore.TerminalID[TerminalIDNum] = _TerminalIDZero

var _TerminalIDZero = TerminalID{
	application: ApplicationIDZero(),
	user:        UserIDZero(),
	idNum:       TerminalIDNumZero,
}

// TerminalIDZero returns
// a zero-valued instance of TerminalID.
func TerminalIDZero() TerminalID {
	return _TerminalIDZero
}

// Clone returns a copy of self.
func (id TerminalID) Clone() TerminalID { return id }

// AZID is required by azid.ID interface.
func (TerminalID) AZID() {}

// AZAdjunctEntityID is required
// by azcore.AdjunctEntityID interface.
func (TerminalID) AZAdjunctEntityID() {}

// IDNum returns the scoped identifier of the entity.
func (id TerminalID) IDNum() TerminalIDNum {
	return id.idNum
}

// IDNumPtr returns a pointer to a copy of the id-num if it's considered valid
// otherwise it returns nil.
func (id TerminalID) IDNumPtr() *TerminalIDNum {
	if id.IsNotStaticallyValid() {
		return nil
	}
	i := id.IDNum()
	return &i
}

// AZIDNum is required for conformance with azid.ID.
func (id TerminalID) AZIDNum() TerminalIDNum {
	return id.idNum
}

// TerminalIDNum is required for conformance
// with azcore.TerminalID.
func (id TerminalID) TerminalIDNum() TerminalIDNum {
	return id.idNum
}

// IsZero is required as TerminalID is a value-object.
func (id TerminalID) IsZero() bool {
	return id.application.IsZero() &&
		id.user.IsZero() &&
		id.idNum == TerminalIDNumZero
}

// IsStaticallyValid returns true if this instance is valid as an isolated value
// of TerminalID.
// It doesn't tell whether it refers to a valid instance of Terminal.
func (id TerminalID) IsStaticallyValid() bool {
	return id.application.IsStaticallyValid() &&
		id.user.IsStaticallyValid() &&
		id.idNum.IsStaticallyValid()
}

// IsNotStaticallyValid returns the negation of value returned by IsStaticallyValid.
func (id TerminalID) IsNotStaticallyValid() bool {
	return !id.IsStaticallyValid()
}

// Equals is required for conformance with azcore.AdjunctEntityID.
func (id TerminalID) Equals(other interface{}) bool {
	if x, ok := other.(TerminalID); ok {
		return id.application.EqualsApplicationID(x.application) &&
			id.user.EqualsUserID(x.user) &&
			id.idNum == x.idNum
	}
	if x, _ := other.(*TerminalID); x != nil {
		return id.application.EqualsApplicationID(x.application) &&
			id.user.EqualsUserID(x.user) &&
			id.idNum == x.idNum
	}
	return false
}

// Equal is required for conformance with azcore.AdjunctEntityID.
func (id TerminalID) Equal(other interface{}) bool {
	return id.Equals(other)
}

// EqualsTerminalID returns true
// if the other value has the same attributes as id.
func (id TerminalID) EqualsTerminalID(
	other TerminalID,
) bool {
	return id.application.EqualsApplicationID(other.application) &&
		id.user.EqualsUserID(other.user) &&
		id.idNum == other.idNum
}

// AZIDBin is required for conformance
// with azid.ID.
func (id TerminalID) AZIDBin() []byte {
	data, typ := id.AZIDBinField()
	out := []byte{typ.Byte()}
	return append(out, data...)
}

// TerminalIDFromAZIDBin creates a new instance of
// TerminalID from its azid-bin form.
func TerminalIDFromAZIDBin(
	b []byte,
) (id TerminalID, readLen int, err error) {
	typ, err := azid.BinDataTypeFromByte(b[0])
	if err != nil {
		return TerminalIDZero(), 0,
			errors.ArgWrap("", "type parsing", err)
	}
	if typ != azid.BinDataTypeArray {
		return TerminalIDZero(), 0,
			errors.Arg("", errors.EntMsg("type", "unsupported"))
	}

	id, readLen, err = TerminalIDFromAZIDBinField(b[1:], typ)
	return id, readLen + 1, err
}

// AZIDBinField is required for conformance
// with azid.ID.
func (id TerminalID) AZIDBinField() ([]byte, azid.BinDataType) {
	var typesBytes []byte
	var dataBytes []byte
	var fieldBytes []byte
	var fieldType azid.BinDataType

	fieldBytes, fieldType = id.application.AZIDBinField()
	typesBytes = append(typesBytes, fieldType.Byte())
	dataBytes = append(dataBytes, fieldBytes...)

	fieldBytes, fieldType = id.user.AZIDBinField()
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

// TerminalIDFromAZIDBinField creates TerminalID from
// its azid-bin field form.
func TerminalIDFromAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (id TerminalID, readLen int, err error) {
	if typeHint != azid.BinDataTypeArray {
		return TerminalIDZero(), 0,
			errors.Arg("", errors.EntMsg("type", "unsupported"))
	}

	arrayLen := int(b[0])
	if arrayLen != _TerminalIDFieldCount {
		return TerminalIDZero(), 0,
			errors.Arg("", errors.EntMsg("field count", "mismatch"))
	}

	typeCursor := 1
	dataCursor := typeCursor + arrayLen

	var fieldType azid.BinDataType

	fieldType, err = azid.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return TerminalIDZero(), 0,
			errors.ArgWrap("", "application ref-key type parsing", err)
	}
	typeCursor++
	applicationID, readLen, err := ApplicationIDFromAZIDBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return TerminalIDZero(), 0,
			errors.ArgWrap("", "application ref-key data parsing", err)
	}
	dataCursor += readLen

	fieldType, err = azid.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return TerminalIDZero(), 0,
			errors.ArgWrap("", "user ref-key type parsing", err)
	}
	typeCursor++
	userID, readLen, err := UserIDFromAZIDBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return TerminalIDZero(), 0,
			errors.ArgWrap("", "user ref-key data parsing", err)
	}
	dataCursor += readLen

	fieldType, err = azid.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return TerminalIDZero(), 0,
			errors.ArgWrap("", "id-num type parsing", err)
	}
	typeCursor++
	idNum, readLen, err := TerminalIDNumFromAZIDBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return TerminalIDZero(), 0,
			errors.ArgWrap("", "id-num data parsing", err)
	}
	dataCursor += readLen

	return TerminalID{
		application: applicationID,
		user:        userID,
		idNum:       idNum,
	}, dataCursor, nil
}

// UnmarshalAZIDBinField is required for conformance
// with azcore.BinFieldUnmarshalable.
func (id *TerminalID) UnmarshalAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (readLen int, err error) {
	i, readLen, err := TerminalIDFromAZIDBinField(b, typeHint)
	if err == nil {
		*id = i
	}
	return readLen, err
}

const _TerminalIDAZIDTextPrefix = "KTx0"

// AZIDText is required for conformance
// with azid.ID.
func (id TerminalID) AZIDText() string {
	if !id.IsStaticallyValid() {
		return ""
	}

	return _TerminalIDAZIDTextPrefix +
		azid.TextEncode(id.AZIDBin())
}

// TerminalIDFromAZIDText creates a new instance of
// TerminalID from its azid-text form.
func TerminalIDFromAZIDText(s string) (TerminalID, error) {
	if s == "" {
		return TerminalIDZero(), nil
	}
	if !strings.HasPrefix(s, _TerminalIDAZIDTextPrefix) {
		return TerminalIDZero(),
			errors.Arg("", errors.EntMsg("prefix", "mismatch"))
	}
	s = strings.TrimPrefix(s, _TerminalIDAZIDTextPrefix)
	b, err := azid.TextDecode(s)
	if err != nil {
		return TerminalIDZero(),
			errors.ArgWrap("", "data parsing", err)
	}
	id, _, err := TerminalIDFromAZIDBin(b)
	if err != nil {
		return TerminalIDZero(),
			errors.ArgWrap("", "data decoding", err)
	}
	return id, nil
}

// UnmarshalAZIDText is required for conformance
// with azid.TextUnmarshalable.
func (id *TerminalID) UnmarshalAZIDText(s string) error {
	r, err := TerminalIDFromAZIDText(s)
	if err == nil {
		*id = r
	}
	return err
}

// MarshalText is for compatibility with Go's encoding.TextMarshaler
func (id TerminalID) MarshalText() ([]byte, error) {
	return []byte(id.AZIDText()), nil
}

// UnmarshalText is for conformance with Go's encoding.TextUnmarshaler
func (id *TerminalID) UnmarshalText(b []byte) error {
	r, err := TerminalIDFromAZIDText(string(b))
	if err == nil {
		*id = r
	}
	return err
}

// MarshalJSON makes this type JSON-marshalable.
func (id TerminalID) MarshalJSON() ([]byte, error) {
	// We assume that there are no symbols in azid-text
	return []byte("\"" + id.AZIDText() + "\""), nil
}

// UnmarshalJSON parses a JSON value.
func (id *TerminalID) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		*id = TerminalIDZero()
		return nil
	}
	i, err := TerminalIDFromAZIDText(s)
	if err == nil {
		*id = i
	}
	return err
}

// Application returns instance's Application value.
func (id TerminalID) Application() ApplicationID {
	return id.application
}

// ApplicationPtr returns a pointer to a copy of
// ApplicationID if it's considered valid.
func (id TerminalID) ApplicationPtr() *ApplicationID {
	if id.application.IsStaticallyValid() {
		rk := id.application
		return &rk
	}
	return nil
}

// WithApplication returns a copy
// of TerminalID
// with its application attribute set to the provided value.
func (id TerminalID) WithApplication(
	application ApplicationID,
) TerminalID {
	return TerminalID{
		application: application,
		user:        id.user,
		idNum:       id.idNum,
	}
}

// User returns instance's User value.
func (id TerminalID) User() UserID {
	return id.user
}

// UserPtr returns a pointer to a copy of
// UserID if it's considered valid.
func (id TerminalID) UserPtr() *UserID {
	if id.user.IsStaticallyValid() {
		rk := id.user
		return &rk
	}
	return nil
}

// WithUser returns a copy
// of TerminalID
// with its user attribute set to the provided value.
func (id TerminalID) WithUser(
	user UserID,
) TerminalID {
	return TerminalID{
		application: id.application,
		user:        user,
		idNum:       id.idNum,
	}
}

// TerminalIDError defines an interface for all
// TerminalID-related errors.
type TerminalIDError interface {
	error
	TerminalIDError()
}

//endregion

//region IDNum

// TerminalIDNum is a scoped identifier
// used to identify an instance of adjunct entity Terminal
// scoped within its host entity(s).
type TerminalIDNum int64

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azid.IDNumMethods = TerminalIDNumZero
var _ azid.BinFieldUnmarshalable = &_TerminalIDNumZeroVar
var _ azcore.AdjunctEntityIDNumMethods = TerminalIDNumZero
var _ azcore.ValueObjectAssert[TerminalIDNum] = TerminalIDNumZero
var _ azcore.TerminalIDNumMethods = TerminalIDNumZero

// TerminalIDNumIdentifierBitsMask is used to
// extract identifier bits from an instance of TerminalIDNum.
const TerminalIDNumIdentifierBitsMask uint64 = 0b_00000000_11111111_11111111_11111111_11111111_11111111_11111111_11111111

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

// Clone returns a copy of self.
func (idNum TerminalIDNum) Clone() TerminalIDNum { return idNum }

// AZIDNum is required
// for conformance with azid.IDNum.
func (TerminalIDNum) AZIDNum() {}

// AZAdjunctEntityIDNum is required
// for conformance with azcore.AdjunctEntityIDNum.
func (TerminalIDNum) AZAdjunctEntityIDNum() {}

// AZTerminalIDNum is required for conformance
// with azcore.TerminalIDNum.
func (TerminalIDNum) AZTerminalIDNum() {}

// IsZero is required as TerminalIDNum is a value-object.
func (idNum TerminalIDNum) IsZero() bool {
	return idNum == TerminalIDNumZero
}

// IsStaticallyValid returns true if this instance is valid as an isolated value
// of TerminalIDNum. It doesn't tell whether it refers to
// a valid instance of Terminal.
func (idNum TerminalIDNum) IsStaticallyValid() bool {
	return int64(idNum) > 0 &&
		(uint64(idNum)&TerminalIDNumIdentifierBitsMask) != 0
}

// IsNotStaticallyValid returns the negation of value returned by IsStaticallyValid.
func (idNum TerminalIDNum) IsNotStaticallyValid() bool {
	return !idNum.IsStaticallyValid()
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

// Embedded fields
const (
	TerminalIDNumEmbeddedFieldsMask = 0b_00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
)

//endregion
