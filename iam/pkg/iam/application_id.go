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

// Entity Application.
//
// An Application is an ....

//region IDNum

// ApplicationIDNum is a scoped identifier
// used to identify an instance of entity Application.
type ApplicationIDNum int32

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azfl.IDNum = ApplicationIDNumZero
var _ azfl.EntityID = ApplicationIDNumZero
var _ azer.BinFieldUnmarshalable = &_ApplicationIDNumZeroVar

// ApplicationIDNumSignificantBitsMask is used to
// extract significant bits from an instance of ApplicationIDNum.
const ApplicationIDNumSignificantBitsMask uint32 = 0b11_11111111_11111111_11111111

// ApplicationIDNumZero is the zero value
// for ApplicationIDNum.
const ApplicationIDNumZero = ApplicationIDNum(0)

// _ApplicationIDNumZeroVar is used for testing
// pointer-based interfaces conformance.
var _ApplicationIDNumZeroVar = ApplicationIDNumZero

// ApplicationIDNumFromPrimitiveValue creates an instance
// of ApplicationIDNum from its primitive value.
func ApplicationIDNumFromPrimitiveValue(v int32) ApplicationIDNum {
	return ApplicationIDNum(v)
}

// ApplicationIDNumFromAZERBinField creates ApplicationIDNum from
// its azer-bin-field form.
func ApplicationIDNumFromAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (idNum ApplicationIDNum, readLen int, err error) {
	if typeHint != azer.BinDataTypeUnspecified && typeHint != azer.BinDataTypeInt32 {
		return ApplicationIDNum(0), 0,
			errors.ArgMsg("typeHint", "unsupported")
	}
	i := binary.BigEndian.Uint32(b)
	return ApplicationIDNum(i), 4, nil
}

// PrimitiveValue returns the value in its primitive type. Prefer to use
// this method instead of casting directly.
func (idNum ApplicationIDNum) PrimitiveValue() int32 {
	return int32(idNum)
}

// AZIDNum is required for conformance
// with azfl.IDNum.
func (ApplicationIDNum) AZIDNum() {}

// AZEntityID is required for conformance
// with azfl.EntityID.
func (ApplicationIDNum) AZEntityID() {}

// IsZero is required as ApplicationIDNum is a value-object.
func (idNum ApplicationIDNum) IsZero() bool {
	return idNum == ApplicationIDNumZero
}

// IsValid returns true if this instance is valid independently
// as an ApplicationIDNum. It doesn't tell whether it refers to
// a valid instance of Application.
//
// To elaborate, validity of a data depends on the perspective of the user.
// For example, age 1000 is a valid as an instance of age, but on the context
// of human living age, we can consider it as invalid.
//
// To use some analogy, a ticket has a date of validity for today, but
// after it got checked in to the counter, it turns out that its serial number
// is not registered in the issuer's database. The ticket claims that it's
// valid, but it's considered invalid because it's a fake.
//
// Similarly, what is considered valid in this context here is that the data
// contained in this instance doesn't break any rule for an instance of
// ApplicationIDNum. Whether the instance is valid for certain context,
// it requires case-by-case validation which is out of the scope of this
// method.
func (idNum ApplicationIDNum) IsValid() bool {
	return int32(idNum) > 0 &&
		(uint32(idNum)&ApplicationIDNumSignificantBitsMask) != 0
}

// IsNotValid returns the negation of value returned by IsValid().
func (idNum ApplicationIDNum) IsNotValid() bool {
	return !idNum.IsValid()
}

// Equals is required as ApplicationIDNum is a value-object.
//
// Use EqualsApplicationIDNum method if the other value
// has the same type.
func (idNum ApplicationIDNum) Equals(other interface{}) bool {
	if x, ok := other.(ApplicationIDNum); ok {
		return x == idNum
	}
	if x, _ := other.(*ApplicationIDNum); x != nil {
		return *x == idNum
	}
	return false
}

// Equal is a wrapper for Equals method. It is required for
// compatibility with github.com/google/go-cmp
func (idNum ApplicationIDNum) Equal(other interface{}) bool {
	return idNum.Equals(other)
}

// EqualsApplicationIDNum determines if the other instance is equal
// to this instance.
func (idNum ApplicationIDNum) EqualsApplicationIDNum(
	other ApplicationIDNum,
) bool {
	return idNum == other
}

// AZERBinField is required for conformance
// with azfl.IDNum.
func (idNum ApplicationIDNum) AZERBinField() ([]byte, azer.BinDataType) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(idNum))
	return b, azer.BinDataTypeInt32
}

