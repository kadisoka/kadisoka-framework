package iam

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"sync"

	azcore "github.com/alloyzeus/go-azcore/azcore"
	azer "github.com/alloyzeus/go-azcore/azcore/azer"
	"github.com/alloyzeus/go-azcore/azcore/errors"
	"github.com/richardlehane/crock32"
)

//region ID

// TerminalID is a scoped identifier
// used to identify an instance of adjunct entity Terminal
// scoped within its host entity(s).
type TerminalID int64

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azcore.EID = TerminalIDZero
var _ azcore.AdjunctEntityID = TerminalIDZero
var _ azer.BinFieldUnmarshalable = &_TerminalIDZeroVar
var _ azcore.TerminalID = TerminalIDZero

// _TerminalIDSignificantBitsMask is used to
// extract significant bits from an instance of TerminalID.
const _TerminalIDSignificantBitsMask uint64 = 0b11111111_11111111_11111111_11111111_11111111_11111111_11111111

// TerminalIDZero is the zero value for TerminalID.
const TerminalIDZero = TerminalID(0)

// _TerminalIDZeroVar is used for testing
// pointer-based interfaces conformance.
var _TerminalIDZeroVar = TerminalIDZero

// TerminalIDFromPrimitiveValue creates an instance
// of TerminalID from its primitive value.
func TerminalIDFromPrimitiveValue(v int64) TerminalID {
	return TerminalID(v)
}

// TerminalIDFromAZERBinField creates TerminalID from
// its azer-bin form.
func TerminalIDFromAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (id TerminalID, readLen int, err error) {
	if typeHint != azer.BinDataTypeUnspecified && typeHint != azer.BinDataTypeInt64 {
		return TerminalID(0), 0,
			errors.ArgMsg("typeHint", "unsupported")
	}
	i := binary.BigEndian.Uint64(b)
	return TerminalID(i), 8, nil
}

// PrimitiveValue returns the ID in its primitive type. Prefer to use
// this method instead of casting directly.
func (id TerminalID) PrimitiveValue() int64 {
	return int64(id)
}

// AZEID is required
// for conformance with azcore.EID.
func (TerminalID) AZEID() {}

// AZAdjunctEntityID is required
// for conformance with azcore.AdjunctEntityID.
func (TerminalID) AZAdjunctEntityID() {}

// AZTerminalID is required for conformance
// with azcore.TerminalID.
func (TerminalID) AZTerminalID() {}

// IsZero is required as TerminalID is a value-object.
func (id TerminalID) IsZero() bool {
	return id == TerminalIDZero
}

// IsValid returns true if this instance is valid independently as an ID.
// It doesn't tell whether it refers to a valid instance of Terminal.
func (id TerminalID) IsValid() bool {
	return int64(id) > 0 &&
		(uint64(id)&_TerminalIDSignificantBitsMask) != 0
}

// IsNotValid returns the negation of value returned by IsValid().
func (id TerminalID) IsNotValid() bool {
	return !id.IsValid()
}

// AZERBinField is required for conformance
// with azcore.EID.
func (id TerminalID) AZERBinField() ([]byte, azer.BinDataType) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(id))
	return b, azer.BinDataTypeInt64
}

// UnmarshalAZERBinField is required for conformance
// with azer.BinFieldUnmarshalable.
func (id *TerminalID) UnmarshalAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (readLen int, err error) {
	i, readLen, err := TerminalIDFromAZERBinField(b, typeHint)
	if err == nil {
		*id = i
	}
	return readLen, err
}

// Equals is required as TerminalID is a value-object.
//
// Use EqualsTerminalID method if the other value
// has the same type.
func (id TerminalID) Equals(other interface{}) bool {
	if x, ok := other.(TerminalID); ok {
		return x == id
	}
	if x, _ := other.(*TerminalID); x != nil {
		return *x == id
	}
	return false
}

// Equal is a wrapper for Equals method. It is required for
// compatibility with github.com/google/go-cmp
func (id TerminalID) Equal(other interface{}) bool {
	return id.Equals(other)
}

// EqualsTerminalID determines if the other instance
// is equal to this instance.
func (id TerminalID) EqualsTerminalID(
	other TerminalID,
) bool {
	return id == other
}

