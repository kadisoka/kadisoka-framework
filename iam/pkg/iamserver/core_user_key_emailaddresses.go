package iamserver

import (
	"database/sql"
	"time"

	"github.com/alloyzeus/go-azcore/azcore/errors"
	"github.com/lib/pq"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/eav10n"
)

// Key email address is an email address which can be used to sign in.

const userKeyEmailAddressTableName = `user_key_email_address_dt`

func (core *Core) GetUserKeyEmailAddress(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
) (*iam.EmailAddress, error) {
	var rawInput string
	err := core.db.
		QueryRow(
			`SELECT raw_input `+
				`FROM `+userKeyEmailAddressTableName+` `+
				`WHERE user_id=$1 `+
				`AND d_ts IS NULL AND verification_time IS NOT NULL`,
			userRef.ID().PrimitiveValue()).
		Scan(&rawInput)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	emailAddress, err := iam.EmailAddressFromString(rawInput)
	if err != nil {
		panic(err)
	}
	return &emailAddress, nil
}

// The ID of the user which provided email address is their verified primary.
func (core *Core) getUserIDByKeyEmailAddress(
	emailAddress iam.EmailAddress,
) (ownerUserID iam.UserID, err error) {
	queryStr :=
		`SELECT user_id ` +
			`FROM ` + userKeyEmailAddressTableName + ` ` +
			`WHERE local_part = $1 AND domain_part = $2 ` +
			`AND d_ts IS NULL ` +
			`AND verification_time IS NOT NULL`
	err = core.db.
		QueryRow(queryStr,
			emailAddress.LocalPart(),
			emailAddress.DomainPart()).
		Scan(&ownerUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return iam.UserIDZero, nil
		}
		return iam.UserIDZero, err
	}

	return
}

// The ID of the user which provided email address is their primary,
// verified or not.
func (core *Core) getUserIDByKeyEmailAddressAllowUnverified(
	emailAddress iam.EmailAddress,
) (ownerUserRef iam.UserRefKey, verified bool, err error) {
	queryStr :=
		`SELECT user_id, CASE WHEN verification_time IS NULL THEN false ELSE true END AS verified ` +
			`FROM ` + userKeyEmailAddressTableName + ` ` +
			`WHERE local_part = $1 AND domain_part = $2 ` +
			`AND d_ts IS NULL ` +
			`ORDER BY c_ts DESC LIMIT 1`
	var ownerUserID iam.UserID
	err = core.db.
		QueryRow(queryStr,
			emailAddress.LocalPart(),
			emailAddress.DomainPart()).
		Scan(&ownerUserRef, &verified)
	if err != nil {
		if err == sql.ErrNoRows {
			return iam.UserRefKeyZero(), false, nil
		}
		return iam.UserRefKeyZero(), false, err
	}

	return iam.NewUserRefKey(ownerUserID), verified, nil
}

func (core *Core) SetUserKeyEmailAddress(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
	emailAddress iam.EmailAddress,
	verificationMethods []eav10n.VerificationMethod,
) (verificationID int64, codeExpiry *time.Time, err error) {
	authCtx := callCtx.Authorization()
	if !authCtx.IsUserContext() {
		return 0, nil, iam.ErrUserContextRequired
	}
	// Don't allow changing other user's for now
	if !userRef.EqualsUserRefKey(authCtx.UserRef()) {
		return 0, nil, iam.ErrContextUserNotAllowedToPerformActionOnResource
	}

	existingOwnerUserID, err := core.
		getUserIDByKeyEmailAddress(emailAddress)
	if err != nil {
		return 0, nil, errors.Wrap("getUserIDByKeyEmailAddress", err)
	}
	if existingOwnerUserID.IsValid() {
		if existingOwnerUserID != authCtx.UserID() {
			return 0, nil, errors.ArgMsg("emailAddress", "conflict")
		}
		return 0, nil, nil
	}

	alreadyVerified, err := core.setUserKeyEmailAddress(
		callCtx, authCtx.UserRef(), emailAddress)
	if err != nil {
		panic(err)
	}
	if alreadyVerified {
		return 0, nil, nil
	}

	//TODO: user-set has higher priority over terminal's
	userLanguages, err := core.getTerminalAcceptLanguages(authCtx.TerminalID())

	verificationID, codeExpiry, err = core.eaVerifier.
		StartVerification(callCtx, emailAddress,
			0, userLanguages, verificationMethods)
	if err != nil {
		switch err.(type) {
		case eav10n.InvalidEmailAddressError:
			return 0, nil, errors.Arg("emailAddress", err)
		}
		return 0, nil, errors.Wrap("eaVerifier.StartVerification", err)
	}

	return
}

