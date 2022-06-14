package iamserver

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/alloyzeus/go-azfl/errors"
	goqu "github.com/doug-martin/goqu/v9"
	"golang.org/x/text/language"

	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iamserver/eav10n"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iamserver/pnv10n"
)

var (
	errTerminalVerificationConfirmationReplayed = errors.EntMsg("terminal verification confirmation", "replayed")
)

func (core *Core) AuthenticateTerminal(
	terminalID iam.TerminalID,
	terminalSecret string,
) (authOK bool, ownerUserID iam.UserID, err error) {
	var storedSecret string
	var ownerUserIDNum iam.UserIDNum

	sqlString, _, _ := goqu.
		From(terminalDBTableName).
		Select("user_id", "secret").
		Where(
			goqu.C("id_num").Eq(terminalID.IDNum().PrimitiveValue()),
			goqu.C("_md_ts").IsNull(),
			goqu.C("verification_ts").IsNotNull(),
		).
		ToSQL()

	err = core.db.
		QueryRow(sqlString).
		Scan(&ownerUserIDNum, &storedSecret)
	if err == sql.ErrNoRows {
		return false, iam.UserIDZero(), nil
	}
	if err != nil {
		return false, iam.UserIDZero(), err
	}

	return subtle.ConstantTimeCompare(
			[]byte(storedSecret),
			[]byte(terminalSecret),
		) == 1,
		iam.NewUserID(ownerUserIDNum), nil
}

func (core *Core) StartTerminalAuthorizationByPhoneNumber(
	inputCtx iam.CallInputContext,
	inputData TerminalAuthorizationByPhoneNumberStartInputData,
) (iam.CallOutputContext, TerminalAuthorizationStartOutputData) {
	ctxAuth := inputCtx.Authorization()
	if ctxAuth.IsStaticallyValid() && !ctxAuth.IsUserSubject() {
		return iam.CallOutputContext{Err: iam.ErrAuthorizationInvalid},
			TerminalAuthorizationStartOutputData{}
	}

	phoneNumber := inputData.PhoneNumber
	if !phoneNumber.IsStaticallyValid() && !core.isTestPhoneNumber(phoneNumber) {
		return iam.CallOutputContext{Err: errors.Arg("phoneNumber", nil)},
			TerminalAuthorizationStartOutputData{}
	}

	// Get the existing owner, whether already verified or not.
	ownerUserID, _, err := core.
		getUserIDByKeyPhoneNumberAllowUnverifiedInsecure(phoneNumber)
	if err != nil {
		panic(err)
	}

	if ownerUserID.IsStaticallyValid() {
		// As the request is authenticated, check if the phone number
		// is associated to the authenticated user.
		if ctxAuth.IsUserSubject() && !ctxAuth.IsUser(ownerUserID) {
			return iam.CallOutputContext{Err: errors.ArgMsg("phoneNumber", "conflict")},
				TerminalAuthorizationStartOutputData{}
		}
	} else {
		newUserID, _, err := core.contextUserOrNewInstance(inputCtx)
		if err != nil {
			panic(err)
		}
		_, err = core.
			setUserKeyPhoneNumber(
				inputCtx, newUserID, phoneNumber)
		if err != nil {
			panic(err)
		}
		ownerUserID = newUserID
	}

	userPreferredLanguages := inputCtx.OriginInfo().AcceptLanguage

	verificationID, verificationCodeExpiryTime, err := core.pnVerifier.
		StartVerification(inputCtx, phoneNumber,
			0, userPreferredLanguages, inputData.VerificationMethods)
	if err != nil {
		switch err.(type) {
		case pnv10n.InvalidPhoneNumberError:
			return iam.CallOutputContext{Err: errors.Arg("phoneNumber", err)},
				TerminalAuthorizationStartOutputData{}
		}
		return iam.CallOutputContext{Err: errors.Wrap("pnVerifier.StartVerification", err)},
			TerminalAuthorizationStartOutputData{}
	}

	regOutCtx, regOutData := core.RegisterTerminal(inputCtx,
		TerminalRegistrationInputData{
			ApplicationID:    inputData.ApplicationID,
			UserID:           ownerUserID,
			DisplayName:      inputData.DisplayName,
			VerificationType: iam.TerminalVerificationResourceTypePhoneNumber,
			VerificationID:   verificationID,
		})
	if regOutCtx.Err != nil {
		panic(regOutCtx.Err)
	}

	return iam.CallOutputContext{
			Mutated: regOutCtx.Mutated,
		}, TerminalAuthorizationStartOutputData{
			TerminalID:                 regOutData.TerminalID,
			VerificationID:             verificationID,
			VerificationCodeExpiryTime: verificationCodeExpiryTime,
		}
}

