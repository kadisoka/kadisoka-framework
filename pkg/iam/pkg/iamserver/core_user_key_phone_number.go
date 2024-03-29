package iamserver

import (
	"bytes"
	"database/sql"
	"strconv"
	"time"

	"github.com/alloyzeus/go-azfl/errors"
	goqu "github.com/doug-martin/goqu/v9"
	"github.com/lib/pq"

	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iamserver/pnv10n"
	"github.com/kadisoka/kadisoka-framework/pkg/volib/pkg/telephony"
)

// Interface conformance assertion.
var _ iam.UserKeyPhoneNumberService = &Core{}

const userKeyPhoneNumberDBTableName = `user_key_phone_number_dt`

func (core *Core) GetUserKeyPhoneNumber(
	inputCtx iam.CallInputContext,
	userID iam.UserID,
) (*telephony.PhoneNumber, error) {
	//TODO: access control
	return core.getUserKeyPhoneNumberInsecure(inputCtx, userID)
}

func (core *Core) getUserKeyPhoneNumberInsecure(
	inputCtx iam.CallInputContext,
	userID iam.UserID,
) (*telephony.PhoneNumber, error) {
	var countryCode int32
	var nationalNumber int64

	sqlString, _, _ := goqu.
		From(userKeyPhoneNumberDBTableName).
		Select("country_code", "national_number").
		Where(
			goqu.C("user_id").Eq(userID.IDNum().PrimitiveValue()),
			goqu.C("md_d_ts").IsNull(),
			goqu.C("verification_ts").IsNotNull()).
		ToSQL()

	err := core.db.
		QueryRow(sqlString).
		Scan(&countryCode, &nationalNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	phoneNumber := telephony.NewPhoneNumber(countryCode, nationalNumber)
	return &phoneNumber, nil
}

// The ID of the user which provided phone number is their verified primary.
func (core *Core) getUserIDNumByKeyPhoneNumberInsecure(
	phoneNumber telephony.PhoneNumber,
) (ownerUserIDNum iam.UserIDNum, err error) {
	sqlString, _, _ := goqu.
		From(userKeyPhoneNumberDBTableName).
		Select("user_id").
		Where(
			goqu.C("country_code").Eq(phoneNumber.CountryCode()),
			goqu.C("national_number").Eq(phoneNumber.NationalNumber()),
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

// The ID of the user which provided phone number is their identifier,
// verified or not.
func (core *Core) getUserIDByKeyPhoneNumberAllowUnverifiedInsecure(
	phoneNumber telephony.PhoneNumber,
) (ownerUserID iam.UserID, alreadyVerified bool, err error) {
	sqlString, _, _ := goqu.
		From(userKeyPhoneNumberDBTableName).
		Select(
			"user_id",
			goqu.Case().
				When(goqu.C("verification_ts").IsNull(), false).
				Else(true).
				As("verified"),
		).
		Where(
			goqu.C("country_code").Eq(phoneNumber.CountryCode()),
			goqu.C("national_number").Eq(phoneNumber.NationalNumber()),
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

func (core *Core) SetUserKeyPhoneNumber(
	inputCtx iam.CallInputContext,
	userID iam.UserID,
	phoneNumber telephony.PhoneNumber,
	verificationMethods []pnv10n.VerificationMethod,
) (verificationID int64, verificationCodeExpiry *time.Time, err error) {
	ctxAuth := inputCtx.Authorization()
	if !ctxAuth.IsUserSubject() {
		return 0, nil, iam.ErrUserContextRequired
	}
	// Don't allow changing other user's for now
	if !ctxAuth.IsUser(userID) {
		return 0, nil, iam.ErrOperationNotAllowed
	}

	//TODO: prone to race condition. solution: simply call
	// setUserKeyPhoneNumber and translate the error.
	existingOwnerUserIDNum, err := core.
		getUserIDNumByKeyPhoneNumberInsecure(phoneNumber)
	if err != nil {
		return 0, nil, errors.Wrap("getUserIDNumByKeyPhoneNumberInsecure", err)
	}
	if existingOwnerUserIDNum.IsStaticallyValid() {
		if existingOwnerUserIDNum != ctxAuth.UserIDNum() {
			return 0, nil, errors.ArgMsg("phoneNumber", "conflict")
		}
		return 0, nil, nil
	}

	alreadyVerified, err := core.setUserKeyPhoneNumber(
		inputCtx, ctxAuth.UserID(), phoneNumber)
	if err != nil {
		return 0, nil, errors.Wrap("setUserKeyPhoneNumber", err)
	}
	if alreadyVerified {
		return 0, nil, nil
	}

	//TODO: user-set has higher priority over terminal's
	userLanguages, err := core.getTerminalAcceptLanguagesAllowDeleted(ctxAuth.TerminalIDNum())
	if err != nil {
	}

	verificationID, verificationCodeExpiry, err = core.pnVerifier.
		StartVerification(inputCtx, phoneNumber,
			0, userLanguages, nil)
	if err != nil {
		switch err.(type) {
		case pnv10n.InvalidPhoneNumberError:
			return 0, nil, errors.Arg("phoneNumber", err)
		}
		return 0, nil, errors.Wrap("pnVerifier.StartVerification", err)
	}

	return
}

func (core *Core) setUserKeyPhoneNumber(
	inputCtx iam.CallInputContext,
	userID iam.UserID,
	phoneNumber telephony.PhoneNumber,
) (alreadyVerified bool, err error) {
	ctxTime := inputCtx.CallInputMetadata().ReceiveTime
	ctxAuth := inputCtx.Authorization()

	xres, err := core.db.Exec(
		`INSERT INTO `+userKeyPhoneNumberDBTableName+` (`+
			`user_id, country_code, national_number, raw_input, `+
			`md_c_ts, md_c_uid, md_c_tid `+
			`) VALUES (`+
			`$1, $2, $3, $4, $5, $6, $7`+
			`) `+
			`ON CONFLICT (user_id, country_code, national_number) WHERE md_d_ts IS NULL `+
			`DO NOTHING`,
		userID.IDNum().PrimitiveValue(),
		phoneNumber.CountryCode(),
		phoneNumber.NationalNumber(),
		phoneNumber.RawInput(),
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
			`FROM `+userKeyPhoneNumberDBTableName+` `+
			`WHERE user_id = $1 AND country_code = $2 AND national_number = $3`,
		userID.IDNum().PrimitiveValue(), phoneNumber.CountryCode(), phoneNumber.NationalNumber()).
		Scan(&alreadyVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return
}

func (core *Core) ConfirmUserPhoneNumberVerification(
	inputCtx iam.CallInputContext,
	verificationID int64,
	code string,
) (stateChanged bool, err error) {
	ctxAuth := inputCtx.Authorization()
	err = core.pnVerifier.ConfirmVerification(
		inputCtx, verificationID, code)
	if err != nil {
		switch err {
		case pnv10n.ErrVerificationCodeMismatch:
			return false, errors.ArgMsg("code", "mismatch")
		case pnv10n.ErrVerificationCodeExpired:
			return false, errors.ArgMsg("code", "expired")
		}
		return false, errors.Wrap("pnVerifier.ConfirmVerification", err)
	}

	phoneNumber, err := core.pnVerifier.
		GetPhoneNumberByVerificationID(verificationID)
	// An unexpected condition which could cause bad state
	if err != nil {
		panic(err)
	}

	ctxTime := inputCtx.CallInputMetadata().ReceiveTime
	stateChanged, err = core.
		ensureUserPhoneNumberVerifiedFlag(
			ctxAuth.UserIDNum(), *phoneNumber,
			&ctxTime, verificationID)
	if err != nil {
		panic(err)
	}

	return stateChanged, nil
}

// ensureUserPhoneNumberVerifiedFlag is used to ensure that the a user
// phone number is marked as verified. If it has not been verified, this
// method will update the record, otherwise, it does nothing.
//
//TODO: only the verificationID. We'll use it to look up the verification
// data.
func (core *Core) ensureUserPhoneNumberVerifiedFlag(
	userIDNum iam.UserIDNum,
	phoneNumber telephony.PhoneNumber,
	verificationTime *time.Time,
	verificationID int64,
) (stateChanged bool, err error) {
	var xres sql.Result

	xres, err = core.db.Exec(
		`UPDATE `+userKeyPhoneNumberDBTableName+` SET (`+
			`verification_ts, verification_id`+
			`) = ( `+
			`$1, $2`+
			`) WHERE user_id = $3 `+
			`AND country_code = $4 AND national_number = $5 `+
			`AND md_d_ts IS NULL AND verification_ts IS NULL`,
		verificationTime,
		verificationID,
		userIDNum.PrimitiveValue(),
		phoneNumber.CountryCode(),
		phoneNumber.NationalNumber())
	if err != nil {
		pqErr, _ := err.(*pq.Error)
		if pqErr != nil &&
			pqErr.Code == "23505" &&
			pqErr.Constraint == userKeyPhoneNumberDBTableName+`_country_code_national_number_uidx` {
			return false, errors.ArgMsg("phoneNumber", "conflict")
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

func phoneNumberSliceToSQLSetString(pnSlice []telephony.PhoneNumber) string {
	if len(pnSlice) == 0 {
		return ""
	}
	var r bytes.Buffer
	for idx, iv := range pnSlice {
		if idx != 0 {
			r.WriteByte(',')
		}
		r.WriteByte('(')
		r.WriteString(strconv.FormatInt(int64(iv.CountryCode()), 10))
		r.WriteByte(',')
		r.WriteString(strconv.FormatInt(iv.NationalNumber(), 10))
		r.WriteByte(')')
	}
	return r.String()
}
