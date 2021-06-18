package iamserver

import (
	"crypto/rand"
	"encoding/binary"

	errors "github.com/alloyzeus/go-azfl/azfl/errors"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

const terminalDBTableName = "terminal_dt"

// GenerateTerminalIDNum generates a new iam.TerminalIDNum.
// Note that this function does not consult any database nor registry.
// This method will not create an instance of iam.Terminal, i.e., the
// resulting iam.TerminalIDNum might or might not refer to valid instance
// of iam.Terminal. The resulting iam.TerminalIDNum is designed to be
// used as an argument to create a new instance of iam.Terminal.
//
// The embeddedFieldBits argument could be constructed by combining
// iam.TerminalIDNum*Bits constants. If none are defined,
// use the value of 0.
func GenerateTerminalIDNum(embeddedFieldBits uint64) (iam.TerminalIDNum, error) {
	idBytes := make([]byte, 8)
	_, err := rand.Read(idBytes)
	if err != nil {
		return iam.TerminalIDNumZero, errors.Wrap("random number source reading", err)
	}

	idUint := (embeddedFieldBits & iam.TerminalIDNumEmbeddedFieldsMask) |
		(binary.BigEndian.Uint64(idBytes) & iam.TerminalIDNumIdentifierBitsMask)
	return iam.TerminalIDNum(idUint), nil
}