func (core *Core) StartTerminalAuthorizationByEmailAddress(
	inputCtx iam.CallInputContext,
	inputData TerminalAuthorizationByEmailAddressStartInputData,
) (iam.CallOutputContext, TerminalAuthorizationStartOutputData) {
	ctxAuth := inputCtx.Authorization()
	if ctxAuth.IsStaticallyValid() && !ctxAuth.IsUserSubject() {
		return iam.CallOutputContext{Err: iam.ErrAuthorizationInvalid},
			TerminalAuthorizationStartOutputData{}
	}

	emailAddress := inputData.EmailAddress
	if !emailAddress.IsStaticallyValid() && !core.isTestEmailAddress(emailAddress) {
		return iam.CallOutputContext{Err: errors.Arg("emailAddress", nil)},
			TerminalAuthorizationStartOutputData{}
	}

	// Get the existing owner, whether already verified or not.
	ownerUserID, _, err := core.
		getUserIDByKeyEmailAddressAllowUnverifiedInsecure(emailAddress)
	if err != nil {
		panic(err)
	}

	if ownerUserID.IsStaticallyValid() {
		// Check if it's fully claimed (already verified)
		ownerUserIDNum, err := core.getUserIDNumByKeyEmailAddressInsecure(emailAddress)
		if err != nil {
			panic(err)
		}
		if ownerUserIDNum.IsStaticallyValid() {
			ownerUserID = iam.NewUserID(ownerUserIDNum)
		}
		// As the request is authenticated, check if the phone number
		// is associated to the authenticated user.
		if ctxAuth.IsUserSubject() && !ctxAuth.IsUser(ownerUserID) {
			return iam.CallOutputContext{Err: errors.ArgMsg("emailAddress", "conflict")},
				TerminalAuthorizationStartOutputData{}
		}
	} else {
		newUserID, _, err := core.contextUserOrNewInstance(inputCtx)
		if err != nil {
			panic(err)
		}
		_, err = core.
			setUserKeyEmailAddressInsecure(
				inputCtx, newUserID, emailAddress)
		if err != nil {
			panic(err)
		}
		ownerUserID = newUserID
	}

	userPreferredLanguages := inputCtx.OriginInfo().AcceptLanguage

	verificationID, verificationCodeExpiryTime, err := core.eaVerifier.
		StartVerification(inputCtx, emailAddress,
			0, userPreferredLanguages, inputData.VerificationMethods)
	if err != nil {
		switch err.(type) {
		case eav10n.InvalidEmailAddressError:
			return iam.CallOutputContext{
					Err: errors.Arg("emailAddress", err)},
				TerminalAuthorizationStartOutputData{}
		}
		return iam.CallOutputContext{
				Err: errors.Wrap("eaVerifier.StartVerification", err)},
			TerminalAuthorizationStartOutputData{}
	}

	regOutCtx, regOutData := core.RegisterTerminal(inputCtx,
		TerminalRegistrationInputData{
			ApplicationID:    inputData.ApplicationID,
			UserID:           ownerUserID,
			DisplayName:      inputData.DisplayName,
			VerificationType: iam.TerminalVerificationResourceTypeEmailAddress,
			VerificationID:   verificationID,
		})
	if regOutCtx.Err != nil {
		panic(regOutCtx.Err)
	}

	return iam.CallOutputContext{
			Mutated: regOutCtx.Mutated,
		},
		TerminalAuthorizationStartOutputData{
			TerminalID:                 regOutData.TerminalID,
			VerificationID:             verificationID,
			VerificationCodeExpiryTime: verificationCodeExpiryTime,
		}
}

