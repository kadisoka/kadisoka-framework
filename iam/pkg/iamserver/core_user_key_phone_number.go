package iamserver

import (
	"bytes"
	"database/sql"
	"strconv"
	"time"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	goqu "github.com/doug-martin/goqu/v9"
	"github.com/lib/pq"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/pnv10n"
	"github.com/kadisoka/kadisoka-framework/volib/pkg/telephony"
)

// Interface conformance assertion.
var _ iam.UserKeyPhoneNumberService = &Core{}

const userKeyPhoneNumberDBTableName = `user_key_phone_number_dt`

func (core *Core) ListUsersByPhoneNumber(
	callCtx iam.CallContext,
	phoneNumbers []telephony.PhoneNumber,
) ([]iam.UserKeyPhoneNumber, error) {
	if len(phoneNumbers) == 0 {
		return []iam.UserKeyPhoneNumber{}, nil
	}
	ctxAuth := callCtx.Authorization()

	var err error

	// https://dba.stackexchange.com/questions/91247/optimizing-a-postgres-query-with-a-large-in
	userPhoneNumberRows, err := core.db.
		Queryx(
			`SELECT user_id, country_code, national_number ` +
				`FROM ` + userKeyPhoneNumberDBTableName + ` ` +
				`WHERE (country_code, national_number) ` +
				`IN (VALUES ` + phoneNumberSliceToSQLSetString(phoneNumbers) + `) ` +
				`AND d_ts IS NULL AND verification_ts IS NOT NULL ` +
				`LIMIT ` + strconv.Itoa(len(phoneNumbers)))
	if err != nil {
		panic(err)
	}
	defer userPhoneNumberRows.Close()

	userPhoneNumberList := []iam.UserKeyPhoneNumber{}
	for userPhoneNumberRows.Next() {
		userIDNum := iam.UserIDNumZero
		var countryCode int32
		var nationalNumber int64
		err = userPhoneNumberRows.Scan(
			&userIDNum, &countryCode, &nationalNumber)
		if err != nil {
			panic(err)
		}
		userPhoneNumber := iam.UserKeyPhoneNumber{
			UserRef:     iam.NewUserRefKey(userIDNum),
			PhoneNumber: telephony.NewPhoneNumber(countryCode, nationalNumber),
		}
		userPhoneNumberList = append(userPhoneNumberList, userPhoneNumber)
	}
	if err = userPhoneNumberRows.Err(); err != nil {
		panic(err)
	}
	// Already deferred above but we are closing it now because we want to do more queries
	userPhoneNumberRows.Close()

	if ctxAuth.IsUserContext() {
		go func() {
			for _, pn := range phoneNumbers {
				_, err = core.db.Exec(
					`INSERT INTO `+userContactPhoneNumberDBTableName+` (`+
						"user_id, contact_country_code, contact_national_number, "+
						"c_uid, c_tid"+
						") VALUES ($1, $2, $3, $4, $5) "+
						`ON CONFLICT ON CONSTRAINT `+userContactPhoneNumberDBTableName+`_pkey DO NOTHING`,
					ctxAuth.UserIDNum().PrimitiveValue(),
					pn.CountryCode(),
					pn.NationalNumber(),
					ctxAuth.UserIDNum().PrimitiveValue(),
					ctxAuth.TerminalIDNum().PrimitiveValue())
				if err != nil {
					logCtx(callCtx).
						Error().Err(err).Str("phone_number", pn.String()).
						Msg("User contact phone number store")
				}
			}
		}()
	}

	return userPhoneNumberList, nil
}

func (core *Core) GetUserKeyPhoneNumber(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
) (*telephony.PhoneNumber, error) {
	//TODO: access control
	return core.getUserKeyPhoneNumberNoAC(callCtx, userRef)
}

