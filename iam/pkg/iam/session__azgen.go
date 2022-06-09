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

// Adjunct-entity Session of Terminal.
//
// A Session represents authorization for a time span. While Terminal
// usually provide longer authorization period, a Session is used to
// break down that authorization into smaller time spans.
//
// If a Session is expired or revoked, the previously authorized
// Application instance (Terminal) could request another Session as long the
// Application's authorization is still valid. There's only one instance
// of Session active at a time for a Terminal.
//
// An access token represents a instance of Session.

//region IDNum

// SessionIDNum is a scoped identifier
// used to identify an instance of adjunct entity Session
// scoped within its host entity(s).
type SessionIDNum int32

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azid.IDNumMethods = SessionIDNumZero
var _ azid.BinFieldUnmarshalable = &_SessionIDNumZeroVar
var _ azcore.AdjunctEntityIDNumMethods = SessionIDNumZero
var _ azcore.SessionIDNumMethods = SessionIDNumZero

// SessionIDNumIdentifierBitsMask is used to
// extract identifier bits from an instance of SessionIDNum.
const SessionIDNumIdentifierBitsMask uint32 = 0b_00000000_11111111_11111111_11111111

// SessionIDNumZero is the zero value for SessionIDNum.
const SessionIDNumZero = SessionIDNum(0)

// _SessionIDNumZeroVar is used for testing
// pointer-based interfaces conformance.
var _SessionIDNumZeroVar = SessionIDNumZero

// SessionIDNumFromPrimitiveValue creates an instance
// of SessionIDNum from its primitive value.
func SessionIDNumFromPrimitiveValue(v int32) SessionIDNum {
	return SessionIDNum(v)
}

// SessionIDNumFromAZIDBinField creates SessionIDNum from
// its azid-bin form.
func SessionIDNumFromAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (idNum SessionIDNum, readLen int, err error) {
	if typeHint != azid.BinDataTypeUnspecified && typeHint != azid.BinDataTypeInt32 {
		return SessionIDNum(0), 0,
			errors.ArgMsg("typeHint", "unsupported")
	}
	i := binary.BigEndian.Uint32(b)
	return SessionIDNum(i), 4, nil
}

// PrimitiveValue returns the value in its primitive type. Prefer to use
// this method instead of casting directly.
func (idNum SessionIDNum) PrimitiveValue() int32 {
	return int32(idNum)
}

// AZIDNum is required
// for conformance with azid.IDNum.
func (SessionIDNum) AZIDNum() {}

// AZAdjunctEntityIDNum is required
// for conformance with azcore.AdjunctEntityIDNum.
func (SessionIDNum) AZAdjunctEntityIDNum() {}

// AZSessionIDNum is required for conformance
// with azcore.SessionIDNum.
func (SessionIDNum) AZSessionIDNum() {}

// IsZero is required as SessionIDNum is a value-object.
func (idNum SessionIDNum) IsZero() bool {
	return idNum == SessionIDNumZero
}

// IsStaticallyValid returns true if this instance is valid as an isolated value
// of SessionIDNum. It doesn't tell whether it refers to
// a valid instance of Session.
func (idNum SessionIDNum) IsStaticallyValid() bool {
	return int32(idNum) > 0 &&
		(uint32(idNum)&SessionIDNumIdentifierBitsMask) != 0
}

// IsNotStaticallyValid returns the negation of value returned by IsStaticallyValid.
func (idNum SessionIDNum) IsNotStaticallyValid() bool {
	return !idNum.IsStaticallyValid()
}

// Equals is required as SessionIDNum is a value-object.
//
// Use EqualsSessionIDNum method if the other value
// has the same type.
func (idNum SessionIDNum) Equals(other interface{}) bool {
	if x, ok := other.(SessionIDNum); ok {
		return x == idNum
	}
	if x, _ := other.(*SessionIDNum); x != nil {
		return *x == idNum
	}
	return false
}

