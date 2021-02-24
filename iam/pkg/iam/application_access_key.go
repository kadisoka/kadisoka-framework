package iam

import (
	"bytes"
	"encoding/hex"
	"strconv"
	"strings"

	azcore "github.com/alloyzeus/go-azcore/azcore"
	errors "github.com/alloyzeus/go-azcore/azcore/errors"
	protowire "google.golang.org/protobuf/encoding/protowire"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = azcore.AZCorePackageIsVersion1
var _ = bytes.MinRead
var _ = hex.ErrLength
var _ = strconv.IntSize
var _ = strings.Compare
var _ = protowire.MinValidNumber

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
var _ azcore.AZWireUnmarshalable = &_ApplicationAccessKeyIDZeroVar

// In binary: 0b11111111111111111111111111111111111111111111111111111111
const _ApplicationAccessKeyIDSignificantBitsMask uint64 = 0xffffffffffffff

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

func ApplicationAccessKeyIDFromAZWire(b []byte) (id ApplicationAccessKeyID, readLen int, err error) {
	_, typ, n := protowire.ConsumeTag(b)
	if n <= 0 {
		return ApplicationAccessKeyIDZero, n, ApplicationAccessKeyIDAZWireDecodingArgumentError{}
	}
	readLen = n
	if typ != protowire.VarintType {
		return ApplicationAccessKeyIDZero, readLen, ApplicationAccessKeyIDAZWireDecodingArgumentError{}
	}
	e, n := protowire.ConsumeVarint(b)
	if n <= 0 {
		return ApplicationAccessKeyIDZero, readLen, ApplicationAccessKeyIDAZWireDecodingArgumentError{}
	}
	readLen += n
	return ApplicationAccessKeyID(e), readLen, nil
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

// AZWire returns a binary representation of the instance.
//
// AZWire is required for conformance
// with azcore.AZWireObject.
func (id ApplicationAccessKeyID) AZWire() []byte {
	return id.AZWireField(1)
}

// AZWireField encode this instance as azwire with a specified field number.
//
// AZWire is required for conformance
// with azcore.AZWireObject.
func (id ApplicationAccessKeyID) AZWireField(fieldNum int) []byte {
	var buf []byte
	buf = protowire.AppendTag(buf, protowire.Number(fieldNum), protowire.VarintType)
	buf = protowire.AppendVarint(buf, uint64(id))
	return buf
}

// UnmarshalAZWire is required for conformance
// with azcore.AZWireUnmarshalable.
func (id *ApplicationAccessKeyID) UnmarshalAZWire(b []byte) (readLen int, err error) {
	var i ApplicationAccessKeyID
	i, readLen, err = ApplicationAccessKeyIDFromAZWire(b)
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

type ApplicationAccessKeyIDAZWireDecodingArgumentError struct{}

var _ errors.ArgumentError = ApplicationAccessKeyIDAZWireDecodingArgumentError{}

func (ApplicationAccessKeyIDAZWireDecodingArgumentError) ArgumentName() string {
	return ""
}

func (ApplicationAccessKeyIDAZWireDecodingArgumentError) Error() string {
	return "ApplicationAccessKeyIDAZWireDecodingArgumentError"
}

//endregion

//region RefKey

// ApplicationAccessKeyRefKey is used to identify
// an instance of adjunct entity ApplicationAccessKey system-wide.
type ApplicationAccessKeyRefKey struct {
	application ApplicationRefKey
	id          ApplicationAccessKeyID
}

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
var _ azcore.AZWireUnmarshalable = &_ApplicationAccessKeyRefKeyZero
var _ azcore.AZRSUnmarshalable = &_ApplicationAccessKeyRefKeyZero

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

// AZWire is required for conformance
// with azcore.AZWireObject.
func (refKey ApplicationAccessKeyRefKey) AZWire() []byte {
	return refKey.AZWireField(1)
}

// AZWireField is required for conformance
// with azcore.AZWireObject.
func (refKey ApplicationAccessKeyRefKey) AZWireField(fieldNum int) []byte {
	buf := &bytes.Buffer{}
	var fieldWire []byte
	fieldWire = refKey.application.AZWireField(0 + 1)
	buf.Write(fieldWire)
	fieldWire = refKey.id.AZWireField(1 + 1)
	buf.Write(fieldWire)
	var outBuf []byte
	outBuf = protowire.AppendTag(outBuf,
		protowire.Number(fieldNum), protowire.BytesType)
	outBuf = protowire.AppendBytes(outBuf, buf.Bytes())
	return outBuf
}

// ApplicationAccessKeyRefKeyFromAZWire creates ApplicationAccessKeyRefKey from
// its azwire-encoded form.
func ApplicationAccessKeyRefKeyFromAZWire(
	b []byte,
) (
	refKey ApplicationAccessKeyRefKey,
	readLen int,
	err error,
) {
	var readOffset int = 0
	_, typ, n := protowire.ConsumeTag(b)
	if n <= 0 {
		return ApplicationAccessKeyRefKeyZero(), readOffset, ApplicationAccessKeyRefKeyAZWireDecodingArgumentError{}
	}
	readOffset += n
	if typ != protowire.BytesType {
		return ApplicationAccessKeyRefKeyZero(), readOffset, ApplicationAccessKeyRefKeyAZWireDecodingArgumentError{}
	}
	_, n = protowire.ConsumeVarint(b[readOffset:])
	if n <= 0 {
		return ApplicationAccessKeyRefKeyZero(), readOffset, ApplicationAccessKeyRefKeyAZWireDecodingArgumentError{}
	}
	readOffset += n

	application, fieldLen, err := ApplicationRefKeyFromAZWire(b[readOffset:])
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(), readOffset, ApplicationAccessKeyRefKeyAZWireDecodingArgumentError{}
	}
	readOffset += fieldLen

	id, fieldLen, err := ApplicationAccessKeyIDFromAZWire(b[readOffset:])
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(), readOffset, ApplicationAccessKeyRefKeyAZWireDecodingArgumentError{}
	}
	readOffset += fieldLen

	return ApplicationAccessKeyRefKey{
		application: application,
		id:          id,
	}, readOffset, nil
}

func (refKey *ApplicationAccessKeyRefKey) UnmarshalAZWire(b []byte) (readLen int, err error) {
	var r ApplicationAccessKeyRefKey
	r, readLen, err = ApplicationAccessKeyRefKeyFromAZWire(b)
	if err == nil {
		*refKey = r
	}
	return readLen, err
}

const _ApplicationAccessKeyRefKeyAZRSPrefix = "KAK0"

// ApplicationAccessKeyRefKeyFromAZRS creates ApplicationAccessKeyRefKey from
// its AZRS-encoded form.
func ApplicationAccessKeyRefKeyFromAZRS(s string) (ApplicationAccessKeyRefKey, error) {
	if s == "" {
		return ApplicationAccessKeyRefKeyZero(), nil
	}
	if !strings.HasPrefix(s, _ApplicationAccessKeyRefKeyAZRSPrefix) {
		return ApplicationAccessKeyRefKeyZero(), ApplicationAccessKeyRefKeyAZRSDecodingArgumentError{}
	}
	s = strings.TrimPrefix(s, _ApplicationAccessKeyRefKeyAZRSPrefix)
	b, err := hex.DecodeString(s)
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(), ApplicationAccessKeyRefKeyAZRSDecodingArgumentError{}
	}
	refKey, _, err := ApplicationAccessKeyRefKeyFromAZWire(b)
	if err != nil {
		return ApplicationAccessKeyRefKeyZero(), ApplicationAccessKeyRefKeyAZRSDecodingArgumentError{}
	}
	return refKey, nil
}

