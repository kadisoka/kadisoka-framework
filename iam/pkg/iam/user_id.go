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

var (
	ErrUserIDStringInvalid        = errors.Ent("user ID string", nil)
	ErrServiceUserIDStringInvalid = errors.Ent("service user ID string", nil)
)

//region ID

// UserID is a scoped identifier
// used to identify an instance of entity User.
type UserID int64

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azcore.EID = UserIDZero
var _ azcore.EntityID = UserIDZero
var _ azcore.AZWireUnmarshalable = &_UserIDZeroVar
var _ azcore.UserID = UserIDZero

// UserIDZero is the zero value
// for UserID.
const UserIDZero = UserID(0)

// _UserIDZeroVar is used for testing
// pointer-based interfaces conformance.
var _UserIDZeroVar = UserIDZero

// UserIDFromPrimitiveValue creates an instance
// of UserID from its primitive value.
func UserIDFromPrimitiveValue(v int64) UserID {
	return UserID(v)
}

// UserIDFromAZWire creates UserID from
// its azwire-encoded form.
func UserIDFromAZWire(b []byte) (id UserID, readLen int, err error) {
	_, typ, n := protowire.ConsumeTag(b)
	if n <= 0 {
		return UserIDZero, n, UserIDAZWireDecodingArgumentError{}
	}
	readLen = n
	if typ != protowire.VarintType {
		return UserIDZero, readLen, UserIDAZWireDecodingArgumentError{}
	}
	e, n := protowire.ConsumeVarint(b)
	if n <= 0 {
		return UserIDZero, readLen, UserIDAZWireDecodingArgumentError{}
	}
	readLen += n
	return UserID(e), readLen, nil
}

// PrimitiveValue returns the ID in its primitive type. Prefer to use
// this method instead of casting directly.
func (id UserID) PrimitiveValue() int64 {
	return int64(id)
}

// AZEID is required for conformance
// with azcore.EID.
func (UserID) AZEID() {}

// AZEntityID is required for conformance
// with azcore.EntityID.
func (UserID) AZEntityID() {}

// AZUserID is required for conformance
// with azcore.UserID.
func (UserID) AZUserID() {}

// IsZero is required as UserID is a value-object.
func (id UserID) IsZero() bool {
	return id == UserIDZero
}

// Equals is required as UserID is a value-object.
//
// Use EqualsUserID method if the other value
// has the same type.
func (id UserID) Equals(other interface{}) bool {
	if x, ok := other.(UserID); ok {
		return x == id
	}
	if x, _ := other.(*UserID); x != nil {
		return *x == id
	}
	return false
}

// Equal is a wrapper for Equals method. It is required for
// compatibility with github.com/google/go-cmp
func (id UserID) Equal(other interface{}) bool {
	return id.Equals(other)
}

// EqualsUserID determines if the other instance is equal
// to this instance.
func (id UserID) EqualsUserID(
	other UserID,
) bool {
	return id == other
}

// AZWire returns a binary representation of the instance.
//
// AZWire is required for conformance
// with azcore.AZWireObject.
func (id UserID) AZWire() []byte {
	return id.AZWireField(1)
}

// AZWireField encode this instance as azwire with a specified field number.
//
// AZWire is required for conformance
// with azcore.AZWireObject.
func (id UserID) AZWireField(fieldNum int) []byte {
	var buf []byte
	buf = protowire.AppendTag(buf, protowire.Number(fieldNum), protowire.VarintType)
	buf = protowire.AppendVarint(buf, uint64(id))
	return buf
}

// UnmarshalAZWire is required for conformance
// with azcore.AZWireUnmarshalable.
func (id *UserID) UnmarshalAZWire(b []byte) (readLen int, err error) {
	var i UserID
	i, readLen, err = UserIDFromAZWire(b)
	if err == nil {
		*id = i
	}
	return readLen, err
}

// IsBot returns true if the User instance
// this ID is for is a Bot User.
//
// Bot account is ....
func (id UserID) IsBot() bool {
	const mask = uint64(0) |
		(uint64(1) << 61)
	const flags = uint64(0) |
		(uint64(1) << 61)
	return (uint64(id) & mask) == flags
}

