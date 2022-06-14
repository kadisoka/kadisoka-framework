package iamserver

import (
	"github.com/alloyzeus/go-azfl/errors"
	"github.com/jmoiron/sqlx"

	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
)

// Interface conformance assertion.
var _ iam.TerminalFCMRegistrationTokenService = &Core{}

const terminalFCMRegistrationTokenDBTableName = "terminal_fcm_registration_token_dt"

func (core *Core) DisposeTerminalFCMRegistrationToken(
	inputCtx iam.CallInputContext,
	terminalID iam.TerminalID,
	token string,
) error {
	ctxAuth := inputCtx.Authorization()
	_, err := core.db.Exec(
		`UPDATE `+terminalFCMRegistrationTokenDBTableName+` `+
			"SET _md_ts = $1, _md_uid = $2, _md_tid = $3 "+
			"WHERE terminal_id = $4 AND token = $5 AND _md_ts IS NULL",
		inputCtx.CallInputMetadata().ReceiveTime,
		ctxAuth.UserIDNum().PrimitiveValue(),
		ctxAuth.TerminalIDNum().PrimitiveValue(),
		terminalID.IDNum().PrimitiveValue(),
		token)
	return err
}

func (core *Core) ListTerminalFCMRegistrationTokensByUser(
	ownerUserID iam.UserID,
) (tokens map[iam.TerminalID]string, err error) {
	//TODO: use cache service

	termRows, err := core.db.Query(
		`SELECT tid.id_num, tid.application_id, tft.token `+
			`FROM `+terminalDBTableName+` tid `+
			`JOIN `+terminalFCMRegistrationTokenDBTableName+` tft `+
			"ON tft.terminal_id=tid.id_num AND tft._md_ts IS NULL "+
			"WHERE tid.user_id=$1 AND tid._md_ts IS NULL AND tid.verification_ts IS NOT NULL",
		ownerUserID.IDNum().PrimitiveValue())
	if err != nil {
		return nil, err
	}
	defer termRows.Close()

	result := map[iam.TerminalID]string{}
	for termRows.Next() {
		var terminalIDNum iam.TerminalIDNum
		var applicationIDNum iam.ApplicationIDNum
		var token string
		if err = termRows.Scan(&terminalIDNum, &applicationIDNum, &token); err != nil {
			return nil, err
		}
		result[iam.NewTerminalID(
			iam.NewApplicationID(applicationIDNum),
			ownerUserID,
			terminalIDNum,
		)] = token
	}
	if err = termRows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (core *Core) SetTerminalFCMRegistrationToken(
	inputCtx iam.CallInputContext,
	terminalID iam.TerminalID, userID iam.UserID,
	token string,
) error {
	if inputCtx == nil {
		return errors.ArgMsg("inputCtx", "missing")
	}

	ctxTermID := inputCtx.Authorization().TerminalID()
	ctxAppIDNum := ctxTermID.Application().IDNum()
	if !ctxAppIDNum.IsFirstParty() || !ctxAppIDNum.IsService() {
		return errors.ArgMsg("inputCtx", "unauthorized application type")
	}
	if !ctxTermID.User().EqualsUserID(userID) {
		return errors.ArgMsg("inputCtx", "terminal user mismatch")
	}

	ctxAuth := inputCtx.Authorization()

	return doTx(core.db, func(tx *sqlx.Tx) error {
		_, err := tx.Exec(
			`UPDATE `+terminalFCMRegistrationTokenDBTableName+` `+
				"SET _md_ts = $1, _md_uid = $2, _md_tid = $3 "+
				"WHERE terminal_id = $4 AND _md_ts IS NULL",
			inputCtx.CallInputMetadata().ReceiveTime,
			ctxAuth.UserIDNum().PrimitiveValue(),
			ctxAuth.TerminalIDNum().PrimitiveValue(),
			terminalID.IDNum().PrimitiveValue())
		if err != nil {
			return err
		}
		if token == "" {
			return nil
		}
		_, err = tx.Exec(
			`INSERT INTO `+terminalFCMRegistrationTokenDBTableName+` `+
				"(terminal_id, user_id, _mc_uid, _mc_tid, token) "+
				"VALUES ($1, $2, $3, $4, $5)",
			terminalID.IDNum().PrimitiveValue(),
			userID.IDNum().PrimitiveValue(),
			ctxAuth.UserIDNum().PrimitiveValue(),
			ctxAuth.TerminalIDNum().PrimitiveValue(),
			token)
		return err
	})
}
