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

// Adjunct-entity Session of Terminal, User.
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

//region ID

// SessionID is a scoped identifier
// used to identify an instance of adjunct entity Session
// scoped within its host entity(s).
type SessionID int32

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azfl.EID = SessionIDZero
var _ azfl.AdjunctEntityID = SessionIDZero
var _ azer.BinFieldUnmarshalable = &_SessionIDZeroVar
var _ azfl.SessionID = SessionIDZero

// SessionIDSignificantBitsMask is used to
// extract significant bits from an instance of SessionID.
const SessionIDSignificantBitsMask uint32 = 0b11111111_11111111_11111111

// SessionIDZero is the zero value for SessionID.
const SessionIDZero = SessionID(0)

// _SessionIDZeroVar is used for testing
// pointer-based interfaces conformance.
var _SessionIDZeroVar = SessionIDZero

// SessionIDFromPrimitiveValue creates an instance
// of SessionID from its primitive value.
func SessionIDFromPrimitiveValue(v int32) SessionID {
	return SessionID(v)
}

// SessionIDFromAZERBinField creates SessionID from
// its azer-bin form.
func SessionIDFromAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (id SessionID, readLen int, err error) {
	if typeHint != azer.BinDataTypeUnspecified && typeHint != azer.BinDataTypeInt32 {
		return SessionID(0), 0,
			errors.ArgMsg("typeHint", "unsupported")
	}
	i := binary.BigEndian.Uint32(b)
	return SessionID(i), 4, nil
}

// PrimitiveValue returns the ID in its primitive type. Prefer to use
// this method instead of casting directly.
func (id SessionID) PrimitiveValue() int32 {
	return int32(id)
}

// AZEID is required
// for conformance with azfl.EID.
func (SessionID) AZEID() {}

// AZAdjunctEntityID is required
// for conformance with azfl.AdjunctEntityID.
func (SessionID) AZAdjunctEntityID() {}

// AZSessionID is required for conformance
// with azfl.SessionID.
func (SessionID) AZSessionID() {}

// IsZero is required as SessionID is a value-object.
func (id SessionID) IsZero() bool {
	return id == SessionIDZero
}

// IsValid returns true if this instance is valid independently as an ID.
// It doesn't tell whether it refers to a valid instance of Session.
func (id SessionID) IsValid() bool {
	return int32(id) > 0 &&
		(uint32(id)&SessionIDSignificantBitsMask) != 0
}

// IsNotValid returns the negation of value returned by IsValid().
func (id SessionID) IsNotValid() bool {
	return !id.IsValid()
}

// AZERBinField is required for conformance
// with azfl.EID.
func (id SessionID) AZERBinField() ([]byte, azer.BinDataType) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(id))
	return b, azer.BinDataTypeInt32
}

// UnmarshalAZERBinField is required for conformance
// with azer.BinFieldUnmarshalable.
func (id *SessionID) UnmarshalAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (readLen int, err error) {
	i, readLen, err := SessionIDFromAZERBinField(b, typeHint)
	if err == nil {
		*id = i
	}
	return readLen, err
}

// Equals is required as SessionID is a value-object.
//
// Use EqualsSessionID method if the other value
// has the same type.
func (id SessionID) Equals(other interface{}) bool {
	if x, ok := other.(SessionID); ok {
		return x == id
	}
	if x, _ := other.(*SessionID); x != nil {
		return *x == id
	}
	return false
}

// Equal is a wrapper for Equals method. It is required for
// compatibility with github.com/google/go-cmp
func (id SessionID) Equal(other interface{}) bool {
	return id.Equals(other)
}

// EqualsSessionID determines if the other instance
// is equal to this instance.
func (id SessionID) EqualsSessionID(
	other SessionID,
) bool {
	return id == other
}

//endregion

//region RefKey

// SessionRefKey is used to identify
// an instance of adjunct entity Session system-wide.
type SessionRefKey struct {
	terminal TerminalRefKey
	id       SessionID
}

// The total number of fields in the struct.
const _SessionRefKeyFieldCount = 1 + 1

// NewSessionRefKey returns a new instance
// of SessionRefKey with the provided attribute values.
func NewSessionRefKey(
	terminal TerminalRefKey,
	id SessionID,
) SessionRefKey {
	return SessionRefKey{
		terminal: terminal,
		id:       id,
	}
}

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azfl.RefKey = _SessionRefKeyZero
var _ azfl.AdjunctEntityRefKey = _SessionRefKeyZero
var _ azfl.SessionRefKey = _SessionRefKeyZero

