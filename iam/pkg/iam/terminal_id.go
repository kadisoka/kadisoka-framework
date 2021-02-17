package iam

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"sync"

	azcore "github.com/alloyzeus/go-azcore/azcore"
	"github.com/alloyzeus/go-azcore/azcore/errors"
	"github.com/richardlehane/crock32"
)

//region ID

// TerminalID is a scoped identifier
// used to identify an instance of adjunct entity Terminal
// scoped within its host entity(s).
type TerminalID int64

var _ azcore.EID = TerminalIDZero
var _ azcore.AdjunctEntityID = TerminalIDZero
var _ azcore.TerminalID = TerminalIDZero

// TerminalIDZero is the zero value for TerminalID.
const TerminalIDZero = TerminalID(0)

// TerminalIDFromPrimitiveValue creates an instance
// of TerminalID from its primitive value.
func TerminalIDFromPrimitiveValue(v int64) TerminalID {
	return TerminalID(v)
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

// AZEIDBinary returns a binary representation
// of the instance as an EID.
func (id TerminalID) AZEIDBinary() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, uint64(id))
	return buf[:n]
}

// AZEIDString returns a string representation
// of the instance as an EID.
func (id TerminalID) AZEIDString() string {
	return id.AZAdjunctEntityIDString()
}

// AZAdjunctEntityIDString returns a string representation
// of the instance as an AdjunctEntityID.
func (id TerminalID) AZAdjunctEntityIDString() string {
	return "" + strconv.FormatInt(int64(id), 10)
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

// To ensure that it conforms the interfaces
var _ azcore.RefKey = _TerminalRefKeyZero
var _ azcore.AdjunctEntityRefKey = _TerminalRefKeyZero
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

// RefKeyString returns an encoded representation of this instance.
//
// RefKeyString is required by azcore.RefKey.
func (refKey TerminalRefKey) RefKeyString() string {
	// TODO: refkeystring should be defined in the source as it needs
	// to be strictly consistent across implementations.
	// something like /host1_type-host1_id/host2_type-host2_id/hostn_type-hostn_id/own_type-own_id
	// or for global adjuncts /own_type-own_id
	//
	// note that a ref key might comprise of other ref keys. so, we will have
	// something like A(B(C(), D()), E()). the default pattern must be able
	// to accomodate such structure.
	return "Terminal(" +
		refKey.id.AZAdjunctEntityIDString() + ")"
}

//endregion