// ConfirmTerminalAuthorization confirms authorization of a
// terminal by providing the verificationCode which was delivered through
// selected channel when the authorization was created.
func (core *Core) ConfirmTerminalAuthorization(
	inputCtx iam.CallInputContext,
	terminalID iam.TerminalID,
	verificationCode string,
) (terminalSecret string, userID iam.UserID, err error) {
	// The code is verified based on the identifier used when the verification
	// was requested. Each of the implementation required to implement
	// limit the number of failed attempts.

	ctxTime := inputCtx.CallInputMetadata().ReceiveTime

	termData, err := core.getTerminalRaw(terminalID.IDNum())
	if err != nil {
		panic(err)
	}
	if termData == nil {
		return "", iam.UserIDZero(), errors.ArgMsg("terminalID", "reference invalid")
	}
	disallowReplay := false

	if termData.UserIDNum.IsStaticallyValid() {
		termUserIDNum := termData.UserIDNum

		switch termData.VerificationType {
		case iam.TerminalVerificationResourceTypeEmailAddress:
			err = core.eaVerifier.
				ConfirmVerification(
					inputCtx,
					termData.VerificationID,
					verificationCode)
			if err != nil {
				switch err {
				case eav10n.ErrVerificationCodeMismatch:
					return "", iam.UserIDZero(), iam.ErrTerminalVerificationCodeMismatch
				case eav10n.ErrVerificationCodeExpired:
					return "", iam.UserIDZero(), iam.ErrTerminalVerificationCodeExpired
				}
				panic(err)
			}

			emailAddress, err := core.eaVerifier.
				GetEmailAddressByVerificationID(
					termData.VerificationID)
			if err != nil {
				panic(err)
			}

			updated, err := core.
				ensureUserEmailAddressVerifiedFlag(
					termUserIDNum,
					*emailAddress,
					&ctxTime,
					termData.VerificationID)
			if err != nil {
				panic(err)
			}
			if !updated {
				// Let's check if the email address is associated to other user
				existingOwnerUserIDNum, err := core.
					getUserIDNumByKeyEmailAddressInsecure(*emailAddress)
				if err != nil {
					panic(err)
				}
				if existingOwnerUserIDNum.IsStaticallyValid() && existingOwnerUserIDNum != termUserIDNum {
					// The email address has been claimed by another user after
					// the current user requested the link.
					return "", iam.UserIDZero(), iam.ErrTerminalVerificationResourceConflict
				}
			}

		case iam.TerminalVerificationResourceTypePhoneNumber:
			err = core.pnVerifier.
				ConfirmVerification(
					inputCtx,
					termData.VerificationID,
					verificationCode)
			if err != nil {
				switch err {
				case pnv10n.ErrVerificationCodeMismatch:
					return "", iam.UserIDZero(), iam.ErrTerminalVerificationCodeMismatch
				case pnv10n.ErrVerificationCodeExpired:
					return "", iam.UserIDZero(), iam.ErrTerminalVerificationCodeExpired
				}
				panic(err)
			}

			phoneNumber, err := core.pnVerifier.
				GetPhoneNumberByVerificationID(
					termData.VerificationID)
			if err != nil {
				panic(err)
			}

			updated, err := core.
				ensureUserPhoneNumberVerifiedFlag(
					termUserIDNum,
					*phoneNumber,
					&ctxTime,
					termData.VerificationID)
			if err != nil {
				panic(err)
			}
			if !updated {
				// Let's check if the phone number is associated to other user
				existingOwnerUserIDNum, err := core.
					getUserIDNumByKeyPhoneNumberInsecure(*phoneNumber)
				if err != nil {
					panic(err)
				}
				if existingOwnerUserIDNum.IsStaticallyValid() && existingOwnerUserIDNum != termUserIDNum {
					// The phone number has been claimed by another user after
					// the current user requested the link.
					return "", iam.UserIDZero(), iam.ErrTerminalVerificationResourceConflict
				}
			}

		case iam.TerminalVerificationResourceTypeOAuthAuthorizationCode:
			disallowReplay = true

		default:
			panic("Unsupported")
		}
	}

	termSecret, err := core.
		setTerminalVerified(inputCtx, termData.IDNum, disallowReplay)
	if err != nil {
		if err == errTerminalVerificationConfirmationReplayed {
			return "", iam.UserIDZero(), iam.ErrAuthorizationCodeAlreadyClaimed
		}
		panic(err)
	}

	return termSecret, iam.NewUserID(termData.UserIDNum), nil
}