// UnmarshalAZERBinField is required for conformance
// with azer.BinFieldUnmarshalable.
func (idNum *ApplicationIDNum) UnmarshalAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (readLen int, err error) {
	i, readLen, err := ApplicationIDNumFromAZERBinField(b, typeHint)
	if err == nil {
		*idNum = i
	}
	return readLen, err
}

// IsFirstParty returns true if
// the Application instance this ApplicationIDNum is for
// is a FirstParty Application.
//
// FirstParty indicates that the application is a first-party
// application, i.e., the application was provided by the realm.
//
// First-party applications, in contrast to third party
// applications, are those of official applications, officially
// supported public applications, and internal-use only
// applications.
//
// First-party applications could be in various forms, e.g.,
// official mobile app, web app, or system dashboard.
func (idNum ApplicationIDNum) IsFirstParty() bool {
	return idNum.IsValid() && idNum.HasFirstPartyBits()
}

const _ApplicationIDNumFirstPartyMask = 0b1000000_00000000_00000000_00000000
const _ApplicationIDNumFirstPartyBits = 0b1000000_00000000_00000000_00000000

// HasFirstPartyBits is only checking the bits
// without validating other information.
func (idNum ApplicationIDNum) HasFirstPartyBits() bool {
	return (uint32(idNum) &
		_ApplicationIDNumFirstPartyMask) ==
		_ApplicationIDNumFirstPartyBits
}

// IsService returns true if
// the Application instance this ApplicationIDNum is for
// is a Service Application.
//
// A service application is an application which does not
// represent user. Note that this is different from, e.g.,
// bot, where it's a specialization of User; a bot is a User.
//
// All service applications are confidential OAuth 2.0
// clients (RFC6749 section 2.1).
func (idNum ApplicationIDNum) IsService() bool {
	return idNum.IsValid() && idNum.HasServiceBits()
}

const _ApplicationIDNumServiceMask = 0b100000_00000000_00000000_00000000
const _ApplicationIDNumServiceBits = 0b0

// HasServiceBits is only checking the bits
// without validating other information.
func (idNum ApplicationIDNum) HasServiceBits() bool {
	return (uint32(idNum) &
		_ApplicationIDNumServiceMask) ==
		_ApplicationIDNumServiceBits
}

// IsUserAgent returns true if
// the Application instance this ApplicationIDNum is for
// is a UserAgent Application.
//
// A user-agent application is an application which represents
// the user it got authorized by.
//
// There are two types of user agent based on the flow they use
// to obtain authorization: public and confidential.
// These types align with OAuth 2.0 client types (RFC6749
// section 2.1).
func (idNum ApplicationIDNum) IsUserAgent() bool {
	return idNum.IsValid() && idNum.HasUserAgentBits()
}

const _ApplicationIDNumUserAgentMask = 0b100000_00000000_00000000_00000000
const _ApplicationIDNumUserAgentBits = 0b100000_00000000_00000000_00000000

// HasUserAgentBits is only checking the bits
// without validating other information.
func (idNum ApplicationIDNum) HasUserAgentBits() bool {
	return (uint32(idNum) &
		_ApplicationIDNumUserAgentMask) ==
		_ApplicationIDNumUserAgentBits
}

// IsUserAgentAuthorizationPublic returns true if
// the Application instance this ApplicationIDNum is for
// is a UserAgentAuthorizationPublic Application.
//
// A public-authorization user-agent application is an
// application which can be used by users to authenticate
// themselves to the system. The application will automatically
// receive authorization upon successful authentication of
// the user.
//
// A public-authorization user-agent credentials should
// never be assumed to be secure, i.e., it can not securely
// authenticate itself. An authorized application of this
// kind has no strong identity. Access control checks should
// be focused on testing the user's claims and less of the
// application's claims.
func (idNum ApplicationIDNum) IsUserAgentAuthorizationPublic() bool {
	return idNum.IsValid() && idNum.HasUserAgentAuthorizationPublicBits()
}

const _ApplicationIDNumUserAgentAuthorizationPublicMask = 0b110000_00000000_00000000_00000000
const _ApplicationIDNumUserAgentAuthorizationPublicBits = 0b100000_00000000_00000000_00000000

// HasUserAgentAuthorizationPublicBits is only checking the bits
// without validating other information.
func (idNum ApplicationIDNum) HasUserAgentAuthorizationPublicBits() bool {
	return (uint32(idNum) &
		_ApplicationIDNumUserAgentAuthorizationPublicMask) ==
		_ApplicationIDNumUserAgentAuthorizationPublicBits
}

