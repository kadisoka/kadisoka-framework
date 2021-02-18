package iam

import (
	"strconv"

	azcore "github.com/alloyzeus/go-azcore/azcore"
	crockford32 "github.com/alloyzeus/go-azcore/azcore/eid/integer/textencodings/crockford32"
	errors "github.com/alloyzeus/go-azcore/azcore/errors"
	protowire "google.golang.org/protobuf/encoding/protowire"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = azcore.AZCorePackageIsVersion1
var _ = strconv.IntSize
var _ = protowire.MinValidNumber

// Entity Application.
//
// An Application is an ....

//region ID

// ApplicationID is a scoped identifier
// used to identify an instance of entity Application.
type ApplicationID int32

var _ azcore.EID = ApplicationIDZero
var _ azcore.EntityID = ApplicationIDZero

// ApplicationIDZero is the zero value
// for ApplicationID.
const ApplicationIDZero = ApplicationID(0)

// ApplicationIDFromPrimitiveValue creates an instance
// of ApplicationID from its primitive value.
func ApplicationIDFromPrimitiveValue(v int32) ApplicationID {
	return ApplicationID(v)
}

func ApplicationIDFromAZWire(b []byte) (ApplicationID, error) {
	_, typ, n := protowire.ConsumeTag(b)
	if n <= 0 {
		return ApplicationIDZero, ApplicationIDWireDecodingArgumentError{}
	}
	if typ != protowire.VarintType {
		return ApplicationIDZero, ApplicationIDWireDecodingArgumentError{}
	}
	e, n := protowire.ConsumeVarint(b)
	if n <= 0 {
		return ApplicationIDZero, ApplicationIDWireDecodingArgumentError{}
	}
	return ApplicationID(e), nil
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
	var buf []byte
	protowire.AppendTag(buf, protowire.Number(1), protowire.VarintType)
	protowire.AppendVarint(buf, uint64(id))
	return buf
}

// UnmarshalAZWire is required for conformance
// with azcore.AZWireObject.
func (id *ApplicationID) UnmarshalAZWire(b []byte) error {
	i, err := ApplicationIDFromAZWire(b)
	if err == nil {
		*id = i
	}
	return err
}

// AZString returns a string representation of the instance.
func (id ApplicationID) AZString() string {
	return "ap-" + crockford32.EncodeInt64(int64(id))
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

type ApplicationIDWireDecodingArgumentError struct{}

var _ errors.ArgumentError = ApplicationIDWireDecodingArgumentError{}

func (ApplicationIDWireDecodingArgumentError) ArgumentName() string {
	return ""
}

func (ApplicationIDWireDecodingArgumentError) Error() string {
	return "ApplicationIDWireDecodingArgumentError"
}

//TODO: FromString, (Un)MarshalText, (Un)MarshalJSON

//endregion

//region RefKey

// ApplicationRefKey is used to identify
// an instance of entity Application system-wide.
type ApplicationRefKey ApplicationID

// To ensure that it conforms the interfaces
var _ azcore.RefKey = _ApplicationRefKeyZero
var _ azcore.EntityRefKey = _ApplicationRefKeyZero

const _ApplicationRefKeyZero = ApplicationRefKey(ApplicationIDZero)

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
	v := ApplicationID(refKey)
	return &v
}

// IsZero is required as ApplicationRefKey is a value-object.
func (refKey ApplicationRefKey) IsZero() bool {
	return ApplicationID(refKey) == ApplicationIDZero
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
	return ApplicationID(refKey).AZWire()
}

// UnmarshalAZWire is required for conformance
// with azcore.AZWireObject.
func (refKey *ApplicationRefKey) UnmarshalAZWire(b []byte) error {
	i, err := ApplicationIDFromAZWire(b)
	if err == nil {
		*refKey = ApplicationRefKey(i)
	}
	return err
}

// AZString returns an encoded representation of this instance.
//
// AZString is required for conformance with azcore.RefKey.
func (refKey ApplicationRefKey) AZString() string {
	return "App(" + ApplicationID(refKey).AZString() + ")"
}

//endregion

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