func (core *Core) getTerminalRaw(idNum iam.TerminalIDNum) (*terminalDBRawModel, error) {
	var err error
	var ut terminalDBRawModel

	sqlString, _, _ := goqu.
		From(terminalDBTableName).
		Select(
			"id_num", "application_id", "user_id",
			"_mc_ts", "_mc_uid", "_mc_tid", "_mc_origin_address",
			"display_name", "accept_language",
			"verification_type", "verification_id", "verification_ts").
		Where(
			goqu.C("id_num").Eq(idNum.PrimitiveValue())).
		ToSQL()

	err = core.db.
		QueryRowx(sqlString).
		StructScan(&ut)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &ut, nil
}

func (core *Core) GetTerminalInfo(
	inputCtx iam.CallInputContext,
	terminalIDKey iam.TerminalID,
) (*iam.TerminalInfo, error) {
	if inputCtx == nil {
		return nil, nil
	}
	ctxAuth := inputCtx.Authorization()
	if !ctxAuth.IsUserSubject() {
		return nil, nil
	}

	var ownerUserIDNum iam.UserIDNum
	var displayName string
	var acceptLanguage string

	sqlString, _, _ := goqu.
		From(terminalDBTableName).
		Select("user_id", "display_name", "accept_language").
		Where(
			goqu.C("id_num").Eq(terminalIDKey.IDNum().PrimitiveValue()),
			goqu.C("_md_ts").IsNull(),
		).
		ToSQL()

	err := core.db.
		QueryRow(sqlString).
		Scan(&ownerUserIDNum, &displayName, &acceptLanguage)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if !ctxAuth.UserIDNum().EqualsUserIDNum(ownerUserIDNum) {
		return nil, nil
	}

	tags, _, err := language.ParseAcceptLanguage(acceptLanguage)
	if err != nil {
		return nil, err
	}

	return &iam.TerminalInfo{
		DisplayName:    displayName,
		AcceptLanguage: tags,
	}, nil
}

