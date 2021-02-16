package iam

import (
	"strconv"

	azcore "github.com/alloyzeus/go-azcore/azcore"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = azcore.AZCorePackageIsVersion1

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

// IDString returns a string representation of this instance.
func (id ApplicationID) IDString() string {
	return id.AZEIDString()
}

// AZEIDString returns a string representation
// of the instance as an EID.
func (id ApplicationID) AZEIDString() string {
	//TODO: custom encoding
	return strconv.FormatInt(int64(id), 10)
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
// confidential user-agent can not be used for performing
// user authentication.
func (id ApplicationID) IsConfidentialAuthorizationUserAgent() bool {
	const mask = uint32(0) |
		(uint32(1) << 29) |
		(uint32(1) << 28)
	const flags = uint32(0) |
		(uint32(1) << 29) |
		(uint32(1) << 28)
	return (uint32(id) & mask) == flags
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
	return ApplicationID(refKey)
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

// RefKeyString returns an encoded representation of this instance.
//
// RefKeyString is required by azcore.RefKey.
func (refKey ApplicationRefKey) RefKeyString() string {
	// TODO:
	// something like /<pluralized type_name>/<id> or
	// /<type_name><separator><id>
	//
	// might need to include version information (actually, use prefix option instead.).
	return "Application(" + ApplicationID(refKey).AZEIDString() + ")"
}

//endregion
