package iam

import (
	"encoding/hex"
	"strconv"
	"strings"

	azcore "github.com/alloyzeus/go-azcore/azcore"
	errors "github.com/alloyzeus/go-azcore/azcore/errors"
	protowire "google.golang.org/protobuf/encoding/protowire"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = azcore.AZCorePackageIsVersion1
var _ = hex.ErrLength
var _ = strconv.IntSize
var _ = strings.Compare
var _ = protowire.MinValidNumber

// Entity Application.
//
// An Application is an ....

//region ID

// ApplicationID is a scoped identifier
// used to identify an instance of entity Application.
type ApplicationID int32

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azcore.EID = ApplicationIDZero
var _ azcore.EntityID = ApplicationIDZero
var _ azcore.AZWireUnmarshalable = &_ApplicationIDZeroVar

// In binary: 0b11111111111111111111111111
const _ApplicationIDSignificantBitsMask uint32 = 0x3ffffff

// ApplicationIDZero is the zero value
// for ApplicationID.
const ApplicationIDZero = ApplicationID(0)

// _ApplicationIDZeroVar is used for testing
// pointer-based interfaces conformance.
var _ApplicationIDZeroVar = ApplicationIDZero

// ApplicationIDFromPrimitiveValue creates an instance
// of ApplicationID from its primitive value.
func ApplicationIDFromPrimitiveValue(v int32) ApplicationID {
	return ApplicationID(v)
}

// ApplicationIDFromAZWire creates ApplicationID from
// its azwire-encoded form.
func ApplicationIDFromAZWire(b []byte) (id ApplicationID, readLen int, err error) {
	_, typ, n := protowire.ConsumeTag(b)
	if n <= 0 {
		return ApplicationIDZero, n, ApplicationIDAZWireDecodingArgumentError{}
	}
	readLen = n
	if typ != protowire.VarintType {
		return ApplicationIDZero, readLen, ApplicationIDAZWireDecodingArgumentError{}
	}
	e, n := protowire.ConsumeVarint(b[readLen:])
	if n <= 0 {
		return ApplicationIDZero, readLen, ApplicationIDAZWireDecodingArgumentError{}
	}
	readLen += n
	return ApplicationID(e), readLen, nil
}

// PrimitiveValue returns the ID in its primitive type. Prefer to use
// this method instead of casting directly.
func (id ApplicationID) PrimitiveValue() int32 {
	return int32(id)
}

// AZEID is required for conformance
// with azcore.EID.
func (ApplicationID) AZEID() {}

// AZEntityID is required for conformance
// with azcore.EntityID.
func (ApplicationID) AZEntityID() {}

// IsZero is required as ApplicationID is a value-object.
func (id ApplicationID) IsZero() bool {
	return id == ApplicationIDZero
}

// IsValid returns true if this instance is valid independently as an ID.
// It doesn't tell whether it refers to a valid instance of Application.
func (id ApplicationID) IsValid() bool {
	return int32(id) > 0 &&
		(uint32(id)&_ApplicationIDSignificantBitsMask) != 0
}

// Equals is required as ApplicationID is a value-object.
//
// Use EqualsApplicationID method if the other value
// has the same type.
func (id ApplicationID) Equals(other interface{}) bool {
	if x, ok := other.(ApplicationID); ok {
		return x == id
	}
	if x, _ := other.(*ApplicationID); x != nil {
		return *x == id
	}
	return false
}

// Equal is a wrapper for Equals method. It is required for
// compatibility with github.com/google/go-cmp
func (id ApplicationID) Equal(other interface{}) bool {
	return id.Equals(other)
}

// EqualsApplicationID determines if the other instance is equal
// to this instance.
func (id ApplicationID) EqualsApplicationID(
	other ApplicationID,
) bool {
	return id == other
}

// AZWire returns a binary representation of the instance.
//
// AZWire is required for conformance
// with azcore.AZWireObject.
func (id ApplicationID) AZWire() []byte {
	return id.AZWireField(1)
}

// AZWireField encode this instance as azwire with a specified field number.
//
// AZWire is required for conformance
// with azcore.AZWireObject.
func (id ApplicationID) AZWireField(fieldNum int) []byte {
	var buf []byte
	buf = protowire.AppendTag(buf, protowire.Number(fieldNum), protowire.VarintType)
	buf = protowire.AppendVarint(buf, uint64(id))
	return buf
}

// UnmarshalAZWire is required for conformance
// with azcore.AZWireUnmarshalable.
func (id *ApplicationID) UnmarshalAZWire(b []byte) (readLen int, err error) {
	var i ApplicationID
	i, readLen, err = ApplicationIDFromAZWire(b)
	if err == nil {
		*id = i
	}
	return readLen, err
}

// IsFirstParty returns true if the Application instance
// this ID is for is a FirstParty Application.
//
// FirstParty indicates that the application is a first-party
// application, i.e., an application that we official support.
//
// First-party applications, in contrast to third party
// applications, are those of official applications, officially
// supported applications, or internal-use only applications.
//
// First-party applications could be in various forms, e.g.,
// official mobile app, web app, or system dashboard.
func (id ApplicationID) IsFirstParty() bool {
	const mask = uint32(0) |
		(uint32(1) << 30)
	const flags = uint32(0) |
		(uint32(1) << 30)
	return (uint32(id) & mask) == flags
}

// IsService returns true if the Application instance
// this ID is for is a Service Application.
//
// Service application is an application which does not
// represent user. Note that this is different from, e.g.,
// bot, where it's a specialization of User; a bot is a User.
//
// All service applications are confidential as defined by
// OAuth 2.0.
func (id ApplicationID) IsService() bool {
	const mask = uint32(0) |
		(uint32(1) << 29) |
		(uint32(1) << 28)
	const flags = uint32(0) |
		(uint32(0) << 29) |
		(uint32(1) << 28)
	return (uint32(id) & mask) == flags
}

// IsUserAgent returns true if the Application instance
// this ID is for is a UserAgent Application.
//
// A user-agent application is an application which represents
// the user it got authorized by.
//
// There are two types of user agent: public and confidential.
// These types align with OAuth 2.0's.
func (id ApplicationID) IsUserAgent() bool {
	const mask = uint32(0) |
		(uint32(1) << 29)
	const flags = uint32(0) |
		(uint32(1) << 29)
	return (uint32(id) & mask) == flags
}

// IsPublicAuthorizationUserAgent returns true if the Application instance
// this ID is for is a PublicAuthorizationUserAgent Application.
//
// A direct user-agent application is an application which
// can be used by users to authenticate themselves to the system.
// The application will automatically receive authorization upon
// successful user authentication.
//
// A direct user-agent credentials should never be assumed to
// be secure, i.e., authorized direct user-agent application
// has no strong identity. Access control checks should be
// focused on testing the user's claims and less of the
// application's claims.
func (id ApplicationID) IsPublicAuthorizationUserAgent() bool {
	const mask = uint32(0) |
		(uint32(1) << 29) |
		(uint32(1) << 28)
	const flags = uint32(0) |
		(uint32(1) << 29) |
		(uint32(0) << 28)
	return (uint32(id) & mask) == flags
}

// IsConfidentialAuthorizationUserAgent returns true if the Application instance
// this ID is for is a ConfidentialAuthorizationUserAgent Application.
//
// A confidential user-agent application is a user-agent
// application which could receive authorization from a
// consenting user through 3-legged authorization flow. A
// confidential user-agent can not be used for directly
// performing user authentication.
func (id ApplicationID) IsConfidentialAuthorizationUserAgent() bool {
	const mask = uint32(0) |
		(uint32(1) << 29) |
		(uint32(1) << 28)
	const flags = uint32(0) |
		(uint32(1) << 29) |
		(uint32(1) << 28)
	return (uint32(id) & mask) == flags
}

type ApplicationIDError interface {
	error
	ApplicationIDError()
}

type ApplicationIDAZWireDecodingArgumentError struct{}

var _ ApplicationIDError = ApplicationIDAZWireDecodingArgumentError{}
var _ errors.ArgumentError = ApplicationIDAZWireDecodingArgumentError{}

func (ApplicationIDAZWireDecodingArgumentError) ApplicationIDError()  {}
func (ApplicationIDAZWireDecodingArgumentError) ArgumentName() string { return "" }

func (ApplicationIDAZWireDecodingArgumentError) Error() string {
	return "ApplicationIDAZWireDecodingArgumentError"
}

//TODO: FromString, (Un)MarshalText, (Un)MarshalJSON

//endregion

//region RefKey

// ApplicationRefKey is used to identify
// an instance of entity Application system-wide.
type ApplicationRefKey ApplicationID

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azcore.RefKey = _ApplicationRefKeyZero
var _ azcore.EntityRefKey = _ApplicationRefKeyZero
var _ azcore.AZWireUnmarshalable = &_ApplicationRefKeyZeroVar
var _ azcore.AZISUnmarshalable = &_ApplicationRefKeyZeroVar

const _ApplicationRefKeyZero = ApplicationRefKey(ApplicationIDZero)

var _ApplicationRefKeyZeroVar = _ApplicationRefKeyZero

// ApplicationRefKeyZero returns
// a zero-valued instance of ApplicationRefKey.
func ApplicationRefKeyZero() ApplicationRefKey {
	return _ApplicationRefKeyZero
}

// AZRefKey is required for conformance with azcore.RefKey.
func (ApplicationRefKey) AZRefKey() {}

// AZEntityRefKey is required for conformance
// with azcore.EntityRefKey.
func (ApplicationRefKey) AZEntityRefKey() {}

// ID is required for conformance with azcore.RefKey.
func (refKey ApplicationRefKey) ID() azcore.EID {
	return ApplicationID(refKey)
}

// IsZero is required as ApplicationRefKey is a value-object.
func (refKey ApplicationRefKey) IsZero() bool {
	return ApplicationID(refKey) == ApplicationIDZero
}

// IsValid returns true if this instance is valid independently as a ref-key.
// It doesn't tell whether it refers to a valid instance of Application.
func (refKey ApplicationRefKey) IsValid() bool {
	return ApplicationID(refKey).IsValid()
}

// Equals is required for conformance with azcore.EntityRefKey.
func (refKey ApplicationRefKey) Equals(other interface{}) bool {
	if x, ok := other.(ApplicationRefKey); ok {
		return x == refKey
	}
	if x, _ := other.(*ApplicationRefKey); x != nil {
		return *x == refKey
	}
	return false
}

// Equal is required for conformance with azcore.EntityRefKey.
func (refKey ApplicationRefKey) Equal(other interface{}) bool {
	return refKey.Equals(other)
}

// EqualsApplicationRefKey returs true
// if the other value has the same attributes as refKey.
func (refKey ApplicationRefKey) EqualsApplicationRefKey(
	other ApplicationRefKey,
) bool {
	return other == refKey
}

// AZWire is required for conformance
// with azcore.AZWireObject.
func (refKey ApplicationRefKey) AZWire() []byte {
	return refKey.AZWireField(1)
}

// AZWireField is required for conformance
// with azcore.AZWireObject.
func (refKey ApplicationRefKey) AZWireField(fieldNum int) []byte {
	return ApplicationID(refKey).AZWireField(fieldNum)
}

// ApplicationRefKeyFromAZWire creates ApplicationRefKey from
// its azwire-encoded form.
func ApplicationRefKeyFromAZWire(b []byte) (refKey ApplicationRefKey, readLen int, err error) {
	var id ApplicationID
	id, readLen, err = ApplicationIDFromAZWire(b)
	if err != nil {
		return ApplicationRefKeyZero(), readLen, ApplicationRefKeyAZWireDecodingArgumentError{}
	}
	return ApplicationRefKey(id), readLen, nil
}

// UnmarshalAZWire is required for conformance
// with azcore.AZWireUnmarshalable.
func (refKey *ApplicationRefKey) UnmarshalAZWire(b []byte) (readLen int, err error) {
	var i ApplicationRefKey
	i, readLen, err = ApplicationRefKeyFromAZWire(b)
	if err == nil {
		*refKey = i
	}
	return readLen, err
}

const _ApplicationRefKeyAZISPrefix = "KAp0"

// ApplicationRefKeyFromAZIS creates ApplicationRefKey from
// its AZIS-encoded form.
func ApplicationRefKeyFromAZIS(s string) (ApplicationRefKey, error) {
	if s == "" {
		return ApplicationRefKeyZero(), nil
	}
	if !strings.HasPrefix(s, _ApplicationRefKeyAZISPrefix) {
		return ApplicationRefKeyZero(), ApplicationRefKeyAZISDecodingArgumentError{}
	}
	s = strings.TrimPrefix(s, _ApplicationRefKeyAZISPrefix)
	b, err := hex.DecodeString(s)
	if err != nil {
		return ApplicationRefKeyZero(), ApplicationRefKeyAZISDecodingArgumentError{}
	}
	refKey, _, err := ApplicationRefKeyFromAZWire(b)
	if err != nil {
		return ApplicationRefKeyZero(), ApplicationRefKeyAZISDecodingArgumentError{}
	}
	return refKey, nil
}

// AZIS returns an encoded representation of this instance. It returns empty
// if IsValid returned false.
//
// AZIS is required for conformance
// with azcore.RefKey.
func (refKey ApplicationRefKey) AZIS() string {
	if !refKey.IsValid() {
		return ""
	}
	wire := refKey.AZWire()
	//TODO: configurable encoding
	return _ApplicationRefKeyAZISPrefix +
		hex.EncodeToString(wire)
}

// UnmarshalAZIS is required for conformance
// with azcore.AZISUnmarshalable.
func (refKey *ApplicationRefKey) UnmarshalAZIS(s string) error {
	r, err := ApplicationRefKeyFromAZIS(s)
	if err == nil {
		*refKey = r
	}
	return err
}

type ApplicationRefKeyError interface {
	error
	ApplicationRefKeyError()
}

type ApplicationRefKeyAZWireDecodingArgumentError struct{}

var _ ApplicationRefKeyError = ApplicationRefKeyAZWireDecodingArgumentError{}
var _ errors.ArgumentError = ApplicationRefKeyAZWireDecodingArgumentError{}

func (ApplicationRefKeyAZWireDecodingArgumentError) ApplicationRefKeyError() {}
func (ApplicationRefKeyAZWireDecodingArgumentError) ArgumentName() string    { return "" }

func (ApplicationRefKeyAZWireDecodingArgumentError) Error() string {
	return "ApplicationRefKeyAZWireDecodingArgumentError"
}

type ApplicationRefKeyAZISDecodingArgumentError struct{}

var _ ApplicationRefKeyError = ApplicationRefKeyAZISDecodingArgumentError{}
var _ errors.ArgumentError = ApplicationRefKeyAZISDecodingArgumentError{}

func (ApplicationRefKeyAZISDecodingArgumentError) ApplicationRefKeyError() {}
func (ApplicationRefKeyAZISDecodingArgumentError) ArgumentName() string    { return "" }

func (ApplicationRefKeyAZISDecodingArgumentError) Error() string {
	return "ApplicationRefKeyAZISDecodingArgumentError"
}

//endregion

//

// func ApplicationRefKeyFromString(s string) (ApplicationRefKey, error) {
// 	if s == "" {
// 		return ApplicationRefKeyZero(), nil
// 	}
// }

// func ApplicationIDFromString(s string) (ApplicationID, error) {
// 	if s == "" {
// 		return ApplicationIDZero, nil
// 	}
// 	cid, err := clientIDDecode(s)
// 	if err != nil {
// 		return ApplicationIDZero, err
// 	}
// 	if cid.IsNotValid() {
// 		return ApplicationIDZero, errors.Msg("unexpected")
// 	}
// 	return cid, nil
// }

// func (clientID ApplicationID) String() string {
// 	if clientID.IsNotValid() {
// 		return ""
// 	}
// 	return clientIDEncode(clientID)
// }

// func (clientID ApplicationID) IsValid() bool {
// 	return clientID.validVersion() &&
// 		clientID.validInstance() &&
// 		clientID.validType()
// }
// func (clientID ApplicationID) IsNotValid() bool {
// 	return !clientID.IsValid()
// }