type UserIDError interface {
	error
	UserIDError()
}

type UserIDAZWireDecodingArgumentError struct{}

var _ UserIDError = UserIDAZWireDecodingArgumentError{}
var _ errors.ArgumentError = UserIDAZWireDecodingArgumentError{}

func (UserIDAZWireDecodingArgumentError) UserIDError()         {}
func (UserIDAZWireDecodingArgumentError) ArgumentName() string { return "" }

func (UserIDAZWireDecodingArgumentError) Error() string {
	return "UserIDAZWireDecodingArgumentError"
}

type UserIDWireDecodingArgumentError struct{}

var _ errors.ArgumentError = UserIDWireDecodingArgumentError{}

func (UserIDWireDecodingArgumentError) ArgumentName() string {
	return ""
}

func (UserIDWireDecodingArgumentError) Error() string {
	return "UserIDWireDecodingArgumentError"
}

func UserIDFromString(s string) (UserID, error) {
	if s == "" {
		return UserIDZero, nil
	}
	return userIDDecode(s)
}

func (id UserID) IsValid() bool    { return id > userIDReservedMax && id <= userIDMax }
func (id UserID) IsNotValid() bool { return !id.IsValid() }

func (id UserID) String() string {
	if id.IsNotValid() {
		return ""
	}
	return userIDEncode(id)
}

func (id UserID) IsNormalAccount() bool {
	return id.IsValid() && id > userIDServiceMax
}

func (id UserID) IsServiceAccount() bool {
	return id.IsValid() && id <= userIDServiceMax
}

func (id UserID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

func (id *UserID) UnmarshalText(b []byte) error {
	i, err := UserIDFromString(string(b))
	if err == nil {
		*id = i
	}
	return err
}

func (id UserID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + id.String() + `"`), nil
}

func (id *UserID) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		*id = UserIDZero
		return nil
	}
	i, err := UserIDFromString(s)
	if err == nil {
		*id = i
	}
	return err
}

var (
	userIDEncodingOnce sync.Once

	userIDMax         UserID = userIDV1Max
	userIDServiceMax  UserID = userIDV1ServiceMax
	userIDReservedMax UserID = userIDV1ReservedMax

	userIDEncode func(UserID) string          = userIDV1Encode
	userIDDecode func(string) (UserID, error) = userIDV1Decode
)

func UseUserIDV0Enconding() {
	userIDEncodingOnce.Do(func() {
		userIDMax = userIDV0Max
		userIDServiceMax = userIDV0ServiceMax
		userIDReservedMax = userIDV0ReservedMax
		userIDEncode = userIDV0Encode
		userIDDecode = userIDV0Decode
	})
}

const (
	userIDV1Max           = 0x0000ffffffffffff
	userIDV1ReservedMax   = 0x000000000000ffff
	userIDV1ServiceMax    = 0x00000000ffffffff
	userIDV1Prefix        = "INo0T"
	userIDV1ServicePrefix = "ISv0T"
)

func userIDV1Encode(userID UserID) string {
	var prefix string
	if userID.IsServiceAccount() {
		prefix = userIDV1ServicePrefix
	} else {
		prefix = userIDV1Prefix
	}
	return prefix + crock32.Encode(uint64(userID))
}

func userIDV1Decode(s string) (UserID, error) {
	isService := strings.HasPrefix(s, userIDV1ServicePrefix)
	if isService {
		s = strings.TrimPrefix(s, userIDV1ServicePrefix)
	} else {
		s = strings.TrimPrefix(s, userIDV1Prefix)
	}

	i, err := crock32.Decode(s)
	if err != nil {
		return UserIDZero, errors.Arg("", err)
	}
	// To ensure we can safely treat it as signed
	if i > uint64(0x7fffffffffffffff) {
		return UserIDZero, errors.ArgMsg("", "overflow")
	}

	if isService {
		if i > userIDV1ServiceMax {
			return UserIDZero, errors.Arg("", nil)
		}
	} else {
		if i != 0 && i <= userIDV1ServiceMax {
			return UserIDZero, errors.Arg("", nil)
		}
	}

	return UserID(i), nil
}

