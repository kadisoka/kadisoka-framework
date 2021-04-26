package iamserver

import (
	"crypto/rand"
	"encoding/binary"

	errors "github.com/alloyzeus/go-azfl/azfl/errors"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

// ApplicationServiceServerbase is the server-side
// base implementation for ApplicationService.
type ApplicationServiceServerBase struct {
}

//var _ iam.ApplicationService = &ApplicationServiceServer{}

// GenerateApplicationIDNum generates a new iam.ApplicationIDNum.
// Note that this function does not consulting any database nor registry.
// This method will not create an instance of iam.Application, i.e., the
// resulting iam.ApplicationIDNum might or might not refer to valid instance
// of iam.Application. The resulting iam.ApplicationIDNum is designed to be
// used as an argument to create a new instance of iam.Application.
//
// The embeddedFieldBits argument could be constructed by combining
// iam.ApplicationIDNum*Bits constants. If none are defined,
// use the value of 0.
func GenerateApplicationIDNum(embeddedFieldBits uint32) (iam.ApplicationIDNum, error) {
	idBytes := make([]byte, 4)
	_, err := rand.Read(idBytes)
	if err != nil {
		return iam.ApplicationIDNumZero, errors.ArgWrap("", "random source reading", err)
	}

	idUint := (embeddedFieldBits & iam.ApplicationIDNumEmbeddedFieldsMask) |
		(binary.BigEndian.Uint32(idBytes) & iam.ApplicationIDNumIdentifierBitsMask)
	return iam.ApplicationIDNum(idUint), nil
}
