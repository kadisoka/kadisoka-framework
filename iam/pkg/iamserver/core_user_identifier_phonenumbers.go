package iamserver

import (
	"bytes"
	"database/sql"
	"strconv"
	"time"

	"github.com/alloyzeus/go-azcore/azcore/errors"
	"github.com/lib/pq"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/pnv10n"
)

const userIdentifierPhoneNumberTableName = `user_identifier_phone_numbers`

func (core *Core) ListUsersByPhoneNumber(
	callCtx iam.CallContext,
	phoneNumbers []iam.PhoneNumber,
) ([]iam.UserIdentifierPhoneNumber, error) {
	if len(phoneNumbers) == 0 {
		return []iam.UserIdentifierPhoneNumber{}, nil
	}
	authCtx := callCtx.Authorization()

	var err error

	// https://dba.stackexchange.com/questions/91247/optimizing-a-postgres-query-with-a-large-in
	userPhoneNumberRows, err := core.db.
		Queryx(
			`SELECT user_id, country_code, national_number ` +
				`FROM ` + userIdentifierPhoneNumberTableName + ` ` +
				`WHERE (country_code, national_number) ` +
				`IN (VALUES ` + phoneNumberSliceToSQLSetString(phoneNumbers) + `) ` +
				`AND deletion_time IS NULL AND verification_time IS NOT NULL ` +
				`LIMIT ` + strconv.Itoa(len(phoneNumbers)))
	if err != nil {
		panic(err)
	}
	defer userPhoneNumberRows.Close()

	userPhoneNumberList := []iam.UserIdentifierPhoneNumber{}
	for userPhoneNumberRows.Next() {
		var userPhoneNumber iam.UserIdentifierPhoneNumber
		var countryCode int32
		var nationalNumber int64
		err = userPhoneNumberRows.Scan(
			&userPhoneNumber.UserID, &countryCode, &nationalNumber)
		if err != nil {
			panic(err)
		}
		userPhoneNumber.PhoneNumber = iam.NewPhoneNumber(countryCode, nationalNumber)
		userPhoneNumberList = append(userPhoneNumberList, userPhoneNumber)
	}
	if err = userPhoneNumberRows.Err(); err != nil {
		panic(err)
	}
	// Already deferred above but we are closing it now because we want to do more queries
	userPhoneNumberRows.Close()

	if authCtx.IsUserContext() {
		go func() {
			for _, pn := range phoneNumbers {
				_, err = core.db.Exec(
					"INSERT INTO user_contact_phone_numbers ("+
						"user_id, contact_country_code, contact_national_number, "+
						"creation_user_id, creation_terminal_id"+
						") VALUES ($1, $2, $3, $4, $5) "+
						"ON CONFLICT ON CONSTRAINT user_contact_phone_numbers_pkey DO NOTHING",
					authCtx.UserID, pn.CountryCode(), pn.NationalNumber(), authCtx.UserID, authCtx.TerminalID())
				if err != nil {
					logCtx(callCtx).Err(err).Str("phone_number", pn.String()).
						Msg("User contact phone number store")
				}
			}
		}()
	}

	return userPhoneNumberList, nil
}

//TODO: allow non-verified (let the caller decide with the status)
// there should be getters for different purpose (e.g.,
// for login, for display, for notification, for recovery, etc)
func (core *Core) GetUserIdentifierPhoneNumber(
	callCtx iam.CallContext,
	userID iam.UserID,
) (*iam.PhoneNumber, error) {
	var countryCode int32
	var nationalNumber int64
	err := core.db.
		QueryRow(
			`SELECT country_code, national_number `+
				`FROM `+userIdentifierPhoneNumberTableName+` `+
				`WHERE user_id=$1 `+
				`AND deletion_time IS NULL AND verification_time IS NOT NULL`,
			userID).
		Scan(&countryCode, &nationalNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	phoneNumber := iam.NewPhoneNumber(countryCode, nationalNumber)
	return &phoneNumber, nil
}

// The ID of the user which provided phone number is their verified primary.
func (core *Core) getUserIDByIdentifierPhoneNumber(
	phoneNumber iam.PhoneNumber,
) (ownerUserID iam.UserID, err error) {
	queryStr :=
		`SELECT user_id ` +
			`FROM ` + userIdentifierPhoneNumberTableName + ` ` +
			`WHERE country_code = $1 AND national_number = $2 ` +
			`AND deletion_time IS NULL ` +
			`AND verification_time IS NOT NULL`
	err = core.db.
		QueryRow(queryStr,
			phoneNumber.CountryCode(),
			phoneNumber.NationalNumber()).
		Scan(&ownerUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return iam.UserIDZero, nil
		}
		return iam.UserIDZero, err
	}

	return
}

// The ID of the user which provided phone number is their primary,
// verified or not.
func (core *Core) getUserIDByIdentifierPhoneNumberAllowUnverified(
	phoneNumber iam.PhoneNumber,
) (ownerUserID iam.UserID, verified bool, err error) {
	queryStr :=
		`SELECT user_id, CASE WHEN verification_time IS NULL THEN false ELSE true END AS verified ` +
			`FROM ` + userIdentifierPhoneNumberTableName + ` ` +
			`WHERE country_code = $1 AND national_number = $2 ` +
			`AND deletion_time IS NULL ` +
			`ORDER BY creation_time DESC LIMIT 1`
	err = core.db.
		QueryRow(queryStr,
			phoneNumber.CountryCode(),
			phoneNumber.NationalNumber()).
		Scan(&ownerUserID, &verified)
	if err != nil {
		if err == sql.ErrNoRows {
			return iam.UserIDZero, false, nil
		}
		return iam.UserIDZero, false, err
	}

	return
}