const (
	userIDV0Max = 0x0000ffffffffffff

	// userIDV0ReservedMax is maximum value for reserved user IDs. IDs within
	// this range should never be considered as valid user IDs in client
	// applications.
	userIDV0ReservedMax = 0x00000000000fffff

	// userIDV0ServiceMax is a constant which we use to separate service user IDs
	// and normal user IDs.
	//
	// We are reserving user IDs up to this value. We will use these user ID for
	// various purpose in the future. Possible usage: service applications, bots,
	// service notifications.
	userIDV0ServiceMax = 0x00000003ffffffff

	// userIDV0ServicePrefix is a prefix we use to differentiate normal
	// user (human-representing) account and service user account.
	userIDV0ServicePrefix = "is-0x"

	// userIDV0Prefix is the prefix for normal users.
	userIDV0Prefix = "i-0x"

	userIDV0EncodingRadix = 16
)

func userIDV0Encode(userID UserID) string {
	var prefix string
	if userID.IsServiceAccount() {
		prefix = userIDV0ServicePrefix
	} else {
		prefix = userIDV0Prefix
	}
	return prefix + fmt.Sprintf("%016x", userID.PrimitiveValue())
}

func userIDV0Decode(s string) (UserID, error) {
	isService := strings.HasPrefix(s, userIDV0ServicePrefix)
	if isService {
		s = strings.TrimPrefix(s, userIDV0ServicePrefix)
	} else {
		s = strings.TrimPrefix(s, userIDV0Prefix)
	}

	i, err := strconv.ParseInt(s, userIDV0EncodingRadix, 64)
	if err != nil {
		return UserIDZero, errors.Arg("", err)
	}

	if isService {
		if i > userIDV0ServiceMax {
			return UserIDZero, errors.Arg("", nil)
		}
	} else {
		if i != 0 && i <= userIDV0ServiceMax {
			return UserIDZero, errors.Arg("", nil)
		}
	}

	return UserID(i), nil
}

//endregion

//region RefKey

// UserRefKey is used to identify
// an instance of entity User system-wide.
type UserRefKey UserID

// To ensure that it conforms the interfaces. If any of these is failing,
// there's a bug in the generator.
var _ azcore.RefKey = _UserRefKeyZero
var _ azcore.EntityRefKey = _UserRefKeyZero
var _ azcore.AZWireUnmarshalable = &_UserRefKeyZeroVar
var _ azcore.AZRSUnmarshalable = &_UserRefKeyZeroVar
var _ azcore.UserRefKey = _UserRefKeyZero

const _UserRefKeyZero = UserRefKey(UserIDZero)

var _UserRefKeyZeroVar = _UserRefKeyZero

// UserRefKeyZero returns
// a zero-valued instance of UserRefKey.
func UserRefKeyZero() UserRefKey {
	return _UserRefKeyZero
}

// AZRefKey is required for conformance with azcore.RefKey.
func (UserRefKey) AZRefKey() {}

// AZEntityRefKey is required for conformance
// with azcore.EntityRefKey.
func (UserRefKey) AZEntityRefKey() {}

// ID is required for conformance with azcore.RefKey.
func (refKey UserRefKey) ID() azcore.EID {
	return UserID(refKey)
}

// UserID is required for conformance with azcore.UserRefKey.
func (refKey UserRefKey) UserID() azcore.UserID {
	return UserID(refKey)
}

// IsZero is required as UserRefKey is a value-object.
func (refKey UserRefKey) IsZero() bool {
	return UserID(refKey) == UserIDZero
}

// Equals is required for conformance with azcore.EntityRefKey.
func (refKey UserRefKey) Equals(other interface{}) bool {
	if x, ok := other.(UserRefKey); ok {
		return x == refKey
	}
	if x, _ := other.(*UserRefKey); x != nil {
		return *x == refKey
	}
	return false
}

// Equal is required for conformance with azcore.EntityRefKey.
func (refKey UserRefKey) Equal(other interface{}) bool {
	return refKey.Equals(other)
}