// Equal is a wrapper for Equals method. It is required for
// compatibility with github.com/google/go-cmp
func (idNum SessionIDNum) Equal(other interface{}) bool {
	return idNum.Equals(other)
}

// EqualsSessionIDNum determines if the other instance
// is equal to this instance.
func (idNum SessionIDNum) EqualsSessionIDNum(
	other SessionIDNum,
) bool {
	return idNum == other
}

// AZIDBinField is required for conformance
// with azid.IDNum.
func (idNum SessionIDNum) AZIDBinField() ([]byte, azid.BinDataType) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(idNum))
	return b, azid.BinDataTypeInt32
}

// UnmarshalAZIDBinField is required for conformance
// with azid.BinFieldUnmarshalable.
func (idNum *SessionIDNum) UnmarshalAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (readLen int, err error) {
	i, readLen, err := SessionIDNumFromAZIDBinField(b, typeHint)
	if err == nil {
		*idNum = i
	}
	return readLen, err
}

// Embedded fields
const (
	SessionIDNumEmbeddedFieldsMask = 0b_00000000_00000000_00000000_00000000
)

//endregion

//region RefKey

// SessionRefKey is used to identify
// an instance of adjunct entity Session system-wide.
type SessionRefKey struct {
	terminal TerminalRefKey
	idNum    SessionIDNum
}

// The total number of fields in the struct.
const _SessionRefKeyFieldCount = 1 + 1

// NewSessionRefKey returns a new instance
// of SessionRefKey with the provided attribute values.
func NewSessionRefKey(
	terminal TerminalRefKey,
	idNum SessionIDNum,
) SessionRefKey {
	return SessionRefKey{
		terminal: terminal,
		idNum:    idNum,
	}
}

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azid.RefKey[SessionIDNum] = _SessionRefKeyZero
var _ azid.BinFieldUnmarshalable = &_SessionRefKeyZero
var _ azid.TextUnmarshalable = &_SessionRefKeyZero
var _ azcore.AdjunctEntityRefKey[SessionIDNum] = _SessionRefKeyZero
var _ azcore.SessionRefKey[SessionIDNum] = _SessionRefKeyZero

var _SessionRefKeyZero = SessionRefKey{
	terminal: TerminalRefKeyZero(),
	idNum:    SessionIDNumZero,
}

// SessionRefKeyZero returns
// a zero-valued instance of SessionRefKey.
func SessionRefKeyZero() SessionRefKey {
	return _SessionRefKeyZero
}

// AZRefKey is required by azid.RefKey interface.
func (SessionRefKey) AZRefKey() {}

// AZAdjunctEntityRefKey is required
// by azcore.AdjunctEntityRefKey interface.
func (SessionRefKey) AZAdjunctEntityRefKey() {}

// IDNum returns the scoped identifier of the entity.
func (refKey SessionRefKey) IDNum() SessionIDNum {
	return refKey.idNum
}

// IDNumPtr returns a pointer to a copy of the id-num if it's considered valid
// otherwise it returns nil.
func (refKey SessionRefKey) IDNumPtr() *SessionIDNum {
	if refKey.IsNotStaticallyValid() {
		return nil
	}
	i := refKey.IDNum()
	return &i
}

// AZIDNum is required for conformance with azid.RefKey.
func (refKey SessionRefKey) AZIDNum() SessionIDNum {
	return refKey.idNum
}

// SessionIDNum is required for conformance
// with azcore.SessionRefKey.
func (refKey SessionRefKey) SessionIDNum() SessionIDNum {
	return refKey.idNum
}

// IsZero is required as SessionRefKey is a value-object.
func (refKey SessionRefKey) IsZero() bool {
	return refKey.terminal.IsZero() &&
		refKey.idNum == SessionIDNumZero
}

// IsStaticallyValid returns true if this instance is valid as an isolated value
// of SessionRefKey.
// It doesn't tell whether it refers to a valid instance of Session.
func (refKey SessionRefKey) IsStaticallyValid() bool {
	return refKey.terminal.IsStaticallyValid() &&
		refKey.idNum.IsStaticallyValid()
}

