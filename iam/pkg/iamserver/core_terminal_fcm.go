package iamserver

import (
	"github.com/alloyzeus/go-azfl/azfl/errors"
	"github.com/jmoiron/sqlx"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

// Interface conformance assertion.
var _ iam.TerminalFCMRegistrationTokenService = &Core{}

const terminalFCMRegistrationTokenTableName = "terminal_fcm_registration_token_dt"

func (core *Core) DisposeTerminalFCMRegistrationToken(
	callCtx iam.CallContext,
	terminalID iam.TerminalID,
	token string,
) error {
	ctxAuth := callCtx.Authorization()
	_, err := core.db.Exec(
		`UPDATE `+terminalFCMRegistrationTokenTableName+` `+
			"SET d_ts = $1, d_uid = $2, d_tid = $3 "+
			"WHERE terminal_id = $4 AND token = $5 AND d_ts IS NULL",
		callCtx.RequestInfo().ReceiveTime, ctxAuth.UserID().PrimitiveValue(), ctxAuth.TerminalID().PrimitiveValue(),
		terminalID.PrimitiveValue(), token)
	return err
}

func (core *Core) ListTerminalFCMRegistrationTokensByUser(
	ownerUserRef iam.UserRefKey,
) (tokens map[iam.TerminalID]string, err error) {
	//TODO: use cache service

	userTermRows, err := core.db.Query(
		`SELECT tid.id, tft.token `+
			`FROM `+terminalTableName+` tid `+
			`JOIN `+terminalFCMRegistrationTokenTableName+` tft `+
			"ON tft.terminal_id=tid.id AND tft.d_ts IS NULL "+
			"WHERE tid.user_id=$1 AND tid.verification_time IS NOT NULL",
		ownerUserRef.ID().PrimitiveValue())
	if err != nil {
		return nil, err
	}
	defer userTermRows.Close()

	result := map[iam.TerminalID]string{}
	for userTermRows.Next() {
		var terminalID iam.TerminalID
		var token string
		if err = userTermRows.Scan(&terminalID, &token); err != nil {
			return nil, err
		}
		result[terminalID] = token
	}
	if err = userTermRows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (core *Core) SetTerminalFCMRegistrationToken(
	callCtx iam.CallContext,
	terminalRef iam.TerminalRefKey, userRef iam.UserRefKey,
	token string,
) error {
	if callCtx == nil {
		return errors.ArgMsg("callCtx", "missing")
	}

	ctxTermRef := callCtx.Authorization().TerminalRef()
	ctxAppID := ctxTermRef.Application().ID()
	if !ctxAppID.IsFirstParty() || !ctxAppID.IsService() {
		return errors.ArgMsg("callCtx", "unauthorized application type")
	}
	if !ctxTermRef.User().EqualsUserRefKey(userRef) {
		return errors.ArgMsg("callCtx", "terminal user mismatch")
	}

	ctxAuth := callCtx.Authorization()

	return doTx(core.db, func(tx *sqlx.Tx) error {
		_, err := tx.Exec(
			`UPDATE `+terminalFCMRegistrationTokenTableName+` `+
				"SET d_ts = $1, d_uid = $2, d_tid = $3 "+
				"WHERE terminal_id = $4 AND d_ts IS NULL",
			callCtx.RequestInfo().ReceiveTime,
			ctxAuth.UserID().PrimitiveValue(),
			ctxAuth.TerminalID().PrimitiveValue(),
			terminalRef.ID().PrimitiveValue())
		if err != nil {
			return err
		}
		if token == "" {
			return nil
		}
		_, err = tx.Exec(
			`INSERT INTO `+terminalFCMRegistrationTokenTableName+` `+
				"(terminal_id, user_id, c_uid, c_tid, token) "+
				"VALUES ($1, $2, $3, $4, $5)",
			terminalRef.ID().PrimitiveValue(), userRef.ID().PrimitiveValue(),
			ctxAuth.UserID().PrimitiveValue(), ctxAuth.TerminalID().PrimitiveValue(),
			token)
		return err
	})
}