// EqualsUserRefKey returs true
// if the other value has the same attributes as refKey.
func (refKey UserRefKey) EqualsUserRefKey(
	other UserRefKey,
) bool {
	return other == refKey
}

// AZWire is required for conformance
// with azcore.AZWireObject.
func (refKey UserRefKey) AZWire() []byte {
	return refKey.AZWireField(1)
}

// AZWireField is required for conformance
// with azcore.AZWireObject.
func (refKey UserRefKey) AZWireField(fieldNum int) []byte {
	return UserID(refKey).AZWireField(fieldNum)
}

// UserRefKeyFromAZWire creates UserRefKey from
// its azwire-encoded form.
func UserRefKeyFromAZWire(b []byte) (refKey UserRefKey, readLen int, err error) {
	var id UserID
	id, readLen, err = UserIDFromAZWire(b)
	if err != nil {
		return UserRefKeyZero(), readLen, UserRefKeyAZWireDecodingArgumentError{}
	}
	return UserRefKey(id), readLen, nil
}

// UnmarshalAZWire is required for conformance
// with azcore.AZWireUnmarshalable.
func (refKey *UserRefKey) UnmarshalAZWire(b []byte) (readLen int, err error) {
	var i UserRefKey
	i, readLen, err = UserRefKeyFromAZWire(b)
	if err == nil {
		*refKey = i
	}
	return readLen, err
}

const _UserRefKeyAZRSPrefix = "KUs0"

// UserRefKeyFromAZRS creates UserRefKey from
// its AZRS-encoded form.
func UserRefKeyFromAZRS(s string) (UserRefKey, error) {
	if !strings.HasPrefix(s, _UserRefKeyAZRSPrefix) {
		return UserRefKeyZero(), UserRefKeyAZRSDecodingArgumentError{}
	}
	s = strings.TrimPrefix(s, _UserRefKeyAZRSPrefix)
	b, err := hex.DecodeString(s)
	if err != nil {
		return UserRefKeyZero(), UserRefKeyAZRSDecodingArgumentError{}
	}
	refKey, _, err := UserRefKeyFromAZWire(b)
	if err != nil {
		return UserRefKeyZero(), UserRefKeyAZRSDecodingArgumentError{}
	}
	return refKey, nil
}

// AZRS returns an encoded representation of this instance.
//
// AZRS is required for conformance
// with azcore.RefKey.
func (refKey UserRefKey) AZRS() string {
	wire := refKey.AZWire()
	//TODO: configurable encoding
	return _UserRefKeyAZRSPrefix +
		hex.EncodeToString(wire)
}

// UnmarshalAZRS is required for conformance
// with azcore.AZRSUnmarshalable.
func (refKey *UserRefKey) UnmarshalAZRS(s string) error {
	r, err := UserRefKeyFromAZRS(s)
	if err == nil {
		*refKey = r
	}
	return err
}

type UserRefKeyError interface {
	error
	UserRefKeyError()
}

type UserRefKeyAZWireDecodingArgumentError struct{}

var _ UserRefKeyError = UserRefKeyAZWireDecodingArgumentError{}
var _ errors.ArgumentError = UserRefKeyAZWireDecodingArgumentError{}

func (UserRefKeyAZWireDecodingArgumentError) UserRefKeyError()     {}
func (UserRefKeyAZWireDecodingArgumentError) ArgumentName() string { return "" }

func (UserRefKeyAZWireDecodingArgumentError) Error() string {
	return "UserRefKeyAZWireDecodingArgumentError"
}

type UserRefKeyAZRSDecodingArgumentError struct{}

var _ UserRefKeyError = UserRefKeyAZRSDecodingArgumentError{}
var _ errors.ArgumentError = UserRefKeyAZRSDecodingArgumentError{}

func (UserRefKeyAZRSDecodingArgumentError) UserRefKeyError()     {}
func (UserRefKeyAZRSDecodingArgumentError) ArgumentName() string { return "" }

func (UserRefKeyAZRSDecodingArgumentError) Error() string {
	return "UserRefKeyAZRSDecodingArgumentError"
}

//endregion
