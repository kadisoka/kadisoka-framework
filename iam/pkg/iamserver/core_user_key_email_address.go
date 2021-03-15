package iamserver

import (
	"database/sql"
	"time"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	goqu "github.com/doug-martin/goqu/v9"
	"github.com/lib/pq"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/eav10n"
)

// Interface conformance assertion.
var _ iam.UserKeyEmailAddressService = &Core{}

const userKeyEmailAddressDBTableName = `user_key_email_address_dt`

func (core *Core) GetUserKeyEmailAddress(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
) (*iam.EmailAddress, error) {
	//TODO: access control
	return core.getUserKeyEmailAddressNoAC(callCtx, userRef)
}

func (core *Core) getUserKeyEmailAddressNoAC(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
) (*iam.EmailAddress, error) {
	var rawInput string

	sqlString, _, _ := goqu.
		From(userKeyEmailAddressDBTableName).
		Select("raw_input").
		Where(
			goqu.C("user_id").Eq(userRef.ID().PrimitiveValue()),
			goqu.C("d_ts").IsNull(),
			goqu.C("verification_time").IsNotNull()).
		ToSQL()

	err := core.db.
		QueryRow(sqlString).
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
	sqlString, _, _ := goqu.
		From(userKeyEmailAddressDBTableName).
		Select("user_id").
		Where(
			goqu.C("domain_part").Eq(emailAddress.DomainPart()),
			goqu.C("local_part").Eq(emailAddress.LocalPart()),
			goqu.C("d_ts").IsNull(),
			goqu.C("verification_time").IsNotNull(),
		).
		ToSQL()
	err = core.db.
		QueryRow(sqlString).
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
func (core *Core) getUserRefByKeyEmailAddressAllowUnverified(
	emailAddress iam.EmailAddress,
) (ownerUserRef iam.UserRefKey, alreadyVerified bool, err error) {
	sqlString, _, _ := goqu.
		From(userKeyEmailAddressDBTableName).
		Select(
			"user_id",
			goqu.Case().
				When(goqu.C("verification_time").IsNull(), false).
				Else(true).
				As("verified"),
		).
		Where(
			goqu.C("domain_part").Eq(emailAddress.DomainPart()),
			goqu.C("local_part").Eq(emailAddress.LocalPart()),
			goqu.C("d_ts").IsNull(),
		).
		Order(goqu.C("c_ts").Desc()).
		Limit(1).
		ToSQL()

	var ownerUserID iam.UserID
	err = core.db.
		QueryRow(sqlString).
		Scan(&ownerUserID, &alreadyVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return iam.UserRefKeyZero(), false, nil
		}
		return iam.UserRefKeyZero(), false, err
	}

	return iam.NewUserRefKey(ownerUserID), alreadyVerified, nil
}

func (core *Core) SetUserKeyEmailAddress(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
	emailAddress iam.EmailAddress,
	verificationMethods []eav10n.VerificationMethod,
) (verificationID int64, verificationCodeExpiry *time.Time, err error) {
	ctxAuth := callCtx.Authorization()
	if !ctxAuth.IsUserContext() {
		return 0, nil, iam.ErrUserContextRequired
	}
	// Don't allow changing other user's for now
	if !ctxAuth.IsUser(userRef) {
		return 0, nil, iam.ErrContextUserNotAllowedToPerformActionOnResource
	}

	existingOwnerUserID, err := core.
		getUserIDByKeyEmailAddress(emailAddress)
	if err != nil {
		return 0, nil, errors.Wrap("getUserIDByKeyEmailAddress", err)
	}
	if existingOwnerUserID.IsValid() {
		if existingOwnerUserID != ctxAuth.UserID() {
			return 0, nil, errors.ArgMsg("emailAddress", "conflict")
		}
		return 0, nil, nil
	}

	alreadyVerified, err := core.setUserKeyEmailAddress(
		callCtx, ctxAuth.UserRef(), emailAddress)
	if err != nil {
		panic(err)
	}
	if alreadyVerified {
		return 0, nil, nil
	}

	//TODO: user-set has higher priority over terminal's
	userLanguages, err := core.getTerminalAcceptLanguages(ctxAuth.TerminalID())
	if err != nil {
	}

	verificationID, verificationCodeExpiry, err = core.eaVerifier.
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
	ctxTime := callCtx.RequestInfo().ReceiveTime
	ctxAuth := callCtx.Authorization()

	xres, err := core.db.Exec(
		`INSERT INTO `+userKeyEmailAddressDBTableName+` (`+
			`user_id, domain_part, local_part, raw_input, `+
			`c_ts, c_uid, c_tid `+
			`) VALUES (`+
			`$1, $2, $3, $4, $5, $6, $7`+
			`) `+
			`ON CONFLICT (user_id, domain_part, local_part) WHERE d_ts IS NULL `+
			`DO NOTHING`,
		userRef.ID().PrimitiveValue(),
		emailAddress.DomainPart(),
		emailAddress.LocalPart(),
		emailAddress.RawInput(),
		ctxTime,
		ctxAuth.UserID().PrimitiveValue(),
		ctxAuth.TerminalID().PrimitiveValue())
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
			`FROM `+userKeyEmailAddressDBTableName+` `+
			`WHERE user_id = $1 AND domain_part = $2 AND local_part = $3`,
		userRef.ID().PrimitiveValue(), emailAddress.DomainPart(), emailAddress.LocalPart()).
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
) (stateChanged bool, err error) {
	ctxAuth := callCtx.Authorization()
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

	ctxTime := callCtx.RequestInfo().ReceiveTime
	stateChanged, err = core.
		ensureUserEmailAddressVerifiedFlag(
			ctxAuth.UserID(), *emailAddress,
			&ctxTime, verificationID)
	if err != nil {
		panic(err)
	}

	return stateChanged, nil
}

func (core *Core) ensureUserEmailAddressVerifiedFlag(
	userID iam.UserID,
	emailAddress iam.EmailAddress,
	verificationTime *time.Time,
	verificationID int64,
) (stateChanged bool, err error) {
	var xres sql.Result

	xres, err = core.db.Exec(
		`UPDATE `+userKeyEmailAddressDBTableName+` SET (`+
			`verification_time, verification_id`+
			`) = ( `+
			`$1, $2`+
			`) WHERE user_id = $3 `+
			`AND domain_part = $4 AND local_part = $5 `+
			`AND d_ts IS NULL AND verification_time IS NULL`,
		verificationTime,
		verificationID,
		userID,
		emailAddress.DomainPart(),
		emailAddress.LocalPart())
	if err != nil {
		pqErr, _ := err.(*pq.Error)
		if pqErr != nil &&
			pqErr.Code == "23505" &&
			pqErr.Constraint == userKeyEmailAddressDBTableName+`_local_part_domain_part_uidx` {
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