func TerminalIDFromString(s string) (TerminalID, error) {
	if s == "" {
		return TerminalIDZero, nil
	}
	tid, err := terminalIDDecode(s)
	if err != nil {
		return TerminalIDZero, err
	}
	if tid.IsNotValid() {
		return TerminalIDZero, errors.Msg("unexpeted")
	}
	return tid, nil
}

func (id TerminalID) String() string {
	if id.IsNotValid() {
		return ""
	}
	return terminalIDEncode(id)
}

func (id TerminalID) ClientID() ClientID {
	return ClientID(int64(id) >> terminalClientIDShift)
}

func (id TerminalID) InstanceID() int32 {
	return int32(id & terminalInstanceIDMask)
}

func (id TerminalID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

func (id *TerminalID) UnmarshalText(b []byte) error {
	i, err := TerminalIDFromString(string(b))
	if err == nil {
		*id = i
	}
	return err
}

func (id TerminalID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + id.String() + `"`), nil
}

func (id *TerminalID) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	i, err := TerminalIDFromString(s)
	if err == nil {
		*id = i
	}
	return err
}

const (
	terminalInstanceIDMask = 0x00000000ffffffff
	terminalIDMax          = 0x7fffffffffffffff
	terminalClientIDShift  = 32
)

var (
	terminalIDEncodingOnce sync.Once

	terminalIDEncode func(TerminalID) string          = terminalIDV1Encode
	terminalIDDecode func(string) (TerminalID, error) = terminalIDV1Decode
)

func UseTerminalIDV0Enconding() {
	terminalIDEncodingOnce.Do(func() {
		terminalIDEncode = terminalIDV0Encode
		terminalIDDecode = terminalIDV0Decode
	})
}

const (
	terminalIDV1Prefix = "TZZ0T"
)

func terminalIDV1Encode(tid TerminalID) string {
	return terminalIDV1Prefix + crock32.Encode(uint64(tid))
}

func terminalIDV1Decode(s string) (TerminalID, error) {
	if len(s) <= len(terminalIDV1Prefix) {
		return TerminalIDZero, errors.Arg("", errors.Ent("length", nil))
	}
	pfx := s[:len(terminalIDV1Prefix)]
	if pfx != terminalIDV1Prefix {
		return TerminalIDZero, errors.Arg("", errors.Ent("prefix", nil))
	}
	instIDStr := s[len(pfx):]
	instIDU64, err := crock32.Decode(instIDStr)
	if err != nil {
		return TerminalIDZero, errors.Arg("", err)
	}
	if instIDU64 > terminalIDMax {
		return TerminalIDZero, errors.ArgMsg("", "overflow")
	}
	return TerminalID(instIDU64), nil
}

const (
	terminalIDV0Prefix = "tl-0x"
)

func terminalIDV0Encode(tid TerminalID) string {
	return fmt.Sprintf("%s%016x", terminalIDV0Prefix, int64(tid))
}

func terminalIDV0Decode(s string) (TerminalID, error) {
	s = strings.TrimPrefix(s, terminalIDV0Prefix)
	i, err := strconv.ParseInt(s, 16, 64)
	return TerminalID(i), err
}

//endregion

//region RefKey

// TerminalRefKey is used to identify
// an instance of adjunct entity Terminal system-wide.
type TerminalRefKey struct {
	application ApplicationRefKey
	user        UserRefKey
	id          TerminalID
}

// The total number of fields in the struct.
const _TerminalRefKeyFieldCount = 2 + 1

// NewTerminalRefKey returns a new instance
// of TerminalRefKey with the provided attribute values.
func NewTerminalRefKey(
	application ApplicationRefKey,
	user UserRefKey,
	id TerminalID,
) TerminalRefKey {
	return TerminalRefKey{
		application: application,
		user:        user,
		id:          id,
	}
}

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azcore.RefKey = _TerminalRefKeyZero
var _ azcore.AdjunctEntityRefKey = _TerminalRefKeyZero
var _ azcore.TerminalRefKey = _TerminalRefKeyZero