// IsNotStaticallyValid returns the negation of value returned by IsStaticallyValid.
func (refKey SessionRefKey) IsNotStaticallyValid() bool {
	return !refKey.IsStaticallyValid()
}

// Equals is required for conformance with azcore.AdjunctEntityRefKey.
func (refKey SessionRefKey) Equals(other interface{}) bool {
	if x, ok := other.(SessionRefKey); ok {
		return refKey.terminal.EqualsTerminalRefKey(x.terminal) &&
			refKey.idNum == x.idNum
	}
	if x, _ := other.(*SessionRefKey); x != nil {
		return refKey.terminal.EqualsTerminalRefKey(x.terminal) &&
			refKey.idNum == x.idNum
	}
	return false
}

// Equal is required for conformance with azcore.AdjunctEntityRefKey.
func (refKey SessionRefKey) Equal(other interface{}) bool {
	return refKey.Equals(other)
}

// EqualsSessionRefKey returns true
// if the other value has the same attributes as refKey.
func (refKey SessionRefKey) EqualsSessionRefKey(
	other SessionRefKey,
) bool {
	return refKey.terminal.EqualsTerminalRefKey(other.terminal) &&
		refKey.idNum == other.idNum
}

// AZIDBin is required for conformance
// with azid.RefKey.
func (refKey SessionRefKey) AZIDBin() []byte {
	data, typ := refKey.AZIDBinField()
	out := []byte{typ.Byte()}
	return append(out, data...)
}

// SessionRefKeyFromAZIDBin creates a new instance of
// SessionRefKey from its azid-bin form.
func SessionRefKeyFromAZIDBin(
	b []byte,
) (refKey SessionRefKey, readLen int, err error) {
	typ, err := azid.BinDataTypeFromByte(b[0])
	if err != nil {
		return SessionRefKeyZero(), 0,
			errors.ArgWrap("", "type parsing", err)
	}
	if typ != azid.BinDataTypeArray {
		return SessionRefKeyZero(), 0,
			errors.Arg("", errors.EntMsg("type", "unsupported"))
	}

	refKey, readLen, err = SessionRefKeyFromAZIDBinField(b[1:], typ)
	return refKey, readLen + 1, err
}

// AZIDBinField is required for conformance
// with azid.RefKey.
func (refKey SessionRefKey) AZIDBinField() ([]byte, azid.BinDataType) {
	var typesBytes []byte
	var dataBytes []byte
	var fieldBytes []byte
	var fieldType azid.BinDataType

	fieldBytes, fieldType = refKey.terminal.AZIDBinField()
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

// SessionRefKeyFromAZIDBinField creates SessionRefKey from
// its azid-bin field form.
func SessionRefKeyFromAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (refKey SessionRefKey, readLen int, err error) {
	if typeHint != azid.BinDataTypeArray {
		return SessionRefKeyZero(), 0,
			errors.Arg("", errors.EntMsg("type", "unsupported"))
	}

	arrayLen := int(b[0])
	if arrayLen != _SessionRefKeyFieldCount {
		return SessionRefKeyZero(), 0,
			errors.Arg("", errors.EntMsg("field count", "mismatch"))
	}

	typeCursor := 1
	dataCursor := typeCursor + arrayLen

	var fieldType azid.BinDataType

	fieldType, err = azid.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return SessionRefKeyZero(), 0,
			errors.ArgWrap("", "terminal ref-key type parsing", err)
	}
	typeCursor++
	terminalRefKey, readLen, err := TerminalRefKeyFromAZIDBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return SessionRefKeyZero(), 0,
			errors.ArgWrap("", "terminal ref-key data parsing", err)
	}
	dataCursor += readLen

	fieldType, err = azid.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return SessionRefKeyZero(), 0,
			errors.ArgWrap("", "id-num type parsing", err)
	}
	typeCursor++
	idNum, readLen, err := SessionIDNumFromAZIDBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return SessionRefKeyZero(), 0,
			errors.ArgWrap("", "id-num data parsing", err)
	}
	dataCursor += readLen

	return SessionRefKey{
		terminal: terminalRefKey,
		idNum:    idNum,
	}, dataCursor, nil
}