func (core *Core) SetUserIdentifierPhoneNumber(
	callCtx iam.CallContext,
	userID iam.UserID,
	phoneNumber iam.PhoneNumber,
	verificationMethods []pnv10n.VerificationMethod,
) (verificationID int64, codeExpiry *time.Time, err error) {
	authCtx := callCtx.Authorization()
	if !authCtx.IsUserContext() {
		return 0, nil, iam.ErrUserContextRequired
	}
	// Don't allow changing other user's for now
	if userID != authCtx.UserID {
		return 0, nil, iam.ErrContextUserNotAllowedToPerformActionOnResource
	}

	//TODO: prone to race condition. solution: simply call
	// setUserIdentifierPhoneNumber and translate the error.
	existingOwnerUserID, err := core.
		getUserIDByIdentifierPhoneNumber(phoneNumber)
	if err != nil {
		return 0, nil, errors.Wrap("getUserIDByIdentifierPhoneNumber", err)
	}
	if existingOwnerUserID.IsValid() {
		if existingOwnerUserID != authCtx.UserID {
			return 0, nil, errors.ArgMsg("phoneNumber", "conflict")
		}
		return 0, nil, nil
	}

	alreadyVerified, err := core.setUserIdentifierPhoneNumber(
		authCtx.Actor(), authCtx.UserID, phoneNumber)
	if err != nil {
		return 0, nil, errors.Wrap("setUserIdentifierPhoneNumber", err)
	}
	if alreadyVerified {
		return 0, nil, nil
	}

	//TODO: user-set has higher priority over terminal's
	userLanguages, err := core.getTerminalAcceptLanguages(authCtx.TerminalID())

	verificationID, codeExpiry, err = core.pnVerifier.
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

func (core *Core) setUserIdentifierPhoneNumber(
	actor iam.Actor,
	userID iam.UserID,
	phoneNumber iam.PhoneNumber,
) (alreadyVerified bool, err error) {
	tNow := time.Now().UTC()

	xres, err := core.db.Exec(
		`INSERT INTO `+userIdentifierPhoneNumberTableName+` (`+
			`user_id, country_code, national_number, raw_input, `+
			`creation_time, creation_user_id, creation_terminal_id `+
			`) VALUES (`+
			`$1, $2, $3, $4, $5, $6, $7`+
			`) `+
			`ON CONFLICT (user_id, country_code, national_number) WHERE deletion_time IS NULL `+
			`DO NOTHING`,
		userID,
		phoneNumber.CountryCode(),
		phoneNumber.NationalNumber(),
		phoneNumber.RawInput(),
		tNow,
		actor.UserID,
		actor.TerminalID)
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
			`FROM `+userIdentifierPhoneNumberTableName+` `+
			`WHERE user_id = $1 AND country_code = $2 AND national_number = $3`,
		userID, phoneNumber.CountryCode(), phoneNumber.NationalNumber()).
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
) (updated bool, err error) {
	authCtx := callCtx.Authorization()
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

	tNow := time.Now().UTC()
	phoneNumber, err := core.pnVerifier.
		GetPhoneNumberByVerificationID(verificationID)
	// An unexpected condition which could cause bad state
	if err != nil {
		panic(err)
	}

	updated, err = core.
		ensureUserPhoneNumberVerifiedFlag(
			authCtx.UserID, *phoneNumber,
			&tNow, verificationID)
	if err != nil {
		panic(err)
	}

	return updated, nil
}

// ensureUserPhoneNumberVerifiedFlag is used to ensure that the a user
// phone number is marked as verified. If it has not been verified, this
// method will update the record, otherwise, it does nothing.
//
//TODO: only the verificationID. We'll use it to look up the verification
// data.
func (core *Core) ensureUserPhoneNumberVerifiedFlag(
	userID iam.UserID,
	phoneNumber iam.PhoneNumber,
	verificationTime *time.Time,
	verificationID int64,
) (bool, error) {
	var err error
	var xres sql.Result

	xres, err = core.db.Exec(
		`UPDATE `+userIdentifierPhoneNumberTableName+` SET (`+
			`verification_time, verification_id`+
			`) = ( `+
			`$1, $2`+
			`) WHERE user_id = $3 `+
			`AND country_code = $4 AND national_number = $5 `+
			`AND deletion_time IS NULL AND verification_time IS NULL`,
		verificationTime,
		verificationID,
		userID,
		phoneNumber.CountryCode(),
		phoneNumber.NationalNumber())
	if err != nil {
		pqErr, _ := err.(*pq.Error)
		if pqErr != nil &&
			pqErr.Code == "23505" &&
			pqErr.Constraint == userIdentifierPhoneNumberTableName+`_country_code_national_number_uidx` {
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

func phoneNumberSliceToSQLSetString(pnSlice []iam.PhoneNumber) string {
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