// RegisterTerminal registers a terminal. This function returns terminal's
// secret if the verification type is set to 'implicit'.
func (core *Core) RegisterTerminal(
	inputCtx iam.CallInputContext,
	inputData TerminalRegistrationInputData,
) (iam.CallOutputContext, TerminalRegistrationOutputData) {
	if inputData.ApplicationID.IsNotStaticallyValid() {
		return iam.CallOutputContext{
				Err: errors.ArgMsg("input", "application ID check", errors.Ent("ApplicationID", nil))},
			TerminalRegistrationOutputData{}
	}

	// Allow zero or a valid user ref.
	if !inputData.UserID.IsZero() && inputData.UserID.IsNotStaticallyValid() {
		return iam.CallOutputContext{
				Err: errors.ArgMsg("input", "user ID check", errors.Ent("Data.UserID", nil))},
			TerminalRegistrationOutputData{}
	}

	appInfo, err := core.ApplicationByID(inputData.ApplicationID)
	if err != nil {
		return iam.CallOutputContext{
				Err: errors.Wrap("ApplicationByID", err)},
			TerminalRegistrationOutputData{}
	}
	if appInfo == nil {
		return iam.CallOutputContext{
				Err: errors.ArgMsg("input", "application check",
					errors.EntMsg("ApplicationID", "reference invalid"))},
			TerminalRegistrationOutputData{}
	}

	if (inputData.UserID.IsStaticallyValid() && !inputData.ApplicationID.IDNum().IsUserAgent()) ||
		(inputData.UserID.IsNotStaticallyValid() && inputData.ApplicationID.IDNum().IsUserAgent()) {
		return iam.CallOutputContext{
				Err: errors.ArgMsg("input", "user and application combination invalid",
					errors.EntMsg("Data.UserID", fmt.Sprintf("statically valid: %v", inputData.UserID.IsStaticallyValid())),
					errors.EntMsg("ApplicationID", fmt.Sprintf("user agent: %v", inputData.ApplicationID.AZIDNum().IsUserAgent())))},
			TerminalRegistrationOutputData{}
	}

	//TODO:SEC:
	// - check verification type against client type
	// - check user ref validity against verification type and client type
	// - assert platform type againts client data

	return core.registerTerminalInsecure(inputCtx, inputData)
}

func (core *Core) registerTerminalInsecure(
	inputCtx iam.CallInputContext,
	inputData TerminalRegistrationInputData,
) (iam.CallOutputContext, TerminalRegistrationOutputData) {
	ctxAuth := inputCtx.Authorization()
	ctxTime := inputCtx.CallInputMetadata().ReceiveTime
	originInfo := inputCtx.OriginInfo()

	var termSecret string
	generateSecret :=
		inputData.VerificationType == iam.TerminalVerificationResourceTypeOAuthClientCredentials ||
			inputData.VerificationType == iam.TerminalVerificationResourceTypeOAuthPassword
	if generateSecret {
		termSecret = core.generateTerminalSecret()
		inputData.VerificationTime = &ctxTime
	} else {
		termSecret = ""
		inputData.VerificationTime = nil
	}

	acceptLangStrings := make([]string, 0, len(originInfo.AcceptLanguage))
	for _, tag := range originInfo.AcceptLanguage {
		acceptLangStrings = append(acceptLangStrings, tag.String())
	}

	//TODO: if id conflict, generate another id and retry
	termID, err := GenerateTerminalIDNum(0)
	if err != nil {
		return iam.CallOutputContext{
				Err: errors.Wrap("ID generation", err)},
			TerminalRegistrationOutputData{}
	}

	sqlString, _, _ := goqu.
		Insert(terminalDBTableName).
		Rows(
			goqu.Record{
				"id_num":             termID.PrimitiveValue(),
				"application_id":     inputData.ApplicationID.IDNum().PrimitiveValue(),
				"user_id":            inputData.UserID.IDNum().PrimitiveValue(),
				"secret":             termSecret,
				"_mc_ts":             ctxTime,
				"_mc_uid":            ctxAuth.UserIDNumPtr(),
				"_mc_tid":            ctxAuth.TerminalIDNumPtr(),
				"_mc_origin_address": originInfo.Address,
				"_mc_origin_env":     originInfo.EnvironmentString,
				"display_name":       strings.TrimSpace(inputData.DisplayName),
				"accept_language":    strings.Join(acceptLangStrings, ","),
				"verification_type":  inputData.VerificationType,
				"verification_id":    inputData.VerificationID,
				"verification_ts":    inputData.VerificationTime,
			}).
		ToSQL()

	_, err = core.db.Exec(sqlString)
	if err != nil {
		return iam.CallOutputContext{
				Err: errors.Wrap("data insert", err)},
			TerminalRegistrationOutputData{}
	}

	terminalID := iam.NewTerminalID(inputData.ApplicationID, inputData.UserID, termID)

	return iam.CallOutputContext{
			Mutated: true,
		}, TerminalRegistrationOutputData{
			TerminalID:     terminalID,
			TerminalSecret: termSecret,
		}
}