var _SessionRefKeyZero = SessionRefKey{
	terminal: TerminalRefKeyZero(),
	id:       SessionIDZero,
}

// SessionRefKeyZero returns
// a zero-valued instance of SessionRefKey.
func SessionRefKeyZero() SessionRefKey {
	return _SessionRefKeyZero
}

// AZRefKey is required by azfl.RefKey interface.
func (SessionRefKey) AZRefKey() {}

// AZAdjunctEntityRefKey is required
// by azfl.AdjunctEntityRefKey interface.
func (SessionRefKey) AZAdjunctEntityRefKey() {}

// ID returns the scoped identifier of the entity.
func (refKey SessionRefKey) ID() SessionID {
	return refKey.id
}

// IDPtr returns a pointer to a copy of the ID if it's considered valid
// otherwise it returns nil.
func (refKey SessionRefKey) IDPtr() *SessionID {
	if refKey.IsNotValid() {
		return nil
	}
	i := refKey.ID()
	return &i
}

// EID is required for conformance with azfl.RefKey.
func (refKey SessionRefKey) EID() azfl.EID {
	return refKey.id
}

// SessionID is required for conformance
// with azfl.SessionRefKey.
func (refKey SessionRefKey) SessionID() azfl.SessionID {
	return refKey.id
}

// IsZero is required as SessionRefKey is a value-object.
func (refKey SessionRefKey) IsZero() bool {
	return refKey.terminal.IsZero() &&
		refKey.id == SessionIDZero
}

// IsValid returns true if this instance is valid independently as a ref-key.
// It doesn't tell whether it refers to a valid instance of Session.
func (refKey SessionRefKey) IsValid() bool {
	return refKey.terminal.IsValid() &&
		refKey.id.IsValid()
}

// IsNotValid returns the negation of value returned by IsValid().
func (refKey SessionRefKey) IsNotValid() bool {
	return !refKey.IsValid()
}

// Equals is required for conformance with azfl.AdjunctEntityRefKey.
func (refKey SessionRefKey) Equals(other interface{}) bool {
	if x, ok := other.(SessionRefKey); ok {
		return refKey.terminal.EqualsTerminalRefKey(x.terminal) &&
			refKey.id == x.id
	}
	if x, _ := other.(*SessionRefKey); x != nil {
		return refKey.terminal.EqualsTerminalRefKey(x.terminal) &&
			refKey.id == x.id
	}
	return false
}

// Equal is required for conformance with azfl.AdjunctEntityRefKey.
func (refKey SessionRefKey) Equal(other interface{}) bool {
	return refKey.Equals(other)
}

// EqualsSessionRefKey returns true
// if the other value has the same attributes as refKey.
func (refKey SessionRefKey) EqualsSessionRefKey(
	other SessionRefKey,
) bool {
	return refKey.terminal.EqualsTerminalRefKey(other.terminal) &&
		refKey.id == other.id
}

// AZERBin is required for conformance
// with azfl.RefKey.
func (refKey SessionRefKey) AZERBin() []byte {
	data, typ := refKey.AZERBinField()
	out := []byte{typ.Byte()}
	return append(out, data...)
}

// SessionRefKeyFromAZERBin creates a new instance of
// SessionRefKey from its azer-bin form.
func SessionRefKeyFromAZERBin(
	b []byte,
) (refKey SessionRefKey, readLen int, err error) {
	typ, err := azer.BinDataTypeFromByte(b[0])
	if err != nil {
		return SessionRefKeyZero(), 0,
			errors.ArgWrap("", "type parsing", err)
	}
	if typ != azer.BinDataTypeArray {
		return SessionRefKeyZero(), 0,
			errors.Arg("", errors.EntMsg("type", "unsupported"))
	}

	refKey, readLen, err = SessionRefKeyFromAZERBinField(b[1:], typ)
	return refKey, readLen + 1, err
}