// IsUserAgentAuthorizationConfidential returns true if
// the Application instance this ApplicationIDNum is for
// is a UserAgentAuthorizationConfidential Application.
//
// A confidential-authorization user-agent application
// is a user-agent application which could receive authorization
// from a consenting user through 3-legged authorization flow.
// A confidential user-agent can not be used for directly
// performing user authentication.
func (idNum ApplicationIDNum) IsUserAgentAuthorizationConfidential() bool {
	return idNum.IsValid() && idNum.HasUserAgentAuthorizationConfidentialBits()
}

const _ApplicationIDNumUserAgentAuthorizationConfidentialMask = 0b110000_00000000_00000000_00000000
const _ApplicationIDNumUserAgentAuthorizationConfidentialBits = 0b110000_00000000_00000000_00000000

// HasUserAgentAuthorizationConfidentialBits is only checking the bits
// without validating other information.
func (idNum ApplicationIDNum) HasUserAgentAuthorizationConfidentialBits() bool {
	return (uint32(idNum) &
		_ApplicationIDNumUserAgentAuthorizationConfidentialMask) ==
		_ApplicationIDNumUserAgentAuthorizationConfidentialBits
}

type ApplicationIDNumError interface {
	error
	ApplicationIDNumError()
}

//TODO: (Un)MarshalText (for SQL?)

//endregion

//region RefKey

// ApplicationRefKey is used to identify
// an instance of entity Application system-wide.
type ApplicationRefKey ApplicationIDNum

// NewApplicationRefKey returns a new instance
// of ApplicationRefKey with the provided attribute values.
func NewApplicationRefKey(
	idNum ApplicationIDNum,
) ApplicationRefKey {
	return ApplicationRefKey(idNum)
}

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azfl.RefKey = _ApplicationRefKeyZero
var _ azfl.EntityRefKey = _ApplicationRefKeyZero

const _ApplicationRefKeyZero = ApplicationRefKey(ApplicationIDNumZero)

var _ApplicationRefKeyZeroVar = _ApplicationRefKeyZero

// ApplicationRefKeyZero returns
// a zero-valued instance of ApplicationRefKey.
func ApplicationRefKeyZero() ApplicationRefKey {
	return _ApplicationRefKeyZero
}

// AZRefKey is required for conformance with azfl.RefKey.
func (ApplicationRefKey) AZRefKey() {}

// AZEntityRefKey is required for conformance
// with azfl.EntityRefKey.
func (ApplicationRefKey) AZEntityRefKey() {}

// IDNum returns the scoped identifier of the entity.
func (refKey ApplicationRefKey) IDNum() ApplicationIDNum {
	return ApplicationIDNum(refKey)
}

// IDNumPtr returns a pointer to a copy of the IDNum if it's considered valid
// otherwise it returns nil.
func (refKey ApplicationRefKey) IDNumPtr() *ApplicationIDNum {
	if refKey.IsNotValid() {
		return nil
	}
	i := refKey.IDNum()
	return &i
}

// IDNum is required for conformance with azfl.RefKey.
func (refKey ApplicationRefKey) AZIDNum() azfl.IDNum {
	return ApplicationIDNum(refKey)
}

// IsZero is required as ApplicationRefKey is a value-object.
func (refKey ApplicationRefKey) IsZero() bool {
	return ApplicationIDNum(refKey) == ApplicationIDNumZero
}

// IsValid returns true if this instance is valid independently as a ref-key.
// It doesn't tell whether it refers to a valid instance of Application.
func (refKey ApplicationRefKey) IsValid() bool {
	return ApplicationIDNum(refKey).IsValid()
}

// IsNotValid returns the negation of value returned by IsValid().
func (refKey ApplicationRefKey) IsNotValid() bool {
	return !refKey.IsValid()
}

// Equals is required for conformance with azfl.EntityRefKey.
func (refKey ApplicationRefKey) Equals(other interface{}) bool {
	if x, ok := other.(ApplicationRefKey); ok {
		return x == refKey
	}
	if x, _ := other.(*ApplicationRefKey); x != nil {
		return *x == refKey
	}
	return false
}

// Equal is required for conformance with azfl.EntityRefKey.
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

func (refKey ApplicationRefKey) AZERBin() []byte {
	b := make([]byte, 4+1)
	b[0] = azer.BinDataTypeInt32.Byte()
	binary.BigEndian.PutUint32(b[1:], uint32(refKey))
	return b
}

func ApplicationRefKeyFromAZERBin(b []byte) (refKey ApplicationRefKey, readLen int, err error) {
	typ, err := azer.BinDataTypeFromByte(b[0])
	if err != nil {
		return _ApplicationRefKeyZero, 0,
			errors.ArgWrap("", "type parsing", err)
	}
	if typ != azer.BinDataTypeInt32 {
		return _ApplicationRefKeyZero, 0,
			errors.Arg("", errors.EntMsg("type", "unsupported"))
	}

	i, readLen, err := ApplicationRefKeyFromAZERBinField(b[1:], typ)
	if err != nil {
		return _ApplicationRefKeyZero, 0,
			errors.ArgWrap("", "idnum data parsing", err)
	}

	return ApplicationRefKey(i), 1 + readLen, nil
}

