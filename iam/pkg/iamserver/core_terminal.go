package iamserver

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"strings"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	goqu "github.com/doug-martin/goqu/v9"
	"golang.org/x/text/language"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/eav10n"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/pnv10n"
)

var (
	errTerminalVerificationConfirmationReplayed = errors.EntMsg("terminal verification confirmation", "replayed")
)

const terminalDBTableName = "terminal_dt"

func (core *Core) AuthenticateTerminal(
	terminalRef iam.TerminalRefKey,
	terminalSecret string,
) (authOK bool, ownerUserRef iam.UserRefKey, err error) {
	var storedSecret string
	var ownerUserIDNum iam.UserIDNum

	sqlString, _, _ := goqu.
		From(terminalDBTableName).
		Select("user_id", "secret").
		Where(
			goqu.C("id").Eq(terminalRef.IDNum().PrimitiveValue()),
			goqu.C("d_ts").IsNull(),
			goqu.C("verification_ts").IsNotNull(),
		).
		ToSQL()

	err = core.db.
		QueryRow(sqlString).
		Scan(&ownerUserIDNum, &storedSecret)
	if err == sql.ErrNoRows {
		return false, iam.UserRefKeyZero(), nil
	}
	if err != nil {
		return false, iam.UserRefKeyZero(), err
	}

	return subtle.ConstantTimeCompare(
			[]byte(storedSecret),
			[]byte(terminalSecret),
		) == 1,
		iam.NewUserRefKey(ownerUserIDNum), nil
}

func (core *Core) StartTerminalAuthorizationByPhoneNumber(
	input TerminalAuthorizationByPhoneNumberStartInput,
) TerminalAuthorizationStartOutput {
	callCtx := input.Context
	ctxAuth := callCtx.Authorization()
	if ctxAuth.IsValid() && !ctxAuth.IsUserContext() {
		return TerminalAuthorizationStartOutput{
			Context: iam.OpOutputContext{
				Err: iam.ErrAuthorizationInvalid}}
	}

	phoneNumber := input.Data.PhoneNumber
	if !phoneNumber.IsSound() && !core.isTestPhoneNumber(phoneNumber) {
		return TerminalAuthorizationStartOutput{
			Context: iam.OpOutputContext{
				Err: errors.Arg("phoneNumber", nil)}}
	}

	// Get the existing owner, whether already verified or not.
	ownerUserRef, _, err := core.
		getUserRefByKeyPhoneNumberAllowUnverified(phoneNumber)
	if err != nil {
		panic(err)
	}

	if ownerUserRef.IsSound() {
		// As the request is authenticated, check if the phone number
		// is associated to the authenticated user.
		if ctxAuth.IsUserContext() && !ctxAuth.IsUser(ownerUserRef) {
			return TerminalAuthorizationStartOutput{
				Context: iam.OpOutputContext{
					Err: errors.ArgMsg("phoneNumber", "conflict")}}
		}
	} else {
		newUserRef, _, err := core.contextUserOrNewInstance(callCtx)
		if err != nil {
			panic(err)
		}
		_, err = core.
			setUserKeyPhoneNumber(
				callCtx, newUserRef, phoneNumber)
		if err != nil {
			panic(err)
		}
		ownerUserRef = newUserRef
	}

	userPreferredLanguages := input.Context.OriginInfo().AcceptLanguage

	verificationID, verificationCodeExpiryTime, err := core.pnVerifier.
		StartVerification(callCtx, phoneNumber,
			0, userPreferredLanguages, input.Data.VerificationMethods)
	if err != nil {
		switch err.(type) {
		case pnv10n.InvalidPhoneNumberError:
			return TerminalAuthorizationStartOutput{
				Context: iam.OpOutputContext{
					Err: errors.Arg("phoneNumber", err)}}
		}
		return TerminalAuthorizationStartOutput{
			Context: iam.OpOutputContext{
				Err: errors.Wrap("pnVerifier.StartVerification", err)}}
	}

	regOutput := core.RegisterTerminal(TerminalRegistrationInput{
		Context:        callCtx,
		ApplicationRef: input.ApplicationRef,
		Data: TerminalRegistrationInputData{
			UserRef:          ownerUserRef,
			DisplayName:      input.Data.DisplayName,
			VerificationType: iam.TerminalVerificationResourceTypePhoneNumber,
			VerificationID:   verificationID,
		}})
	if regOutput.Context.Err != nil {
		panic(regOutput.Context.Err)
	}

	return TerminalAuthorizationStartOutput{
		Context: iam.OpOutputContext{
			Mutated: true,
		},
		Data: TerminalAuthorizationStartOutputData{
			TerminalRef:                regOutput.Data.TerminalRef,
			VerificationID:             verificationID,
			VerificationCodeExpiryTime: verificationCodeExpiryTime,
		},
	}
}