// AZERBinField is required for conformance
// with azfl.RefKey.
func (refKey SessionRefKey) AZERBinField() ([]byte, azer.BinDataType) {
	var typesBytes []byte
	var dataBytes []byte
	var fieldBytes []byte
	var fieldType azer.BinDataType

	fieldBytes, fieldType = refKey.terminal.AZERBinField()
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

// SessionRefKeyFromAZERBinField creates SessionRefKey from
// its azer-bin field form.
func SessionRefKeyFromAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (refKey SessionRefKey, readLen int, err error) {
	if typeHint != azer.BinDataTypeArray {
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

	var fieldType azer.BinDataType

	fieldType, err = azer.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return SessionRefKeyZero(), 0,
			errors.ArgWrap("", "terminal ref-key type parsing", err)
	}
	typeCursor++
	terminalRefKey, readLen, err := TerminalRefKeyFromAZERBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return SessionRefKeyZero(), 0,
			errors.ArgWrap("", "terminal ref-key data parsing", err)
	}
	dataCursor += readLen

	fieldType, err = azer.BinDataTypeFromByte(b[typeCursor])
	if err != nil {
		return SessionRefKeyZero(), 0,
			errors.ArgWrap("", "id type parsing", err)
	}
	typeCursor++
	id, readLen, err := SessionIDFromAZERBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return SessionRefKeyZero(), 0,
			errors.ArgWrap("", "id data parsing", err)
	}
	dataCursor += readLen

	return SessionRefKey{
		terminal: terminalRefKey,
		id:       id,
	}, dataCursor, nil
}

// UnmarshalAZERBinField is required for conformance
// with azfl.BinFieldUnmarshalable.
func (refKey *SessionRefKey) UnmarshalAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (readLen int, err error) {
	i, readLen, err := SessionRefKeyFromAZERBinField(b, typeHint)
	if err == nil {
		*refKey = i
	}
	return readLen, err
}

const _SessionRefKeyAZERTextPrefix = "KSe0"

// AZERText is required for conformance
// with azfl.RefKey.
func (refKey SessionRefKey) AZERText() string {
	if !refKey.IsValid() {
		return ""
	}

	return _SessionRefKeyAZERTextPrefix +
		azer.TextEncode(refKey.AZERBin())
}

// SessionRefKeyFromAZERText creates a new instance of
// SessionRefKey from its azer-text form.
func SessionRefKeyFromAZERText(s string) (SessionRefKey, error) {
	if s == "" {
		return SessionRefKeyZero(), nil
	}
	if !strings.HasPrefix(s, _SessionRefKeyAZERTextPrefix) {
		return SessionRefKeyZero(),
			errors.Arg("", errors.EntMsg("prefix", "mismatch"))
	}
	s = strings.TrimPrefix(s, _SessionRefKeyAZERTextPrefix)
	b, err := azer.TextDecode(s)
	if err != nil {
		return SessionRefKeyZero(),
			errors.ArgWrap("", "data parsing", err)
	}
	refKey, _, err := SessionRefKeyFromAZERBin(b)
	if err != nil {
		return SessionRefKeyZero(),
			errors.ArgWrap("", "data decoding", err)
	}
	return refKey, nil
}

// UnmarshalAZERText is required for conformance
// with azer.TextUnmarshalable.
func (refKey *SessionRefKey) UnmarshalAZERText(s string) error {
	r, err := SessionRefKeyFromAZERText(s)
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalText is for compatibility with Go's encoding.TextMarshaler
func (refKey SessionRefKey) MarshalText() ([]byte, error) {
	return []byte(refKey.AZERText()), nil
}

// UnmarshalText is for conformance with Go's encoding.TextUnmarshaler
func (refKey *SessionRefKey) UnmarshalText(b []byte) error {
	r, err := SessionRefKeyFromAZERText(string(b))
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalJSON makes this type JSON-marshalable.
func (refKey SessionRefKey) MarshalJSON() ([]byte, error) {
	// We assume that there's no symbols in azer-text
	return []byte("\"" + refKey.AZERText() + "\""), nil
}

// UnmarshalJSON parses a JSON value.
func (refKey *SessionRefKey) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		*refKey = SessionRefKeyZero()
		return nil
	}
	i, err := SessionRefKeyFromAZERText(s)
	if err == nil {
		*refKey = i
	}
	return err
}

// Terminal returns instance's Terminal value.
func (refKey SessionRefKey) Terminal() TerminalRefKey {
	return refKey.terminal
}

// WithTerminal returns a copy
// of SessionRefKey
// with its terminal attribute set to the provided value.
func (refKey SessionRefKey) WithTerminal(
	terminal TerminalRefKey,
) SessionRefKey {
	return SessionRefKey{
		terminal: terminal,
		id:       refKey.id,
	}
}

// SessionRefKeyError defines an interface for all
// SessionRefKey-related errors.
type SessionRefKeyError interface {
	error
	SessionRefKeyError()
}

//endregion
