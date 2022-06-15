package iamserver

import (
	"time"

	errors "github.com/alloyzeus/go-azfl/errors"
	"github.com/square/go-jose/v3/jwt"

	apperrs "github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/app/errors"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
)

func (core *Core) GenerateTokenSetJWT(
	inputCtx iam.CallInputContext,
	terminalID iam.TerminalID,
	userID iam.UserID,
	terminalSecret string,
) (accessToken string, refreshToken string, err error) {
	if inputCtx == nil {
		return "", "", errors.ArgMsg("inputCtx", "missing")
	}

	if terminalID.IsNotStaticallyValid() {
		return "", "", errors.ArgMsg("terminalID", "invalid")
	}
	if (terminalID.Application().AZIDNum().IsUserAgent() && userID.IsNotStaticallyValid()) ||
		(!terminalID.Application().AZIDNum().IsUserAgent() && userID.IsStaticallyValid()) {
		return "", "", errors.ArgMsg("userID", "invalid combination with terminalID")
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

	sessionID, issueTime, expiry, err := core.
		issueSession(inputCtx, terminalID, userID)
	if err != nil {
		return "", "", err
	}

	accessTokenClaims := &iam.AccessTokenClaims{
		Claims: jwt.Claims{
			ID:       sessionID.AZIDText(),
			IssuedAt: jwt.NewNumericDate(issueTime),
			Issuer:   core.RealmName(),
			Expiry:   jwt.NewNumericDate(expiry),
			Subject:  userID.AZIDText(),
		},
		AuthorizedParty: terminalID.Application().AZIDText(),
		TerminalID:      terminalID.AZIDText(),
	}

	accessToken, err = jwt.Signed(signer).Claims(accessTokenClaims).
		CompactSerialize()
	if err != nil {
		return "", "", errors.Wrap("access token signing", err)
	}

	tokenClaims := &iam.RefreshTokenClaims{
		NotBefore:      issueTime.Unix(),
		ExpiresAt:      issueTime.Add(iam.RefreshTokenTTLDefault).Unix(),
		TerminalID:     terminalID.AZIDText(),
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
	inputCtx iam.CallInputContext,
	terminalID iam.TerminalID,
	userID iam.UserID,
) (tokenString string, err error) {
	if inputCtx == nil {
		return "", errors.ArgMsg("inputCtx", "missing")
	}

	if terminalID.IsNotStaticallyValid() {
		return "", errors.ArgMsg("terminalID", "invalid")
	}
	if (terminalID.Application().AZIDNum().IsUserAgent() && userID.IsNotStaticallyValid()) ||
		(!terminalID.Application().AZIDNum().IsUserAgent() && userID.IsStaticallyValid()) {
		return "", errors.ArgMsg("userID", "invalid combination with terminalID")
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

	sessionID, issueTime, expiry, err := core.
		issueSession(inputCtx, terminalID, userID)
	if err != nil {
		return "", err
	}

	tokenClaims := &iam.AccessTokenClaims{
		Claims: jwt.Claims{
			ID:       sessionID.AZIDText(),
			IssuedAt: jwt.NewNumericDate(issueTime),
			Issuer:   core.RealmName(),
			Expiry:   jwt.NewNumericDate(expiry),
			Subject:  userID.AZIDText(),
		},
		AuthorizedParty: terminalID.Application().AZIDText(),
		TerminalID:      terminalID.AZIDText(),
	}

	tokenString, err = jwt.Signed(signer).Claims(tokenClaims).
		CompactSerialize()
	if err != nil {
		return "", errors.Wrap("signing", err)
	}
	return
}

func (core *Core) GenerateRefreshTokenJWT(
	inputCtx iam.CallInputContext,
	terminalID iam.TerminalID,
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
		TerminalID:     terminalID.AZIDText(),
		TerminalSecret: terminalSecret,
	}

	tokenString, err = jwt.Signed(signer).Claims(tokenClaims).
		CompactSerialize()
	if err != nil {
		return "", errors.Wrap("signing", err)
	}

	return
}
