package iamserver

import (
	"github.com/alloyzeus/go-azfl/azfl/errors"
	"github.com/jmoiron/sqlx"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

// Interface conformance assertion.
var _ iam.TerminalFCMRegistrationTokenService = &Core{}

const terminalFCMRegistrationTokenDBTableName = "terminal_fcm_registration_token_dt"

func (core *Core) DisposeTerminalFCMRegistrationToken(
	callCtx iam.OpInputContext,
	terminalRef iam.TerminalRefKey,
	token string,
) error {
	ctxAuth := callCtx.Authorization()
	_, err := core.db.Exec(
		`UPDATE `+terminalFCMRegistrationTokenDBTableName+` `+
			"SET d_ts = $1, d_uid = $2, d_tid = $3 "+
			"WHERE terminal_id = $4 AND token = $5 AND d_ts IS NULL",
		callCtx.OpInputMetadata().ReceiveTime,
		ctxAuth.UserIDNum().PrimitiveValue(),
		ctxAuth.TerminalIDNum().PrimitiveValue(),
		terminalRef.IDNum().PrimitiveValue(),
		token)
	return err
}

func (core *Core) ListTerminalFCMRegistrationTokensByUser(
	ownerUserRef iam.UserRefKey,
) (tokens map[iam.TerminalRefKey]string, err error) {
	//TODO: use cache service

	termRows, err := core.db.Query(
		`SELECT tid.id, tid.application_id, tft.token `+
			`FROM `+terminalDBTableName+` tid `+
			`JOIN `+terminalFCMRegistrationTokenDBTableName+` tft `+
			"ON tft.terminal_id=tid.id AND tft.d_ts IS NULL "+
			"WHERE tid.user_id=$1 AND tid.d_ts IS NULL AND tid.verification_ts IS NOT NULL",
		ownerUserRef.IDNum().PrimitiveValue())
	if err != nil {
		return nil, err
	}
	defer termRows.Close()

	result := map[iam.TerminalRefKey]string{}
	for termRows.Next() {
		var terminalIDNum iam.TerminalIDNum
		var applicationIDNum iam.ApplicationIDNum
		var token string
		if err = termRows.Scan(&terminalIDNum, &applicationIDNum, &token); err != nil {
			return nil, err
		}
		result[iam.NewTerminalRefKey(
			iam.NewApplicationRefKey(applicationIDNum),
			ownerUserRef,
			terminalIDNum,
		)] = token
	}
	if err = termRows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (core *Core) SetTerminalFCMRegistrationToken(
	callCtx iam.OpInputContext,
	terminalRef iam.TerminalRefKey, userRef iam.UserRefKey,
	token string,
) error {
	if callCtx == nil {
		return errors.ArgMsg("callCtx", "missing")
	}

	ctxTermRef := callCtx.Authorization().TerminalRef()
	ctxAppIDNum := ctxTermRef.Application().IDNum()
	if !ctxAppIDNum.IsFirstParty() || !ctxAppIDNum.IsService() {
		return errors.ArgMsg("callCtx", "unauthorized application type")
	}
	if !ctxTermRef.User().EqualsUserRefKey(userRef) {
		return errors.ArgMsg("callCtx", "terminal user mismatch")
	}

	ctxAuth := callCtx.Authorization()

	return doTx(core.db, func(tx *sqlx.Tx) error {
		_, err := tx.Exec(
			`UPDATE `+terminalFCMRegistrationTokenDBTableName+` `+
				"SET d_ts = $1, d_uid = $2, d_tid = $3 "+
				"WHERE terminal_id = $4 AND d_ts IS NULL",
			callCtx.OpInputMetadata().ReceiveTime,
			ctxAuth.UserIDNum().PrimitiveValue(),
			ctxAuth.TerminalIDNum().PrimitiveValue(),
			terminalRef.IDNum().PrimitiveValue())
		if err != nil {
			return err
		}
		if token == "" {
			return nil
		}
		_, err = tx.Exec(
			`INSERT INTO `+terminalFCMRegistrationTokenDBTableName+` `+
				"(terminal_id, user_id, c_uid, c_tid, token) "+
				"VALUES ($1, $2, $3, $4, $5)",
			terminalRef.IDNum().PrimitiveValue(),
			userRef.IDNum().PrimitiveValue(),
			ctxAuth.UserIDNum().PrimitiveValue(),
			ctxAuth.TerminalIDNum().PrimitiveValue(),
			token)
		return err
	})
}