// UnmarshalAZIDBinField is required for conformance
// with azcore.BinFieldUnmarshalable.
func (refKey *SessionRefKey) UnmarshalAZIDBinField(
	b []byte, typeHint azid.BinDataType,
) (readLen int, err error) {
	i, readLen, err := SessionRefKeyFromAZIDBinField(b, typeHint)
	if err == nil {
		*refKey = i
	}
	return readLen, err
}

const _SessionRefKeyAZIDTextPrefix = "KSe0"

// AZIDText is required for conformance
// with azid.RefKey.
func (refKey SessionRefKey) AZIDText() string {
	if !refKey.IsStaticallyValid() {
		return ""
	}

	return _SessionRefKeyAZIDTextPrefix +
		azid.TextEncode(refKey.AZIDBin())
}

// SessionRefKeyFromAZIDText creates a new instance of
// SessionRefKey from its azid-text form.
func SessionRefKeyFromAZIDText(s string) (SessionRefKey, error) {
	if s == "" {
		return SessionRefKeyZero(), nil
	}
	if !strings.HasPrefix(s, _SessionRefKeyAZIDTextPrefix) {
		return SessionRefKeyZero(),
			errors.Arg("", errors.EntMsg("prefix", "mismatch"))
	}
	s = strings.TrimPrefix(s, _SessionRefKeyAZIDTextPrefix)
	b, err := azid.TextDecode(s)
	if err != nil {
		return SessionRefKeyZero(),
			errors.ArgWrap("", "data parsing", err)
	}
	refKey, _, err := SessionRefKeyFromAZIDBin(b)
	if err != nil {
		return SessionRefKeyZero(),
			errors.ArgWrap("", "data decoding", err)
	}
	return refKey, nil
}

// UnmarshalAZIDText is required for conformance
// with azid.TextUnmarshalable.
func (refKey *SessionRefKey) UnmarshalAZIDText(s string) error {
	r, err := SessionRefKeyFromAZIDText(s)
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalText is for compatibility with Go's encoding.TextMarshaler
func (refKey SessionRefKey) MarshalText() ([]byte, error) {
	return []byte(refKey.AZIDText()), nil
}

// UnmarshalText is for conformance with Go's encoding.TextUnmarshaler
func (refKey *SessionRefKey) UnmarshalText(b []byte) error {
	r, err := SessionRefKeyFromAZIDText(string(b))
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalJSON makes this type JSON-marshalable.
func (refKey SessionRefKey) MarshalJSON() ([]byte, error) {
	// We assume that there are no symbols in azid-text
	return []byte("\"" + refKey.AZIDText() + "\""), nil
}

// UnmarshalJSON parses a JSON value.
func (refKey *SessionRefKey) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		*refKey = SessionRefKeyZero()
		return nil
	}
	i, err := SessionRefKeyFromAZIDText(s)
	if err == nil {
		*refKey = i
	}
	return err
}

// Terminal returns instance's Terminal value.
func (refKey SessionRefKey) Terminal() TerminalRefKey {
	return refKey.terminal
}

// TerminalPtr returns a pointer to a copy of
// TerminalRefKey if it's considered valid.
func (refKey SessionRefKey) TerminalPtr() *TerminalRefKey {
	if refKey.terminal.IsStaticallyValid() {
		rk := refKey.terminal
		return &rk
	}
	return nil
}

// WithTerminal returns a copy
// of SessionRefKey
// with its terminal attribute set to the provided value.
func (refKey SessionRefKey) WithTerminal(
	terminal TerminalRefKey,
) SessionRefKey {
	return SessionRefKey{
		terminal: terminal,
		idNum:    refKey.idNum,
	}
}

// SessionRefKeyError defines an interface for all
// SessionRefKey-related errors.
type SessionRefKeyError interface {
	error
	SessionRefKeyError()
}

//endregion
