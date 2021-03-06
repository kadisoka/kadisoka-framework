package iamserver

import (
	"github.com/alloyzeus/go-azfl/azfl/errors"
	"github.com/jmoiron/sqlx"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

const terminalFCMRegistrationTokenTableName = "terminal_fcm_registration_token_dt"

func (core *Core) DisposeTerminalFCMRegistrationToken(
	callCtx iam.CallContext,
	terminalID iam.TerminalID,
	token string,
) error {
	authCtx := callCtx.Authorization()
	_, err := core.db.Exec(
		`UPDATE `+terminalFCMRegistrationTokenTableName+` `+
			"SET d_ts = $1, d_uid = $2, d_tid = $3 "+
			"WHERE terminal_id = $4 AND token = $5 AND d_ts IS NULL",
		callCtx.RequestReceiveTime(), authCtx.UserID().PrimitiveValue(), authCtx.TerminalID().PrimitiveValue(),
		terminalID.PrimitiveValue(), token)
	return err
}

func (core *Core) ListTerminalIDFCMRegistrationTokensByUser(
	ownerUserID iam.UserID,
) ([]iam.TerminalIDFirebaseInstanceToken, error) {
	//TODO: use cache service

	userTermRows, err := core.db.Query(
		`SELECT tid.id, tft.token `+
			`FROM `+terminalTableName+` tid `+
			`JOIN `+terminalFCMRegistrationTokenTableName+` tft `+
			"ON tft.terminal_id=tid.id AND tft.d_ts IS NULL "+
			"WHERE tid.user_id=$1 AND tid.verification_time IS NOT NULL",
		ownerUserID)
	if err != nil {
		return nil, err
	}
	defer userTermRows.Close()

	var result []iam.TerminalIDFirebaseInstanceToken
	for userTermRows.Next() {
		var item iam.TerminalIDFirebaseInstanceToken
		if err = userTermRows.Scan(&item.TerminalID, &item.Token); err != nil {
			return nil, err
		}
		result = append(result, item)
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
	authCtx := callCtx.Authorization()

	return doTx(core.db, func(tx *sqlx.Tx) error {
		_, err := tx.Exec(
			`UPDATE `+terminalFCMRegistrationTokenTableName+` `+
				"SET d_ts = $1, d_uid = $2, d_tid = $3 "+
				"WHERE terminal_id = $4 AND  d_ts IS NULL",
			callCtx.RequestReceiveTime(), authCtx.UserID().PrimitiveValue(), authCtx.TerminalID().PrimitiveValue(),
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
			authCtx.UserID().PrimitiveValue(), authCtx.TerminalID().PrimitiveValue(),
			token)
		return err
	})
}