var _TerminalRefKeyZero = TerminalRefKey{
	application: ApplicationRefKeyZero(),
	user:        UserRefKeyZero(),
	id:          TerminalIDZero,
}

// TerminalRefKeyZero returns
// a zero-valued instance of TerminalRefKey.
func TerminalRefKeyZero() TerminalRefKey {
	return _TerminalRefKeyZero
}

// AZRefKey is required by azcore.RefKey interface.
func (TerminalRefKey) AZRefKey() {}

// AZAdjunctEntityRefKey is required
// by azcore.AdjunctEntityRefKey interface.
func (TerminalRefKey) AZAdjunctEntityRefKey() {}

// ID returns the scoped identifier of the entity.
func (refKey TerminalRefKey) ID() TerminalID {
	return refKey.id
}

// IDPtr returns a pointer to a copy of the ID if it's considered valid
// otherwise it returns nil.
func (refKey TerminalRefKey) IDPtr() *TerminalID {
	if refKey.IsNotValid() {
		return nil
	}
	i := refKey.ID()
	return &i
}

// ID is required for conformance with azcore.RefKey.
func (refKey TerminalRefKey) EID() azcore.EID {
	return refKey.id
}

// TerminalID is required for conformance
// with azcore.TerminalRefKey.
func (refKey TerminalRefKey) TerminalID() azcore.TerminalID {
	return refKey.id
}

// IsZero is required as TerminalRefKey is a value-object.
func (refKey TerminalRefKey) IsZero() bool {
	return refKey.application.IsZero() &&
		refKey.user.IsZero() &&
		refKey.id == TerminalIDZero
}

// IsValid returns true if this instance is valid independently as a ref-key.
// It doesn't tell whether it refers to a valid instance of Terminal.
func (refKey TerminalRefKey) IsValid() bool {
	return refKey.application.IsValid() &&
		refKey.user.IsValid() &&
		refKey.id.IsValid()
}

// IsNotValid returns the negation of value returned by IsValid().
func (refKey TerminalRefKey) IsNotValid() bool {
	return !refKey.IsValid()
}

// Equals is required for conformance with azcore.AdjunctEntityRefKey.
func (refKey TerminalRefKey) Equals(other interface{}) bool {
	if x, ok := other.(TerminalRefKey); ok {
		return refKey.application.EqualsApplicationRefKey(x.application) &&
			refKey.user.EqualsUserRefKey(x.user) &&
			refKey.id == x.id
	}
	if x, _ := other.(*TerminalRefKey); x != nil {
		return refKey.application.EqualsApplicationRefKey(x.application) &&
			refKey.user.EqualsUserRefKey(x.user) &&
			refKey.id == x.id
	}
	return false
}

// Equal is required for conformance with azcore.AdjunctEntityRefKey.
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
		refKey.id == other.id
}

// AZERBin is required for conformance
// with azcore.RefKey.
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
// with azcore.RefKey.
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

	fieldBytes, fieldType = refKey.id.AZERBinField()
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
			errors.ArgWrap("", "id type parsing", err)
	}
	typeCursor++
	id, readLen, err := TerminalIDFromAZERBinField(
		b[dataCursor:], fieldType)
	if err != nil {
		return TerminalRefKeyZero(), 0,
			errors.ArgWrap("", "id data parsing", err)
	}
	dataCursor += readLen

	return TerminalRefKey{
		application: applicationRefKey,
		user:        userRefKey,
		id:          id,
	}, dataCursor, nil
}

// UnmarshalAZERBinField is required for conformance
// with azcore.BinFieldUnmarshalable.
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
// with azcore.RefKey.
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

// WithApplication returns a copy
// of TerminalRefKey
// with its application attribute set to the provided value.
func (refKey TerminalRefKey) WithApplication(
	application ApplicationRefKey,
) TerminalRefKey {
	return TerminalRefKey{
		application: application,
		user:        refKey.user,
	}
}

// User returns instance's User value.
func (refKey TerminalRefKey) User() UserRefKey {
	return refKey.user
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
	}
}

// TerminalRefKeyError defines an interface for all
// TerminalRefKey-related errors.
type TerminalRefKeyError interface {
	error
	TerminalRefKeyError()
}

//endregion
