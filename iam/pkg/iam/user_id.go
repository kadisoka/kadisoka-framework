package iam

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	azcore "github.com/alloyzeus/go-azcore/azcore"
	"github.com/richardlehane/crock32"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/errors"
)

var (
	ErrUserIDStringInvalid        = errors.Ent("user ID string", nil)
	ErrServiceUserIDStringInvalid = errors.Ent("service user ID string", nil)
)

//region ID

// UserID is a scoped identifier
// used to identify an instance of entity User.
type UserID int64

var _ azcore.EID = UserIDZero
var _ azcore.EntityID = UserIDZero
var _ azcore.UserID = UserIDZero

// UserIDZero is the zero value
// for UserID.
const UserIDZero = UserID(0)

// UserIDFromPrimitiveValue creates an instance
// of UserID from its primitive value.
func UserIDFromPrimitiveValue(v int64) UserID {
	return UserID(v)
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

// IDString returns a string representation of this instance.
func (id UserID) IDString() string {
	return id.AZEIDString()
}

// AZEIDString returns a string representation
// of the instance as an EID.
func (id UserID) AZEIDString() string {
	//TODO: custom encoding
	return strconv.FormatInt(int64(id), 10)
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

// To ensure that it conforms the interfaces
var _ azcore.RefKey = _UserRefKeyZero
var _ azcore.EntityRefKey = _UserRefKeyZero
var _ azcore.UserRefKey = _UserRefKeyZero

const _UserRefKeyZero = UserRefKey(UserIDZero)

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

// RefKeyString returns an encoded representation of this instance.
//
// RefKeyString is required by azcore.RefKey.
func (refKey UserRefKey) RefKeyString() string {
	// TODO:
	// something like /<pluralized type_name>/<id> or
	// /<type_name><separator><id>
	//
	// might need to include version information (actually, use prefix option instead.).
	return "User(" + UserID(refKey).AZEIDString() + ")"
}

//endregion
