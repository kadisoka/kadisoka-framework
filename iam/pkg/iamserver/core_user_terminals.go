package iamserver

import (
	"github.com/alloyzeus/go-azcore/azcore/errors"
	"github.com/jmoiron/sqlx"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

func (core *Core) DeleteUserTerminalFCMRegistrationToken(
	callCtx iam.CallContext,
	userID iam.UserID,
	terminalID iam.TerminalID,
	token string,
) error {
	authCtx := callCtx.Authorization()
	_, err := core.db.Exec(
		"UPDATE user_terminal_fcm_registration_tokens "+
			"SET d_ts = $1, d_uid = $2, d_tid = $3 "+
			"WHERE user_id = $4 AND terminal_id = $5 AND token = $6 AND d_ts IS NULL",
		callCtx.RequestReceiveTime(), authCtx.UserID().PrimitiveValue(), authCtx.TerminalID().PrimitiveValue(),
		userID, terminalID, token)
	return err
}

func (core *Core) ListUserTerminalIDFirebaseInstanceTokens(
	ownerUserID iam.UserID,
) ([]iam.TerminalIDFirebaseInstanceToken, error) {
	//TODO: use cache service

	userTermRows, err := core.db.Query(
		`SELECT tid.id, tft.token `+
			`FROM `+terminalTableName+` tid `+
			"JOIN user_terminal_fcm_registration_tokens tft "+
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

func (core *Core) SetUserTerminalFCMRegistrationToken(
	callCtx iam.CallContext,
	userRef iam.UserRefKey, terminalRef iam.TerminalRefKey, token string,
) error {
	if callCtx == nil {
		return errors.ArgMsg("callCtx", "missing")
	}
	authCtx := callCtx.Authorization()

	return doTx(core.db, func(tx *sqlx.Tx) error {
		_, err := tx.Exec(
			"UPDATE user_terminal_fcm_registration_tokens "+
				"SET d_ts = $1, d_uid = $2, d_tid = $3 "+
				"WHERE user_id = $4 AND terminal_id = $5 AND d_ts IS NULL",
			callCtx.RequestReceiveTime(), authCtx.UserID().PrimitiveValue(), authCtx.TerminalID().PrimitiveValue(),
			userRef.ID().PrimitiveValue(), terminalRef.ID().PrimitiveValue())
		if err != nil {
			return err
		}
		if token == "" {
			return nil
		}
		_, err = tx.Exec(
			"INSERT INTO user_terminal_fcm_registration_tokens "+
				"(user_id, terminal_id, c_uid, c_tid, token) "+
				"VALUES ($1, $2, $3, $4, $5)",
			userRef, terminalRef,
			authCtx.UserID().PrimitiveValue(), authCtx.TerminalID().PrimitiveValue(),
			token)
		return err
	})
}
