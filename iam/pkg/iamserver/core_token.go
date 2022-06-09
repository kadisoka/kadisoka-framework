package iamserver

import (
	"time"

	errors "github.com/alloyzeus/go-azfl/errors"
	"github.com/square/go-jose/v3/jwt"

	apperrs "github.com/kadisoka/kadisoka-framework/foundation/pkg/app/errors"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

func (core *Core) GenerateTokenSetJWT(
	callCtx iam.CallInputContext,
	terminalRef iam.TerminalRefKey,
	userRef iam.UserRefKey,
	terminalSecret string,
) (accessToken string, refreshToken string, err error) {
	if callCtx == nil {
		return "", "", errors.ArgMsg("callCtx", "missing")
	}

	jwtKeyChain := core.JWTKeyChain()
	if jwtKeyChain == nil {
		return "", "", apperrs.NewConfigurationMsg("JWT key chain is not configured")
	}
	signer, err := jwtKeyChain.GetSigner()
	if err != nil {
		return "", "", errors.Wrap("signer", err)
	}
	if signer == nil {
		return "", "", apperrs.NewConfigurationMsg("JWT key chain does not have any signing key")
	}

	sessionRef, issueTime, expiry, err := core.
		issueSession(callCtx, terminalRef, userRef)
	if err != nil {
		return "", "", err
	}

	accessTokenClaims := &iam.AccessTokenClaims{
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

	accessToken, err = jwt.Signed(signer).Claims(accessTokenClaims).
		CompactSerialize()
	if err != nil {
		return "", "", errors.Wrap("access token signing", err)
	}

	tokenClaims := &iam.RefreshTokenClaims{
		NotBefore:      issueTime.Unix(),
		ExpiresAt:      issueTime.Add(iam.RefreshTokenTTLDefault).Unix(),
		TerminalID:     terminalRef.AZIDText(),
		TerminalSecret: terminalSecret,
	}

	refreshToken, err = jwt.Signed(signer).Claims(tokenClaims).
		CompactSerialize()
	if err != nil {
		return "", "", errors.Wrap("refresh token signing", err)
	}

	return
}

func (core *Core) GenerateAccessTokenJWT(
	callCtx iam.CallInputContext,
	terminalRef iam.TerminalRefKey,
	userRef iam.UserRefKey,
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

	tokenString, err = jwt.Signed(signer).Claims(tokenClaims).
		CompactSerialize()
	if err != nil {
		return "", errors.Wrap("signing", err)
	}
	return
}

func (core *Core) GenerateRefreshTokenJWT(
	callCtx iam.CallInputContext,
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

	tokenString, err = jwt.Signed(signer).Claims(tokenClaims).
		CompactSerialize()
	if err != nil {
		return "", errors.Wrap("signing", err)
	}

	return
}
