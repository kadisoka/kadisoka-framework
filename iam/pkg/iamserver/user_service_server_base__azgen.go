package iamserver

import (
	"crypto/rand"
	"encoding/binary"

	errors "github.com/alloyzeus/go-azfl/azfl/errors"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

// UserServiceServerbase is the server-side
// base implementation for UserService.
type UserServiceServerBase struct {
}

//var _ iam.UserService = &UserServiceServer{}

// GenerateUserIDNum generates a new iam.UserIDNum.
// Note that this function does not consulting any database nor registry.
// This method will not create an instance of iam.User, i.e., the
// resulting iam.UserIDNum might or might not refer to valid instance
// of iam.User. The resulting iam.UserIDNum is designed to be
// used as an argument to create a new instance of iam.User.
//
// The embeddedFieldBits argument could be constructed by combining
// iam.UserIDNum*Bits constants. If none are defined,
// use the value of 0.
func GenerateUserIDNum(embeddedFieldBits uint64) (iam.UserIDNum, error) {
	idBytes := make([]byte, 8)
	_, err := rand.Read(idBytes)
	if err != nil {
		return iam.UserIDNumZero, errors.ArgWrap("", "random source reading", err)
	}

	idUint := (embeddedFieldBits & iam.UserIDNumEmbeddedFieldsMask) |
		(binary.BigEndian.Uint64(idBytes) & iam.UserIDNumIdentifierBitsMask)
	return iam.UserIDNum(idUint), nil
}
