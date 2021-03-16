package iamserver

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/binary"
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
	var ownerUserID iam.UserID

	sqlString, _, _ := goqu.
		From(terminalDBTableName).
		Select("user_id", "secret").
		Where(
			goqu.C("id").Eq(terminalRef.ID().PrimitiveValue()),
			goqu.C("d_ts").IsNull(),
			goqu.C("verification_ts").IsNotNull(),
		).
		ToSQL()

	err = core.db.
		QueryRow(sqlString).
		Scan(&ownerUserID, &storedSecret)
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
		iam.NewUserRefKey(ownerUserID), nil
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
	if !phoneNumber.IsValid() && !core.isTestPhoneNumber(phoneNumber) {
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

	if ownerUserRef.IsValid() {
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

	userPreferredLanguages := input.Data.UserPreferredLanguages

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

	terminalRef, _, err := core.RegisterTerminal(callCtx,
		TerminalRegistrationInput{
			ApplicationRef:   input.ApplicationRef,
			UserRef:          ownerUserRef,
			DisplayName:      input.Data.DisplayName,
			AcceptLanguage:   userPreferredLanguages,
			VerificationType: iam.TerminalVerificationResourceTypePhoneNumber,
			VerificationID:   verificationID,
		})
	if err != nil {
		panic(err)
	}

	return TerminalAuthorizationStartOutput{
		Context: iam.OpOutputContext{
			Mutated: true,
		},
		Data: TerminalAuthorizationStartOutputData{
			TerminalRef:                terminalRef,
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
	if !emailAddress.IsValid() && !core.isTestEmailAddress(emailAddress) {
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

	if ownerUserRef.IsValid() {
		// Check if it's fully claimed (already verified)
		ownerUserID, err := core.getUserIDByKeyEmailAddress(emailAddress)
		if err != nil {
			panic(err)
		}
		if ownerUserID.IsValid() {
			ownerUserRef = iam.NewUserRefKey(ownerUserID)
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

	userPreferredLanguages := input.Data.UserPreferredLanguages

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

	terminalRef, _, err := core.RegisterTerminal(callCtx,
		TerminalRegistrationInput{
			ApplicationRef:   input.ApplicationRef,
			UserRef:          ownerUserRef,
			DisplayName:      input.Data.DisplayName,
			AcceptLanguage:   userPreferredLanguages,
			VerificationType: iam.TerminalVerificationResourceTypeEmailAddress,
			VerificationID:   verificationID,
		})
	if err != nil {
		panic(err)
	}

	return TerminalAuthorizationStartOutput{
		Context: iam.OpOutputContext{
			Mutated: true,
		},
		Data: TerminalAuthorizationStartOutputData{
			TerminalRef:                terminalRef,
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

	termData, err := core.getTerminalRaw(terminalRef.ID())
	if err != nil {
		panic(err)
	}
	if termData == nil {
		return "", iam.UserRefKeyZero(), errors.ArgMsg("terminalID", "reference invalid")
	}
	disallowReplay := false

	if termData.UserID.IsValid() {
		termUserID := termData.UserID

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
					termUserID,
					*emailAddress,
					&ctxTime,
					termData.VerificationID)
			if err != nil {
				panic(err)
			}
			if !updated {
				// Let's check if the email address is associated to other user
				existingOwnerUserID, err := core.
					getUserIDByKeyEmailAddress(*emailAddress)
				if err != nil {
					panic(err)
				}
				if existingOwnerUserID.IsValid() && existingOwnerUserID != termUserID {
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
					termUserID,
					*phoneNumber,
					&ctxTime,
					termData.VerificationID)
			if err != nil {
				panic(err)
			}
			if !updated {
				// Let's check if the phone number is associated to other user
				existingOwnerUserID, err := core.
					getUserIDByKeyPhoneNumber(*phoneNumber)
				if err != nil {
					panic(err)
				}
				if existingOwnerUserID.IsValid() && existingOwnerUserID != termUserID {
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
		setTerminalVerified(callCtx, termData.ID, disallowReplay)
	if err != nil {
		if err == errTerminalVerificationConfirmationReplayed {
			return "", iam.UserRefKeyZero(), iam.ErrAuthorizationCodeAlreadyClaimed
		}
		panic(err)
	}

	return termSecret, iam.NewUserRefKey(termData.UserID), nil
}

func (core *Core) getTerminalRaw(id iam.TerminalID) (*terminalDBRawModel, error) {
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
			goqu.C("id").Eq(id.PrimitiveValue())).
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
	terminalID iam.TerminalID,
) (*iam.TerminalInfo, error) {
	if callCtx == nil {
		return nil, nil
	}
	ctxAuth := callCtx.Authorization()
	if !ctxAuth.IsUserContext() {
		return nil, nil
	}

	var ownerUserID iam.UserID
	var displayName string
	var acceptLanguage string

	sqlString, _, _ := goqu.
		From(terminalDBTableName).
		Select("user_id", "display_name", "accept_language").
		Where(
			goqu.C("id").Eq(terminalID.PrimitiveValue()),
			goqu.C("d_ts").IsNull(),
		).
		ToSQL()

	err := core.db.
		QueryRow(sqlString).
		Scan(&ownerUserID, &displayName, &acceptLanguage)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if !ctxAuth.UserID().EqualsUserID(ownerUserID) {
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
	callCtx iam.CallContext,
	input TerminalRegistrationInput,
) (terminalRef iam.TerminalRefKey, secret string, err error) {
	ctxAuth := callCtx.Authorization()

	if input.ApplicationRef.IsNotValid() {
		return iam.TerminalRefKeyZero(), "", errors.Arg("input.ClientID", nil)
	}
	// Allow zero or a valid user ref.
	if !input.UserRef.IsZero() && input.UserRef.IsNotValid() {
		return iam.TerminalRefKeyZero(), "", errors.Arg("input.UserID", nil)
	}

	clientInfo, err := core.ApplicationByRefKey(input.ApplicationRef)
	if err != nil {
		return iam.TerminalRefKeyZero(), "", errors.ArgWrap("input.ClientID", "lookup", err)
	}
	if clientInfo == nil {
		return iam.TerminalRefKeyZero(), "", errors.ArgMsg("input.ClientID", "reference invalid")
	}

	//TODO:
	// - check verification type against client type
	// - assert platform type againts client data

	ctxTime := callCtx.RequestInfo().ReceiveTime
	originInfo := callCtx.OriginInfo()

	//var verificationID int64
	var termSecret string
	generateSecret := input.VerificationType == iam.TerminalVerificationResourceTypeOAuthClientCredentials
	if generateSecret {
		termSecret = core.generateTerminalSecret()
		input.VerificationTime = &ctxTime
	} else {
		termSecret = ""
		input.VerificationTime = nil
	}

	termLangStrings := make([]string, 0, len(input.AcceptLanguage))
	for _, tag := range input.AcceptLanguage {
		termLangStrings = append(termLangStrings, tag.String())
	}

	//TODO: if id conflict, generate another id and retry
	termID, err := core.generateTerminalID()

	sqlString, _, _ := goqu.
		Insert(terminalDBTableName).
		Rows(
			goqu.Record{
				"id":                termID.PrimitiveValue(),
				"application_id":    input.ApplicationRef.ID().PrimitiveValue(),
				"user_id":           input.UserRef.ID().PrimitiveValue(),
				"secret":            termSecret,
				"c_ts":              ctxTime,
				"c_uid":             ctxAuth.UserIDPtr(),
				"c_tid":             ctxAuth.TerminalIDPtr(),
				"c_origin_address":  originInfo.Address,
				"c_origin_env":      originInfo.EnvironmentString,
				"display_name":      strings.TrimSpace(input.DisplayName),
				"accept_language":   strings.Join(termLangStrings, ","),
				"verification_type": input.VerificationType,
				"verification_id":   input.VerificationID,
				"verification_ts":   input.VerificationTime,
			}).
		ToSQL()

	_, err = core.db.Exec(sqlString)
	if err != nil {
		return iam.TerminalRefKeyZero(), "", err
	}

	terminalRef = iam.NewTerminalRefKey(input.ApplicationRef, input.UserRef, termID)
	if generateSecret {
		return terminalRef, termSecret, nil
	}
	return terminalRef, "", nil
}

func (core *Core) DeleteTerminal(
	callCtx iam.CallContext,
	termRefToDelete iam.TerminalRefKey,
) (stateChanged bool, err error) {
	authCtx := callCtx.Authorization()

	if !authCtx.IsTerminal(termRefToDelete) {
		return false, iam.ErrOperationNotAllowed
	}

	ctxTime := callCtx.RequestInfo().ReceiveTime

	sqlString, _, _ := goqu.
		From(terminalDBTableName).
		Where(
			goqu.C("id").Eq(termRefToDelete.ID().PrimitiveValue()),
			goqu.C("d_ts").IsNull(),
		).
		Update().
		Set(
			goqu.Record{
				"d_ts":  ctxTime,
				"d_tid": authCtx.TerminalID().PrimitiveValue(),
				"d_uid": authCtx.UserID().PrimitiveValue(),
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
	terminalID iam.TerminalID,
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
			goqu.C("id").Eq(terminalID.PrimitiveValue()),
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
				goqu.C("id").Eq(terminalID.PrimitiveValue())).
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

func (core *Core) generateTerminalID() (iam.TerminalID, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b[2:])
	if err != nil {
		panic(err)
	}
	h := binary.BigEndian.Uint64(b) & iam.TerminalIDSignificantBitsMask
	return iam.TerminalIDFromPrimitiveValue(int64(h)), nil
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
	id iam.TerminalID,
) ([]language.Tag, error) {
	var acceptLanguage string

	sqlString, _, _ := goqu.
		From(terminalDBTableName).
		Select("accept_language").
		Where(
			goqu.C("id").Eq(id.PrimitiveValue()),
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
