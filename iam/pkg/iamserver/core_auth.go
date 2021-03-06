package iamserver

import (
	"crypto/rand"
	"time"

	"github.com/alloyzeus/go-azcore/azcore/errors"
	"github.com/lib/pq"
	"github.com/square/go-jose/v3/jwt"

	apperrs "github.com/kadisoka/kadisoka-framework/foundation/pkg/app/errors"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

const sessionTableName = "session_t"

func (core *Core) GenerateAccessTokenJWT(
	callCtx iam.CallContext,
	terminalRef iam.TerminalRefKey,
	userRef iam.UserRefKey,
	issueTime time.Time,
) (tokenString string, err error) {
	if callCtx == nil {
		return "", errors.ArgMsg("callCtx", "missing")
	}

	jwtKeyChain := core.JWTKeyChain()
	if jwtKeyChain == nil {
		return "", apperrs.NewConfigurationMsg("JWT key chain is not configured")
	}
	signer, err := jwtKeyChain.GetSigner()
	if err != nil {
		return "", errors.Wrap("signer", err)
	}
	if signer == nil {
		return "", apperrs.NewConfigurationMsg("JWT key chain does not have any signing key")
	}

	sessionRef, issueTime, err := core.
		issueSession(callCtx, terminalRef, userRef)
	if err != nil {
		return "", err
	}

	tokenClaims := &iam.AccessTokenClaims{
		Claims: jwt.Claims{
			ID:       sessionRef.AZERText(),
			IssuedAt: jwt.NewNumericDate(issueTime),
			Issuer:   core.RealmName(),
			Expiry:   jwt.NewNumericDate(issueTime.Add(iam.AccessTokenTTLDefault)),
			Subject:  userRef.AZERText(),
		},
		AuthorizedParty: terminalRef.Application().AZERText(),
		TerminalID:      terminalRef.AZERText(),
	}

	tokenString, err = jwt.Signed(signer).Claims(tokenClaims).CompactSerialize()
	if err != nil {
		return "", errors.Wrap("signing", err)
	}
	return
}

func (core *Core) GenerateRefreshTokenJWT(
	callCtx iam.CallContext,
	terminalRef iam.TerminalRefKey,
	terminalSecret string,
	issueTime time.Time,
) (tokenString string, err error) {
	jwtKeyChain := core.JWTKeyChain()
	if jwtKeyChain == nil {
		return "", apperrs.NewConfigurationMsg("JWT key chain is not configured")
	}
	signer, err := jwtKeyChain.GetSigner()
	if err != nil {
		return "", errors.Wrap("signer", err)
	}
	if signer == nil {
		return "", apperrs.NewConfigurationMsg("JWT key chain does not have any signing key")
	}

	tokenClaims := &iam.RefreshTokenClaims{
		NotBefore:      issueTime.Unix(),
		ExpiresAt:      issueTime.Add(iam.RefreshTokenTTLDefault).Unix(),
		TerminalID:     terminalRef.AZERText(),
		TerminalSecret: terminalSecret,
	}

	tokenString, err = jwt.Signed(signer).Claims(tokenClaims).CompactSerialize()
	if err != nil {
		return "", errors.Wrap("signing", err)
	}
	return
}

func (core *Core) issueSession(
	callCtx iam.CallContext,
	terminalRef iam.TerminalRefKey,
	userRef iam.UserRefKey,
) (sessionRef iam.SessionRefKey, issueTime time.Time, err error) {
	authCtx := callCtx.Authorization()

	const attemptNumMax = 3
	timeZero := time.Time{}
	sessionStartTime := timeZero
	var sessionID iam.SessionID

	//TODO: make this more random. using timestamp might cause some security
	// and/or privacy issue.
	// Note:
	// - 0xffffffffffffff00 - timestamp
	// - 0x00000000000000ff - random
	genSessionID := func(ts int64) (iam.SessionID, error) {
		idBytes := make([]byte, 1)
		_, err := rand.Read(idBytes)
		if err != nil {
			return iam.SessionIDZero, errors.Wrap("generation", err)
		}
		bits := ((ts << 8) | int64(idBytes[0])) & int64(iam.SessionIDSignificantBitsMask)
		return iam.SessionID(bits), nil
	}

	for attemptNum := 0; ; attemptNum++ {
		sessionStartTime = time.Now().UTC()
		sessionID, err = genSessionID(sessionStartTime.Unix())
		if err != nil {
			return iam.SessionRefKeyZero(), timeZero, err
		}
		_, err = core.db.
			Exec(
				`INSERT INTO `+sessionTableName+` (`+
					`terminal_id, id, c_ts, c_uid, c_tid`+
					`) VALUES (`+
					`$1, $2, $3, $4, $5`+
					`)`,
				terminalRef.ID().PrimitiveValue(), sessionID.PrimitiveValue(),
				sessionStartTime, authCtx.UserIDPtr(), authCtx.TerminalIDPtr())
		if err == nil {
			break
		}

		pqErr, _ := err.(*pq.Error)
		if pqErr != nil &&
			pqErr.Code == "23505" &&
			pqErr.Constraint == sessionTableName+"_pkey" {
			if attemptNum >= attemptNumMax {
				return iam.SessionRefKeyZero(), timeZero, errors.Wrap("insert max attempts", err)
			}
			continue
		}
		return iam.SessionRefKeyZero(), timeZero, errors.Wrap("insert", err)
	}

	return iam.NewSessionRefKey(terminalRef, sessionID), sessionStartTime, nil
}