// AZRS returns an encoded representation of this instance. It returns empty
// if IsValid returned false.
//
// AZRS is required for conformance
// with azcore.RefKey.
func (refKey ApplicationAccessKeyRefKey) AZRS() string {
	if !refKey.IsValid() {
		return ""
	}
	wire := refKey.AZWire()
	//TODO: configurable encoding
	return _ApplicationAccessKeyRefKeyAZRSPrefix +
		hex.EncodeToString(wire)
}

// UnmarshalAZRS is required for conformance
// with azcore.AZRSUnmarshalable.
func (refKey *ApplicationAccessKeyRefKey) UnmarshalAZRS(s string) error {
	r, err := ApplicationAccessKeyRefKeyFromAZRS(s)
	if err == nil {
		*refKey = r
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

type ApplicationAccessKeyRefKeyError interface {
	error
	ApplicationAccessKeyRefKeyError()
}

type ApplicationAccessKeyRefKeyAZWireDecodingArgumentError struct{}

var _ ApplicationAccessKeyRefKeyError = ApplicationAccessKeyRefKeyAZWireDecodingArgumentError{}
var _ errors.ArgumentError = ApplicationAccessKeyRefKeyAZWireDecodingArgumentError{}

func (ApplicationAccessKeyRefKeyAZWireDecodingArgumentError) ApplicationAccessKeyRefKeyError() {}
func (ApplicationAccessKeyRefKeyAZWireDecodingArgumentError) ArgumentName() string             { return "" }

func (ApplicationAccessKeyRefKeyAZWireDecodingArgumentError) Error() string {
	return "ApplicationAccessKeyRefKeyAZWireDecodingArgumentError"
}

type ApplicationAccessKeyRefKeyAZRSDecodingArgumentError struct{}

var _ ApplicationAccessKeyRefKeyError = ApplicationAccessKeyRefKeyAZRSDecodingArgumentError{}
var _ errors.ArgumentError = ApplicationAccessKeyRefKeyAZRSDecodingArgumentError{}

func (ApplicationAccessKeyRefKeyAZRSDecodingArgumentError) ApplicationAccessKeyRefKeyError() {}
func (ApplicationAccessKeyRefKeyAZRSDecodingArgumentError) ArgumentName() string             { return "" }

func (ApplicationAccessKeyRefKeyAZRSDecodingArgumentError) Error() string {
	return "ApplicationAccessKeyRefKeyAZRSDecodingArgumentError"
}

//endregion