// UnmarshalAZERBin is required for conformance
// with azfl.BinFieldUnmarshalable.
func (refKey *ApplicationRefKey) UnmarshalAZERBin(b []byte) (readLen int, err error) {
	i, readLen, err := ApplicationRefKeyFromAZERBin(b)
	if err == nil {
		*refKey = i
	}
	return readLen, err
}

func (refKey ApplicationRefKey) AZERBinField() ([]byte, azer.BinDataType) {
	return ApplicationIDNum(refKey).AZERBinField()
}

func ApplicationRefKeyFromAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (refKey ApplicationRefKey, readLen int, err error) {
	idNum, n, err := ApplicationIDNumFromAZERBinField(b, typeHint)
	if err != nil {
		return _ApplicationRefKeyZero, n, err
	}
	return ApplicationRefKey(idNum), n, nil
}

// UnmarshalAZERBinField is required for conformance
// with azfl.BinFieldUnmarshalable.
func (refKey *ApplicationRefKey) UnmarshalAZERBinField(
	b []byte, typeHint azer.BinDataType,
) (readLen int, err error) {
	i, readLen, err := ApplicationRefKeyFromAZERBinField(b, typeHint)
	if err == nil {
		*refKey = i
	}
	return readLen, err
}

const _ApplicationRefKeyAZERTextPrefix = "KAp0"

// AZERText is required for conformance
// with azfl.RefKey.
func (refKey ApplicationRefKey) AZERText() string {
	if !refKey.IsValid() {
		return ""
	}

	return _ApplicationRefKeyAZERTextPrefix +
		azer.TextEncode(refKey.AZERBin())
}

// ApplicationRefKeyFromAZERText creates a new instance of
// ApplicationRefKey from its azer-text form.
func ApplicationRefKeyFromAZERText(s string) (ApplicationRefKey, error) {
	if s == "" {
		return ApplicationRefKeyZero(), nil
	}
	if !strings.HasPrefix(s, _ApplicationRefKeyAZERTextPrefix) {
		return ApplicationRefKeyZero(),
			errors.Arg("", errors.EntMsg("prefix", "mismatch"))
	}
	s = strings.TrimPrefix(s, _ApplicationRefKeyAZERTextPrefix)
	b, err := azer.TextDecode(s)
	if err != nil {
		return ApplicationRefKeyZero(),
			errors.ArgWrap("", "data parsing", err)
	}
	refKey, _, err := ApplicationRefKeyFromAZERBin(b)
	if err != nil {
		return ApplicationRefKeyZero(),
			errors.ArgWrap("", "data decoding", err)
	}
	return refKey, nil
}

// UnmarshalAZERText is required for conformance
// with azer.TextUnmarshalable.
func (refKey *ApplicationRefKey) UnmarshalAZERText(s string) error {
	r, err := ApplicationRefKeyFromAZERText(s)
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalText is for compatibility with Go's encoding.TextMarshaler
func (refKey ApplicationRefKey) MarshalText() ([]byte, error) {
	return []byte(refKey.AZERText()), nil
}

// UnmarshalText is for conformance with Go's encoding.TextUnmarshaler
func (refKey *ApplicationRefKey) UnmarshalText(b []byte) error {
	r, err := ApplicationRefKeyFromAZERText(string(b))
	if err == nil {
		*refKey = r
	}
	return err
}

// MarshalJSON makes this type JSON-marshalable.
func (refKey ApplicationRefKey) MarshalJSON() ([]byte, error) {
	// We assume that there's no symbols in azer-text
	return []byte("\"" + refKey.AZERText() + "\""), nil
}

// UnmarshalJSON parses a JSON value.
func (refKey *ApplicationRefKey) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		*refKey = ApplicationRefKeyZero()
		return nil
	}
	i, err := ApplicationRefKeyFromAZERText(s)
	if err == nil {
		*refKey = i
	}
	return err
}

// ApplicationRefKeyService abstracts
// ApplicationRefKey-related services.
type ApplicationRefKeyService interface {
	// IsApplicationRefKey is to check if the ref-key is
	// trully registered to system. It does not check whether the instance
	// is active or not.
	IsApplicationRefKeyRegistered(refKey ApplicationRefKey) bool
}

// ApplicationRefKeyError defines an interface for all
// ApplicationRefKey-related errors.
type ApplicationRefKeyError interface {
	error
	ApplicationRefKeyError()
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
