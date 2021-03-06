package iamserver

import (
	"crypto/rand"
	"encoding/binary"

	errors "github.com/alloyzeus/go-azfl/azfl/errors"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

const sessionDBTableName = "session_dt"

// GenerateSessionIDNum generates a new iam.SessionIDNum.
// Note that this function does not consult any database nor registry.
// This method will not create an instance of iam.Session, i.e., the
// resulting iam.SessionIDNum might or might not refer to valid instance
// of iam.Session. The resulting iam.SessionIDNum is designed to be
// used as an argument to create a new instance of iam.Session.
//
// The embeddedFieldBits argument could be constructed by combining
// iam.SessionIDNum*Bits constants. If none are defined,
// use the value of 0.
func GenerateSessionIDNum(embeddedFieldBits uint32) (iam.SessionIDNum, error) {
	idBytes := make([]byte, 4)
	_, err := rand.Read(idBytes)
	if err != nil {
		return iam.SessionIDNumZero, errors.Wrap("random number source reading", err)
	}

	idUint := (embeddedFieldBits & iam.SessionIDNumEmbeddedFieldsMask) |
		(binary.BigEndian.Uint32(idBytes) & iam.SessionIDNumIdentifierBitsMask)
	return iam.SessionIDNum(idUint), nil
}
