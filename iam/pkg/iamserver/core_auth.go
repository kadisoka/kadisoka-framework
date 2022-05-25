package iamserver

import (
	"time"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/lib/pq"
	"github.com/square/go-jose/v3/jwt"

	apperrs "github.com/kadisoka/kadisoka-framework/foundation/pkg/app/errors"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

func (core *Core) GenerateAccessTokenJWT(
	callCtx iam.OpInputContext,
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

	sessionRef, issueTime, expiry, err := core.
		issueSession(callCtx, terminalRef, userRef)
	if err != nil {
		return "", err
	}

	tokenClaims := &iam.AccessTokenClaims{
		Claims: jwt.Claims{
			ID:       sessionRef.AZIDText(),
			IssuedAt: jwt.NewNumericDate(issueTime),
			Issuer:   core.RealmName(),
			Expiry:   jwt.NewNumericDate(expiry),
			Subject:  userRef.AZIDText(),
		},
		AuthorizedParty: terminalRef.Application().AZIDText(),
		TerminalID:      terminalRef.AZIDText(),
	}

	tokenString, err = jwt.Signed(signer).Claims(tokenClaims).CompactSerialize()
	if err != nil {
		return "", errors.Wrap("signing", err)
	}
	return
}

func (core *Core) GenerateRefreshTokenJWT(
	callCtx iam.OpInputContext,
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
		TerminalID:     terminalRef.AZIDText(),
		TerminalSecret: terminalSecret,
	}

	tokenString, err = jwt.Signed(signer).Claims(tokenClaims).CompactSerialize()
	if err != nil {
		return "", errors.Wrap("signing", err)
	}
	return
}

func (core *Core) issueSession(
	callCtx iam.OpInputContext,
	terminalRef iam.TerminalRefKey,
	userRef iam.UserRefKey,
) (
	sessionRef iam.SessionRefKey,
	issueTime time.Time,
	expiry time.Time,
	err error,
) {
	ctxAuth := callCtx.Authorization()

	const attemptNumMax = 5

	timeZero := time.Time{}
	sessionStartTime := timeZero
	sessionExpiry := timeZero
	var sessionIDNum iam.SessionIDNum

	for attemptNum := 0; ; attemptNum++ {
		sessionStartTime = time.Now().UTC()
		sessionExpiry = sessionStartTime.Add(iam.AccessTokenTTLDefault)
		sessionIDNum, err = GenerateSessionIDNum(0)
		if err != nil {
			return iam.SessionRefKeyZero(), timeZero, timeZero, err
		}
		sqlString, _, _ := goqu.
			Insert(sessionDBTableName).
			Rows(
				goqu.Record{
					"terminal_id": terminalRef.IDNum().PrimitiveValue(),
					"id_num":      sessionIDNum.PrimitiveValue(),
					"expiry":      sessionExpiry,
					"_mc_ts":      sessionStartTime,
					"_mc_tid":     ctxAuth.TerminalIDNumPtr(),
					"_mc_uid":     ctxAuth.UserIDNumPtr(),
				},
			).
			ToSQL()
		_, err = core.db.
			Exec(sqlString)
		if err == nil {
			break
		}

		pqErr, _ := err.(*pq.Error)
		if pqErr != nil &&
			pqErr.Code == "23505" &&
			pqErr.Constraint == sessionDBTableName+"_pkey" {
			if attemptNum >= attemptNumMax {
				return iam.SessionRefKeyZero(), timeZero, timeZero,
					errors.Wrap("insert max attempts", err)
			}
			continue
		}

		return iam.SessionRefKeyZero(), timeZero, timeZero,
			errors.Wrap("insert", err)
	}

	return iam.NewSessionRefKey(terminalRef, sessionIDNum),
		sessionStartTime, sessionExpiry, nil
}
