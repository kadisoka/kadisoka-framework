package iam

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"

	azcore "github.com/alloyzeus/go-azcore/azcore"
	"github.com/alloyzeus/go-azcore/azcore/errors"
	"github.com/richardlehane/crock32"
	protowire "google.golang.org/protobuf/encoding/protowire"
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
var _ azcore.AZWireUnmarshalable = &_TerminalIDZeroVar
var _ azcore.TerminalID = TerminalIDZero

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

func TerminalIDFromAZWire(b []byte) (id TerminalID, readLen int, err error) {
	_, typ, n := protowire.ConsumeTag(b)
	if n <= 0 {
		return TerminalIDZero, n, TerminalIDAZWireDecodingArgumentError{}
	}
	readLen = n
	if typ != protowire.VarintType {
		return TerminalIDZero, readLen, TerminalIDAZWireDecodingArgumentError{}
	}
	e, n := protowire.ConsumeVarint(b)
	if n <= 0 {
		return TerminalIDZero, readLen, TerminalIDAZWireDecodingArgumentError{}
	}
	readLen += n
	return TerminalID(e), readLen, nil
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

// AZTerminalID is required for conformance with azcore.TerminalID.
func (TerminalID) AZTerminalID() {}

// AZWire returns a binary representation of the instance.
//
// AZWire is required for conformance
// with azcore.AZWireObject.
func (id TerminalID) AZWire() []byte {
	return id.AZWireField(1)
}

// AZWireField encode this instance as azwire with a specified field number.
//
// AZWire is required for conformance
// with azcore.AZWireObject.
func (id TerminalID) AZWireField(fieldNum int) []byte {
	var buf []byte
	buf = protowire.AppendTag(buf, protowire.Number(fieldNum), protowire.VarintType)
	buf = protowire.AppendVarint(buf, uint64(id))
	return buf
}

// UnmarshalAZWire is required for conformance
// with azcore.AZWireUnmarshalable.
func (id *TerminalID) UnmarshalAZWire(b []byte) (readLen int, err error) {
	var i TerminalID
	i, readLen, err = TerminalIDFromAZWire(b)
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

type TerminalIDAZWireDecodingArgumentError struct{}

var _ errors.ArgumentError = TerminalIDAZWireDecodingArgumentError{}

func (TerminalIDAZWireDecodingArgumentError) ArgumentName() string {
	return ""
}

func (TerminalIDAZWireDecodingArgumentError) Error() string {
	return "TerminalIDAZWireDecodingArgumentError"
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

func (id TerminalID) IsValid() bool {
	return (id&terminalInstanceIDMask) > 0 &&
		id.ClientID().IsValid()
}

func (id TerminalID) IsNotValid() bool {
	return !id.IsValid()
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
	id TerminalID
}

// NewTerminalRefKey returns a new instance
// of TerminalRefKey with the provided attribute values.
func NewTerminalRefKey(
	id TerminalID,
) TerminalRefKey {
	return TerminalRefKey{
		id: id,
	}
}

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azcore.RefKey = _TerminalRefKeyZero
var _ azcore.AdjunctEntityRefKey = _TerminalRefKeyZero
var _ azcore.AZWireUnmarshalable = &_TerminalRefKeyZero
var _ azcore.AZRSUnmarshalable = &_TerminalRefKeyZero
var _ azcore.TerminalRefKey = _TerminalRefKeyZero

var _TerminalRefKeyZero = TerminalRefKey{
	id: TerminalIDZero,
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

// ID is required for conformance with azcore.RefKey.
func (refKey TerminalRefKey) ID() azcore.EID {
	return refKey.id
}

// TerminalID is required for conformance with azcore.TerminalRefKey.
func (refKey TerminalRefKey) TerminalID() azcore.TerminalID {
	return refKey.id
}

// IsZero is required as TerminalRefKey is a value-object.
func (refKey TerminalRefKey) IsZero() bool {
	return refKey.id == TerminalIDZero
}

// Equals is required for conformance with azcore.AdjunctEntityRefKey.
func (refKey TerminalRefKey) Equals(other interface{}) bool {
	if x, ok := other.(TerminalRefKey); ok {
		return refKey.id == x.id
	}
	if x, _ := other.(*TerminalRefKey); x != nil {
		return refKey.id == x.id
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
	return refKey.id == other.id
}

// AZWire is required for conformance
// with azcore.AZWireObject.
func (refKey TerminalRefKey) AZWire() []byte {
	return refKey.AZWireField(1)
}

// AZWireField is required for conformance
// with azcore.AZWireObject.
func (refKey TerminalRefKey) AZWireField(fieldNum int) []byte {
	return refKey.id.AZWireField(fieldNum)
}

// TerminalRefKeyFromAZWire creates TerminalRefKey from
// its azwire-encoded form.
func TerminalRefKeyFromAZWire(
	b []byte,
) (
	refKey TerminalRefKey,
	readLen int,
	err error,
) {
	var readOffset int = 0
	_, typ, n := protowire.ConsumeTag(b)
	if n <= 0 {
		return TerminalRefKeyZero(), readOffset, TerminalRefKeyAZWireDecodingArgumentError{}
	}
	readOffset += n
	if typ != protowire.BytesType {
		return TerminalRefKeyZero(), readOffset, TerminalRefKeyAZWireDecodingArgumentError{}
	}
	_, n = protowire.ConsumeVarint(b[readOffset:])
	if n <= 0 {
		return TerminalRefKeyZero(), readOffset, TerminalRefKeyAZWireDecodingArgumentError{}
	}
	readOffset += n

	id, fieldLen, err := TerminalIDFromAZWire(b[readOffset:])
	if err != nil {
		return TerminalRefKeyZero(), readOffset, TerminalRefKeyAZWireDecodingArgumentError{}
	}
	readOffset += fieldLen

	return TerminalRefKey{
		id: id,
	}, readOffset, nil
}

func (refKey *TerminalRefKey) UnmarshalAZWire(b []byte) (readLen int, err error) {
	var r TerminalRefKey
	r, readLen, err = TerminalRefKeyFromAZWire(b)
	if err == nil {
		*refKey = r
	}
	return readLen, err
}

const _TerminalRefKeyAZRSPrefix = "KTx0"

// TerminalRefKeyFromAZRS creates TerminalRefKey from
// its AZRS-encoded form.
func TerminalRefKeyFromAZRS(s string) (TerminalRefKey, error) {
	if !strings.HasPrefix(s, _TerminalRefKeyAZRSPrefix) {
		return TerminalRefKeyZero(), TerminalRefKeyAZRSDecodingArgumentError{}
	}
	s = strings.TrimPrefix(s, _TerminalRefKeyAZRSPrefix)
	b, err := hex.DecodeString(s)
	if err != nil {
		return TerminalRefKeyZero(), TerminalRefKeyAZRSDecodingArgumentError{}
	}
	refKey, _, err := TerminalRefKeyFromAZWire(b)
	if err != nil {
		return TerminalRefKeyZero(), TerminalRefKeyAZRSDecodingArgumentError{}
	}
	return refKey, nil
}

// AZRS returns an encoded representation of this instance.
//
// AZRS is required for conformance
// with azcore.RefKey.
func (refKey TerminalRefKey) AZRS() string {
	wire := refKey.AZWire()
	//TODO: configurable encoding
	return _TerminalRefKeyAZRSPrefix +
		hex.EncodeToString(wire)
}

// UnmarshalAZRS is required for conformance
// with azcore.AZRSUnmarshalable.
func (refKey *TerminalRefKey) UnmarshalAZRS(s string) error {
	r, err := TerminalRefKeyFromAZRS(s)
	if err == nil {
		*refKey = r
	}
	return err
}

type TerminalRefKeyError interface {
	error
	TerminalRefKeyError()
}

type TerminalRefKeyAZWireDecodingArgumentError struct{}

var _ TerminalRefKeyError = TerminalRefKeyAZWireDecodingArgumentError{}
var _ errors.ArgumentError = TerminalRefKeyAZWireDecodingArgumentError{}

func (TerminalRefKeyAZWireDecodingArgumentError) TerminalRefKeyError() {}
func (TerminalRefKeyAZWireDecodingArgumentError) ArgumentName() string { return "" }

func (TerminalRefKeyAZWireDecodingArgumentError) Error() string {
	return "TerminalRefKeyAZWireDecodingArgumentError"
}

type TerminalRefKeyAZRSDecodingArgumentError struct{}

var _ TerminalRefKeyError = TerminalRefKeyAZRSDecodingArgumentError{}
var _ errors.ArgumentError = TerminalRefKeyAZRSDecodingArgumentError{}

func (TerminalRefKeyAZRSDecodingArgumentError) TerminalRefKeyError() {}
func (TerminalRefKeyAZRSDecodingArgumentError) ArgumentName() string { return "" }

func (TerminalRefKeyAZRSDecodingArgumentError) Error() string {
	return "TerminalRefKeyAZRSDecodingArgumentError"
}

//endregion