func (core *Core) getUserKeyPhoneNumberNoAC(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
) (*telephony.PhoneNumber, error) {
	var countryCode int32
	var nationalNumber int64

	sqlString, _, _ := goqu.
		From(userKeyPhoneNumberDBTableName).
		Select("country_code", "national_number").
		Where(
			goqu.C("user_id").Eq(userRef.IDNum().PrimitiveValue()),
			goqu.C("d_ts").IsNull(),
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
func (core *Core) getUserIDNumByKeyPhoneNumber(
	phoneNumber telephony.PhoneNumber,
) (ownerUserIDNum iam.UserIDNum, err error) {
	sqlString, _, _ := goqu.
		From(userKeyPhoneNumberDBTableName).
		Select("user_id").
		Where(
			goqu.C("country_code").Eq(phoneNumber.CountryCode()),
			goqu.C("national_number").Eq(phoneNumber.NationalNumber()),
			goqu.C("d_ts").IsNull(),
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
func (core *Core) getUserRefByKeyPhoneNumberAllowUnverified(
	phoneNumber telephony.PhoneNumber,
) (ownerUserRef iam.UserRefKey, alreadyVerified bool, err error) {
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
			goqu.C("d_ts").IsNull(),
		).
		Order(goqu.C("c_ts").Desc()).
		Limit(1).
		ToSQL()

	var ownerUserIDNum iam.UserIDNum
	err = core.db.
		QueryRow(sqlString).
		Scan(&ownerUserIDNum, &alreadyVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return iam.UserRefKeyZero(), false, nil
		}
		return iam.UserRefKeyZero(), false, err
	}

	return iam.NewUserRefKey(ownerUserIDNum), alreadyVerified, nil
}

func (core *Core) SetUserKeyPhoneNumber(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
	phoneNumber telephony.PhoneNumber,
	verificationMethods []pnv10n.VerificationMethod,
) (verificationID int64, verificationCodeExpiry *time.Time, err error) {
	ctxAuth := callCtx.Authorization()
	if !ctxAuth.IsUserContext() {
		return 0, nil, iam.ErrUserContextRequired
	}
	// Don't allow changing other user's for now
	if !ctxAuth.IsUser(userRef) {
		return 0, nil, iam.ErrOperationNotAllowed
	}

	//TODO: prone to race condition. solution: simply call
	// setUserKeyPhoneNumber and translate the error.
	existingOwnerUserIDNum, err := core.
		getUserIDNumByKeyPhoneNumber(phoneNumber)
	if err != nil {
		return 0, nil, errors.Wrap("getUserIDNumByKeyPhoneNumber", err)
	}
	if existingOwnerUserIDNum.IsSound() {
		if existingOwnerUserIDNum != ctxAuth.UserIDNum() {
			return 0, nil, errors.ArgMsg("phoneNumber", "conflict")
		}
		return 0, nil, nil
	}

	alreadyVerified, err := core.setUserKeyPhoneNumber(
		callCtx, ctxAuth.UserRef(), phoneNumber)
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
		StartVerification(callCtx, phoneNumber,
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
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
	phoneNumber telephony.PhoneNumber,
) (alreadyVerified bool, err error) {
	ctxTime := callCtx.RequestInfo().ReceiveTime
	ctxAuth := callCtx.Authorization()

	xres, err := core.db.Exec(
		`INSERT INTO `+userKeyPhoneNumberDBTableName+` (`+
			`user_id, country_code, national_number, raw_input, `+
			`c_ts, c_uid, c_tid `+
			`) VALUES (`+
			`$1, $2, $3, $4, $5, $6, $7`+
			`) `+
			`ON CONFLICT (user_id, country_code, national_number) WHERE d_ts IS NULL `+
			`DO NOTHING`,
		userRef.IDNum().PrimitiveValue(),
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
		userRef.IDNum().PrimitiveValue(), phoneNumber.CountryCode(), phoneNumber.NationalNumber()).
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
	callCtx iam.CallContext,
	verificationID int64,
	code string,
) (stateChanged bool, err error) {
	ctxAuth := callCtx.Authorization()
	err = core.pnVerifier.ConfirmVerification(
		callCtx, verificationID, code)
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

	ctxTime := callCtx.RequestInfo().ReceiveTime
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
			`AND d_ts IS NULL AND verification_ts IS NULL`,
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
