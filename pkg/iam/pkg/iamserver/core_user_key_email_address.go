package iamserver

import (
	"database/sql"
	"time"

	"github.com/alloyzeus/go-azfl/errors"
	goqu "github.com/doug-martin/goqu/v9"
	"github.com/lib/pq"

	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iamserver/eav10n"
	"github.com/kadisoka/kadisoka-framework/pkg/volib/pkg/email"
)

// Interface conformance assertion.
var _ iam.UserKeyEmailAddressService = &Core{}

const userKeyEmailAddressDBTableName = `user_key_email_address_dt`

func (core *Core) GetUserKeyEmailAddress(
	inputCtx iam.CallInputContext,
	userID iam.UserID,
) (*email.Address, error) {
	//TODO: access control
	return core.getUserKeyEmailAddressInsecure(inputCtx, userID)
}

func (core *Core) getUserKeyEmailAddressInsecure(
	inputCtx iam.CallInputContext,
	userID iam.UserID,
) (*email.Address, error) {
	var rawInput string

	sqlString, _, _ := goqu.
		From(userKeyEmailAddressDBTableName).
		Select("raw_input").
		Where(
			goqu.C("user_id").Eq(userID.IDNum().PrimitiveValue()),
			goqu.C("md_d_ts").IsNull(),
			goqu.C("verification_ts").IsNotNull()).
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

	emailAddress, err := email.AddressFromString(rawInput)
	if err != nil {
		panic(err)
	}
	return &emailAddress, nil
}

// The ID of the user which provided email address is their verified primary.
func (core *Core) getUserIDNumByKeyEmailAddressInsecure(
	emailAddress email.Address,
) (ownerUserIDNum iam.UserIDNum, err error) {
	sqlString, _, _ := goqu.
		From(userKeyEmailAddressDBTableName).
		Select("user_id").
		Where(
			goqu.C("domain_part").Eq(emailAddress.DomainPart()),
			goqu.C("local_part").Eq(emailAddress.LocalPart()),
			goqu.C("md_d_ts").IsNull(),
			goqu.C("verification_ts").IsNotNull(),
		).
		ToSQL()
	err = core.db.
		QueryRow(sqlString).
		Scan(&ownerUserIDNum)
	if err != nil {
		if err == sql.ErrNoRows {
			return iam.UserIDNumZero, nil
		}
		return iam.UserIDNumZero, err
	}

	return
}

// The ID of the user which provided email address is their identifier,
// verified or not.
func (core *Core) getUserIDByKeyEmailAddressAllowUnverifiedInsecure(
	emailAddress email.Address,
) (ownerUserID iam.UserID, alreadyVerified bool, err error) {
	sqlString, _, _ := goqu.
		From(userKeyEmailAddressDBTableName).
		Select(
			"user_id",
			goqu.Case().
				When(goqu.C("verification_ts").IsNull(), false).
				Else(true).
				As("verified"),
		).
		Where(
			goqu.C("domain_part").Eq(emailAddress.DomainPart()),
			goqu.C("local_part").Eq(emailAddress.LocalPart()),
			goqu.C("md_d_ts").IsNull(),
		).
		Order(goqu.C("md_c_ts").Desc()).
		Limit(1).
		ToSQL()

	var ownerUserIDNum iam.UserIDNum
	err = core.db.
		QueryRow(sqlString).
		Scan(&ownerUserIDNum, &alreadyVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return iam.UserIDZero(), false, nil
		}
		return iam.UserIDZero(), false, err
	}

	return iam.NewUserID(ownerUserIDNum), alreadyVerified, nil
}