func (core *Core) StartTerminalAuthorizationByEmailAddress(
	input TerminalAuthorizationByEmailAddressStartInput,
) TerminalAuthorizationStartOutput {
	callCtx := input.Context
	ctxAuth := callCtx.Authorization()
	if ctxAuth.IsValid() && !ctxAuth.IsUserContext() {
		return TerminalAuthorizationStartOutput{
			Context: iam.OpOutputContext{
				Err: iam.ErrAuthorizationInvalid}}
	}

	emailAddress := input.Data.EmailAddress
	if !emailAddress.IsSound() && !core.isTestEmailAddress(emailAddress) {
		return TerminalAuthorizationStartOutput{
			Context: iam.OpOutputContext{
				Err: errors.Arg("emailAddress", nil)}}
	}

	// Get the existing owner, whether already verified or not.
	ownerUserRef, _, err := core.
		getUserRefByKeyEmailAddressAllowUnverified(emailAddress)
	if err != nil {
		panic(err)
	}

	if ownerUserRef.IsSound() {
		// Check if it's fully claimed (already verified)
		ownerUserIDNum, err := core.getUserIDNumByKeyEmailAddress(emailAddress)
		if err != nil {
			panic(err)
		}
		if ownerUserIDNum.IsSound() {
			ownerUserRef = iam.NewUserRefKey(ownerUserIDNum)
		}
		// As the request is authenticated, check if the phone number
		// is associated to the authenticated user.
		if ctxAuth.IsUserContext() && !ctxAuth.IsUser(ownerUserRef) {
			return TerminalAuthorizationStartOutput{
				Context: iam.OpOutputContext{
					Err: errors.ArgMsg("emailAddress", "conflict")}}
		}
	} else {
		newUserRef, _, err := core.contextUserOrNewInstance(callCtx)
		if err != nil {
			panic(err)
		}
		_, err = core.
			setUserKeyEmailAddress(
				callCtx, newUserRef, emailAddress)
		if err != nil {
			panic(err)
		}
		ownerUserRef = newUserRef
	}

	userPreferredLanguages := callCtx.OriginInfo().AcceptLanguage

	verificationID, verificationCodeExpiryTime, err := core.eaVerifier.
		StartVerification(callCtx, emailAddress,
			0, userPreferredLanguages, input.Data.VerificationMethods)
	if err != nil {
		switch err.(type) {
		case eav10n.InvalidEmailAddressError:
			return TerminalAuthorizationStartOutput{
				Context: iam.OpOutputContext{
					Err: errors.Arg("emailAddress", err)}}
		}
		return TerminalAuthorizationStartOutput{
			Context: iam.OpOutputContext{
				Err: errors.Wrap("eaVerifier.StartVerification", err)}}
	}

	regOutput := core.RegisterTerminal(TerminalRegistrationInput{
		Context:        callCtx,
		ApplicationRef: input.ApplicationRef,
		Data: TerminalRegistrationInputData{
			UserRef:          ownerUserRef,
			DisplayName:      input.Data.DisplayName,
			VerificationType: iam.TerminalVerificationResourceTypeEmailAddress,
			VerificationID:   verificationID,
		}})
	if regOutput.Context.Err != nil {
		panic(regOutput.Context.Err)
	}

	return TerminalAuthorizationStartOutput{
		Context: iam.OpOutputContext{
			Mutated: true,
		},
		Data: TerminalAuthorizationStartOutputData{
			TerminalRef:                regOutput.Data.TerminalRef,
			VerificationID:             verificationID,
			VerificationCodeExpiryTime: verificationCodeExpiryTime,
		},
	}
}