func (core *Core) setUserKeyEmailAddress(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
	emailAddress iam.EmailAddress,
) (alreadyVerified bool, err error) {
	xres, err := core.db.Exec(
		`INSERT INTO `+userKeyEmailAddressTableName+` (`+
			`user_id, local_part, domain_part, raw_input, `+
			`c_ts, c_uid, c_tid `+
			`) VALUES (`+
			`$1, $2, $3, $4, $5, $6, $7`+
			`) `+
			`ON CONFLICT (user_id, local_part, domain_part) WHERE d_ts IS NULL `+
			`DO NOTHING`,
		userRef.ID().PrimitiveValue(),
		emailAddress.LocalPart(),
		emailAddress.DomainPart(),
		emailAddress.RawInput(),
		callCtx.RequestReceiveTime(),
		callCtx.Authorization().UserID().PrimitiveValue(),
		callCtx.Authorization().TerminalID().PrimitiveValue())
	if err != nil {
		return false, err
	}

	n, err := xres.RowsAffected()
	if err != nil {
		return false, err
	}
	if n == 1 {
		return false, nil
	}

	err = core.db.QueryRow(
		`SELECT CASE WHEN verification_time IS NULL THEN false ELSE true END AS verified `+
			`FROM `+userKeyEmailAddressTableName+` `+
			`WHERE user_id = $1 AND local_part = $2 AND domain_part = $3`,
		userRef, emailAddress.LocalPart(), emailAddress.DomainPart()).
		Scan(&alreadyVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return
}

func (core *Core) ConfirmUserEmailAddressVerification(
	callCtx iam.CallContext,
	verificationID int64,
	code string,
) (updated bool, err error) {
	authCtx := callCtx.Authorization()
	err = core.eaVerifier.ConfirmVerification(
		callCtx, verificationID, code)
	if err != nil {
		switch err {
		case eav10n.ErrVerificationCodeMismatch:
			return false, errors.ArgMsg("code", "mismatch")
		case eav10n.ErrVerificationCodeExpired:
			return false, errors.ArgMsg("code", "expired")
		}
		return false, errors.Wrap("eaVerifier.ConfirmVerification", err)
	}

	emailAddress, err := core.eaVerifier.
		GetEmailAddressByVerificationID(verificationID)
	// An unexpected condition which could cause bad state
	if err != nil {
		panic(err)
	}

	ctxTime := callCtx.RequestReceiveTime()
	updated, err = core.
		ensureUserEmailAddressVerifiedFlag(
			authCtx.UserID(), *emailAddress,
			&ctxTime, verificationID)
	if err != nil {
		panic(err)
	}

	return updated, nil
}

func (core *Core) ensureUserEmailAddressVerifiedFlag(
	userID iam.UserID,
	emailAddress iam.EmailAddress,
	verificationTime *time.Time,
	verificationID int64,
) (bool, error) {
	var err error
	var xres sql.Result

	xres, err = core.db.Exec(
		`UPDATE `+userKeyEmailAddressTableName+` SET (`+
			`verification_time, verification_id`+
			`) = ( `+
			`$1, $2`+
			`) WHERE user_id = $3 `+
			`AND local_part = $4 AND domain_part = $5 `+
			`AND d_ts IS NULL AND verification_time IS NULL`,
		verificationTime,
		verificationID,
		userID,
		emailAddress.LocalPart(),
		emailAddress.DomainPart())
	if err != nil {
		pqErr, _ := err.(*pq.Error)
		if pqErr != nil &&
			pqErr.Code == "23505" &&
			pqErr.Constraint == userKeyEmailAddressTableName+`_local_part_domain_part_uidx` {
			return false, errors.ArgMsg("emailAddress", "conflict")
		}
		return false, err
	}

	var n int64
	n, err = xres.RowsAffected()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}