func (core *Core) SetUserKeyEmailAddress(
	inputCtx iam.CallInputContext,
	userID iam.UserID,
	emailAddress email.Address,
	verificationMethods []eav10n.VerificationMethod,
) (verificationID int64, verificationCodeExpiry *time.Time, err error) {
	ctxAuth := inputCtx.Authorization()
	if !ctxAuth.IsUserSubject() {
		return 0, nil, iam.ErrUserContextRequired
	}
	// Don't allow changing other user's for now
	if !ctxAuth.IsUser(userID) {
		return 0, nil, iam.ErrOperationNotAllowed
	}

	existingOwnerUserIDNum, err := core.
		getUserIDNumByKeyEmailAddressInsecure(emailAddress)
	if err != nil {
		return 0, nil, errors.Wrap("getUserIDNumByKeyEmailAddressInsecure", err)
	}
	if existingOwnerUserIDNum.IsStaticallyValid() {
		if existingOwnerUserIDNum != ctxAuth.UserIDNum() {
			return 0, nil, errors.ArgMsg("emailAddress", "conflict")
		}
		return 0, nil, nil
	}

	alreadyVerified, err := core.setUserKeyEmailAddressInsecure(
		inputCtx, ctxAuth.UserID(), emailAddress)
	if err != nil {
		panic(err)
	}
	if alreadyVerified {
		return 0, nil, nil
	}

	//TODO: user-set has higher priority over terminal's
	userLanguages, err := core.getTerminalAcceptLanguagesAllowDeleted(ctxAuth.TerminalIDNum())
	if err != nil {
	}

	verificationID, verificationCodeExpiry, err = core.eaVerifier.
		StartVerification(inputCtx, emailAddress,
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

func (core *Core) setUserKeyEmailAddressInsecure(
	inputCtx iam.CallInputContext,
	userID iam.UserID,
	emailAddress email.Address,
) (alreadyVerified bool, err error) {
	ctxTime := inputCtx.CallInputMetadata().ReceiveTime
	ctxAuth := inputCtx.Authorization()

	xres, err := core.db.Exec(
		`INSERT INTO `+userKeyEmailAddressDBTableName+` (`+
			`user_id, domain_part, local_part, raw_input, `+
			`md_c_ts, md_c_uid, md_c_tid `+
			`) VALUES (`+
			`$1, $2, $3, $4, $5, $6, $7`+
			`) `+
			`ON CONFLICT (user_id, domain_part, local_part) WHERE md_d_ts IS NULL `+
			`DO NOTHING`,
		userID.IDNum().PrimitiveValue(),
		emailAddress.DomainPart(),
		emailAddress.LocalPart(),
		emailAddress.RawInput(),
		ctxTime,
		ctxAuth.UserIDNum().PrimitiveValue(),
		ctxAuth.TerminalIDNum().PrimitiveValue())
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
		`SELECT CASE WHEN verification_ts IS NULL THEN false ELSE true END AS verified `+
			`FROM `+userKeyEmailAddressDBTableName+` `+
			`WHERE user_id = $1 AND domain_part = $2 AND local_part = $3`,
		userID.IDNum().PrimitiveValue(), emailAddress.DomainPart(), emailAddress.LocalPart()).
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
	inputCtx iam.CallInputContext,
	verificationID int64,
	code string,
) (stateChanged bool, err error) {
	ctxAuth := inputCtx.Authorization()
	err = core.eaVerifier.ConfirmVerification(
		inputCtx, verificationID, code)
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

	ctxTime := inputCtx.CallInputMetadata().ReceiveTime
	stateChanged, err = core.
		ensureUserEmailAddressVerifiedFlag(
			ctxAuth.UserIDNum(), *emailAddress,
			&ctxTime, verificationID)
	if err != nil {
		panic(err)
	}

	return stateChanged, nil
}

func (core *Core) ensureUserEmailAddressVerifiedFlag(
	userIDNum iam.UserIDNum,
	emailAddress email.Address,
	verificationTime *time.Time,
	verificationID int64,
) (stateChanged bool, err error) {
	var xres sql.Result

	xres, err = core.db.Exec(
		`UPDATE `+userKeyEmailAddressDBTableName+` SET (`+
			`verification_ts, verification_id`+
			`) = ( `+
			`$1, $2`+
			`) WHERE user_id = $3 `+
			`AND domain_part = $4 AND local_part = $5 `+
			`AND md_d_ts IS NULL AND verification_ts IS NULL`,
		verificationTime,
		verificationID,
		userIDNum.PrimitiveValue(),
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