// ConfirmTerminalAuthorization confirms authorization of a
// terminal by providing the verificationCode which was delivered through
// selected channel when the authorization was created.
func (core *Core) ConfirmTerminalAuthorization(
	callCtx iam.CallContext,
	terminalRef iam.TerminalRefKey,
	verificationCode string,
) (terminalSecret string, userRef iam.UserRefKey, err error) {
	// The code is verified based on the identifier used when the verification
	// was requested. Each of the implementation required to implement
	// limit the number of failed attempts.

	ctxTime := callCtx.RequestInfo().ReceiveTime

	termData, err := core.getTerminalRaw(terminalRef.IDNum())
	if err != nil {
		panic(err)
	}
	if termData == nil {
		return "", iam.UserRefKeyZero(), errors.ArgMsg("terminalID", "reference invalid")
	}
	disallowReplay := false

	if termData.UserIDNum.IsSound() {
		termUserIDNum := termData.UserIDNum

		switch termData.VerificationType {
		case iam.TerminalVerificationResourceTypeEmailAddress:
			err = core.eaVerifier.
				ConfirmVerification(
					callCtx,
					termData.VerificationID,
					verificationCode)
			if err != nil {
				switch err {
				case eav10n.ErrVerificationCodeMismatch:
					return "", iam.UserRefKeyZero(), iam.ErrTerminalVerificationCodeMismatch
				case eav10n.ErrVerificationCodeExpired:
					return "", iam.UserRefKeyZero(), iam.ErrTerminalVerificationCodeExpired
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
					getUserIDNumByKeyEmailAddress(*emailAddress)
				if err != nil {
					panic(err)
				}
				if existingOwnerUserIDNum.IsSound() && existingOwnerUserIDNum != termUserIDNum {
					// The email address has been claimed by another user after
					// the current user requested the link.
					return "", iam.UserRefKeyZero(), iam.ErrTerminalVerificationResourceConflict
				}
			}

		case iam.TerminalVerificationResourceTypePhoneNumber:
			err = core.pnVerifier.
				ConfirmVerification(
					callCtx,
					termData.VerificationID,
					verificationCode)
			if err != nil {
				switch err {
				case pnv10n.ErrVerificationCodeMismatch:
					return "", iam.UserRefKeyZero(), iam.ErrTerminalVerificationCodeMismatch
				case pnv10n.ErrVerificationCodeExpired:
					return "", iam.UserRefKeyZero(), iam.ErrTerminalVerificationCodeExpired
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
					getUserIDNumByKeyPhoneNumber(*phoneNumber)
				if err != nil {
					panic(err)
				}
				if existingOwnerUserIDNum.IsSound() && existingOwnerUserIDNum != termUserIDNum {
					// The phone number has been claimed by another user after
					// the current user requested the link.
					return "", iam.UserRefKeyZero(), iam.ErrTerminalVerificationResourceConflict
				}
			}

		case iam.TerminalVerificationResourceTypeOAuthAuthorizationCode:
			disallowReplay = true

		default:
			panic("Unsupported")
		}
	}

	termSecret, err := core.
		setTerminalVerified(callCtx, termData.IDNum, disallowReplay)
	if err != nil {
		if err == errTerminalVerificationConfirmationReplayed {
			return "", iam.UserRefKeyZero(), iam.ErrAuthorizationCodeAlreadyClaimed
		}
		panic(err)
	}

	return termSecret, iam.NewUserRefKey(termData.UserIDNum), nil
}

func (core *Core) getTerminalRaw(idNum iam.TerminalIDNum) (*terminalDBRawModel, error) {
	var err error
	var ut terminalDBRawModel

	sqlString, _, _ := goqu.
		From(terminalDBTableName).
		Select(
			"id", "application_id", "user_id",
			"c_ts", "c_uid", "c_tid", "c_origin_address",
			"display_name", "accept_language",
			"verification_type", "verification_id", "verification_ts").
		Where(
			goqu.C("id").Eq(idNum.PrimitiveValue())).
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
	callCtx iam.CallContext,
	terminalRefKey iam.TerminalRefKey,
) (*iam.TerminalInfo, error) {
	if callCtx == nil {
		return nil, nil
	}
	ctxAuth := callCtx.Authorization()
	if !ctxAuth.IsUserContext() {
		return nil, nil
	}

	var ownerUserIDNum iam.UserIDNum
	var displayName string
	var acceptLanguage string

	sqlString, _, _ := goqu.
		From(terminalDBTableName).
		Select("user_id", "display_name", "accept_language").
		Where(
			goqu.C("id").Eq(terminalRefKey.IDNum().PrimitiveValue()),
			goqu.C("d_ts").IsNull(),
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
	input TerminalRegistrationInput,
) TerminalRegistrationOutput {
	if input.ApplicationRef.IsNotSound() {
		return TerminalRegistrationOutput{Context: iam.OpOutputContext{
			Err: errors.Arg("input", nil, errors.Ent("ApplicationRef", nil))}}
	}

	// Allow zero or a valid user ref.
	if !input.Data.UserRef.IsZero() && input.Data.UserRef.IsNotSound() {
		return TerminalRegistrationOutput{Context: iam.OpOutputContext{
			Err: errors.Arg("input", nil, errors.Ent("Data.UserRef", nil))}}
	}

	clientInfo, err := core.ApplicationByRefKey(input.ApplicationRef)
	if err != nil {
		return TerminalRegistrationOutput{Context: iam.OpOutputContext{
			Err: errors.Wrap("ApplicationByRefKey", err)}}
	}
	if clientInfo == nil {
		return TerminalRegistrationOutput{Context: iam.OpOutputContext{
			Err: errors.Arg("input", nil, errors.EntMsg("ApplicationRef", "reference invalid"))}}
	}

	//TODO:
	// - check verification type against client type
	// - check user ref validity against verification type and client type
	// - assert platform type againts client data

	return core.registerTerminalNoAC(input)
}

func (core *Core) registerTerminalNoAC(
	input TerminalRegistrationInput,
) TerminalRegistrationOutput {
	callCtx := input.Context
	ctxAuth := callCtx.Authorization()

	ctxTime := callCtx.RequestInfo().ReceiveTime
	originInfo := callCtx.OriginInfo()

	//var verificationID int64
	var termSecret string
	generateSecret := input.Data.VerificationType == iam.TerminalVerificationResourceTypeOAuthClientCredentials
	if generateSecret {
		termSecret = core.generateTerminalSecret()
		input.Data.VerificationTime = &ctxTime
	} else {
		termSecret = ""
		input.Data.VerificationTime = nil
	}

	acceptLangStrings := make([]string, 0, len(originInfo.AcceptLanguage))
	for _, tag := range originInfo.AcceptLanguage {
		acceptLangStrings = append(acceptLangStrings, tag.String())
	}

	//TODO: if id conflict, generate another id and retry
	termID, err := iam.GenerateTerminalIDNum(0)
	if err != nil {
		return TerminalRegistrationOutput{Context: iam.OpOutputContext{
			Err: errors.Wrap("ID generation", err)}}
	}

	sqlString, _, _ := goqu.
		Insert(terminalDBTableName).
		Rows(
			goqu.Record{
				"id":                termID.PrimitiveValue(),
				"application_id":    input.ApplicationRef.IDNum().PrimitiveValue(),
				"user_id":           input.Data.UserRef.IDNum().PrimitiveValue(),
				"secret":            termSecret,
				"c_ts":              ctxTime,
				"c_uid":             ctxAuth.UserIDNumPtr(),
				"c_tid":             ctxAuth.TerminalIDNumPtr(),
				"c_origin_address":  originInfo.Address,
				"c_origin_env":      originInfo.EnvironmentString,
				"display_name":      strings.TrimSpace(input.Data.DisplayName),
				"accept_language":   strings.Join(acceptLangStrings, ","),
				"verification_type": input.Data.VerificationType,
				"verification_id":   input.Data.VerificationID,
				"verification_ts":   input.Data.VerificationTime,
			}).
		ToSQL()

	_, err = core.db.Exec(sqlString)
	if err != nil {
		return TerminalRegistrationOutput{Context: iam.OpOutputContext{
			Err: errors.Wrap("data insert", err)}}
	}

	terminalRef := iam.NewTerminalRefKey(input.ApplicationRef, input.Data.UserRef, termID)
	if generateSecret {
		return TerminalRegistrationOutput{Data: TerminalRegistrationOutputData{
			TerminalRef:    terminalRef,
			TerminalSecret: termSecret,
		}}
	}
	return TerminalRegistrationOutput{Data: TerminalRegistrationOutputData{
		TerminalRef: terminalRef,
	}}
}

func (core *Core) DeleteTerminal(
	callCtx iam.CallContext,
	termRefToDelete iam.TerminalRefKey,
) (stateChanged bool, err error) {
	ctxAuth := callCtx.Authorization()

	if !ctxAuth.IsTerminal(termRefToDelete) {
		return false, iam.ErrOperationNotAllowed
	}

	ctxTime := callCtx.RequestInfo().ReceiveTime

	sqlString, _, _ := goqu.
		From(terminalDBTableName).
		Where(
			goqu.C("id").Eq(termRefToDelete.IDNum().PrimitiveValue()),
			goqu.C("d_ts").IsNull(),
		).
		Update().
		Set(
			goqu.Record{
				"d_ts":  ctxTime,
				"d_tid": ctxAuth.TerminalIDNum().PrimitiveValue(),
				"d_uid": ctxAuth.UserIDNum().PrimitiveValue(),
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
	callCtx iam.CallContext,
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
			goqu.C("id").Eq(terminalIDNum.PrimitiveValue()),
			goqu.C("verification_ts").IsNull()).
		Update().
		Set(
			goqu.Record{
				"secret":          termSecret,
				"verification_ts": callCtx.RequestInfo().ReceiveTime,
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
				goqu.C("id").Eq(terminalIDNum.PrimitiveValue())).
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
			goqu.C("id").Eq(idNum.PrimitiveValue()),
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