func (core *Core) DeleteTerminal(
	inputCtx iam.CallInputContext,
	terminalIDToDelete iam.TerminalID,
) (stateChanged bool, err error) {
	ctxAuth := inputCtx.Authorization()

	//TODO: allow delete any owned terminal
	if !ctxAuth.IsTerminal(terminalIDToDelete) {
		return false, iam.ErrOperationNotAllowed
	}

	ctxTime := inputCtx.CallInputMetadata().ReceiveTime

	sqlString, _, _ := goqu.
		From(terminalDBTableName).
		Where(
			goqu.C("id_num").Eq(terminalIDToDelete.IDNum().PrimitiveValue()),
			goqu.C("_md_ts").IsNull(),
		).
		Update().
		Set(
			goqu.Record{
				"_md_ts":  ctxTime,
				"_md_tid": ctxAuth.TerminalIDNum().PrimitiveValue(),
				"_md_uid": ctxAuth.UserIDNum().PrimitiveValue(),
			},
		).
		ToSQL()

	xres, err := core.db.
		Exec(sqlString)
	if err != nil {
		return false, err
	}
	n, err := xres.RowsAffected()
	if err != nil {
		panic(err)
	}

	if n == 1 {
		//TODO: push the event
	}

	return n == 1, nil
}

//TODO: error if the terminal is deleted?
func (core *Core) setTerminalVerified(
	inputCtx iam.CallInputContext,
	terminalIDNum iam.TerminalIDNum,
	disallowReplay bool,
) (secret string, err error) {
	// A secret is used to obtain an access token. if an access token is
	// expired, the terminal could request another access token by
	// providing this secret. the secret is only provided after the
	// authorization has been verified.
	termSecret := core.generateTerminalSecret() //TODO(exa): JWT (or something similar)

	sqlString, _, _ := goqu.
		From(terminalDBTableName).
		Where(
			goqu.C("id_num").Eq(terminalIDNum.PrimitiveValue()),
			goqu.C("verification_ts").IsNull()).
		Update().
		Set(
			goqu.Record{
				"secret":          termSecret,
				"verification_ts": inputCtx.CallInputMetadata().ReceiveTime,
			}).
		ToSQL()

	xres, err := core.db.
		Exec(sqlString)
	if err != nil {
		return "", err
	}
	n, err := xres.RowsAffected()
	if err != nil {
		panic(err)
	}

	if n == 0 {
		if disallowReplay {
			return "", errTerminalVerificationConfirmationReplayed
		}

		sqlString, _, _ := goqu.
			From(terminalDBTableName).
			Select("secret").
			Where(
				goqu.C("id_num").Eq(terminalIDNum.PrimitiveValue())).
			ToSQL()
		err = core.db.
			QueryRow(sqlString).
			Scan(&termSecret)
		if err != nil {
			panic(err)
		}
	}

	return termSecret, nil
}

func (core *Core) generateTerminalSecret() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func (core *Core) getTerminalAcceptLanguagesAllowDeleted(
	idNum iam.TerminalIDNum,
) ([]language.Tag, error) {
	var acceptLanguage string

	sqlString, _, _ := goqu.
		From(terminalDBTableName).
		Select("accept_language").
		Where(
			goqu.C("id_num").Eq(idNum.PrimitiveValue()),
		).
		ToSQL()

	err := core.db.
		QueryRow(sqlString).
		Scan(&acceptLanguage)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	tags, _, err := language.ParseAcceptLanguage(acceptLanguage)
	if err != nil {
		return nil, err
	}
	return tags, nil
}
