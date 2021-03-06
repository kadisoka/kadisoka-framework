package iamserver

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"strings"
	"time"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	"golang.org/x/text/language"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/eav10n"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/pnv10n"
)

var (
	errTerminalVerificationConfirmationReplayed = errors.EntMsg("terminal verification confirmation", "replayed")
)

const terminalTableName = "terminal_dt"

func (core *Core) AuthenticateTerminal(
	terminalRef iam.TerminalRefKey,
	terminalSecret string,
) (authOK bool, ownerUserRef iam.UserRefKey, err error) {
	var storedSecret string
	var ownerUserID iam.UserID
	err = core.db.
		QueryRow(
			`SELECT user_id, secret `+
				`FROM `+terminalTableName+` `+
				`WHERE id=$1`,
			terminalRef.ID().PrimitiveValue()).
		Scan(&ownerUserID, &storedSecret)
	if err == sql.ErrNoRows {
		return false, iam.UserRefKeyZero(), nil
	}
	if err != nil {
		return false, iam.UserRefKeyZero(), err
	}

	return storedSecret == terminalSecret, iam.NewUserRefKey(ownerUserID), nil
}

func (core *Core) StartTerminalRegistrationByPhoneNumber(
	callCtx iam.CallContext,
	appRef iam.ApplicationRefKey,
	phoneNumber iam.PhoneNumber,
	displayName string,
	userAgentString string,
	userPreferredLanguages []language.Tag,
	verificationMethods []pnv10n.VerificationMethod,
) (terminalRef iam.TerminalRefKey, verificationID int64, codeExpiry *time.Time, err error) {
	authCtx := callCtx.Authorization()
	if authCtx.IsValid() && !authCtx.IsUserContext() {
		return iam.TerminalRefKeyZero(), 0, nil, iam.ErrAuthorizationInvalid
	}

	if !phoneNumber.IsValid() && !core.isTestPhoneNumber(phoneNumber) {
		return iam.TerminalRefKeyZero(), 0, nil, errors.Arg("phoneNumber", nil)
	}

	// Get the existing owner, whether already verified or not.
	existingOwnerUserRef, _, err := core.
		getUserIDByKeyPhoneNumberAllowUnverified(phoneNumber)
	if err != nil {
		panic(err)
	}

	if existingOwnerUserRef.IsValid() {
		// As the request is authenticated, check if the phone number
		// is associated to the authenticated user.
		if authCtx.IsUserContext() && !existingOwnerUserRef.EqualsUserRefKey(authCtx.UserRef()) {
			return iam.TerminalRefKeyZero(), 0, nil, errors.ArgMsg("phoneNumber", "conflict")
		}
	} else {
		newUserRef, err := core.
			ensureOrNewUserRef(callCtx, authCtx.UserRef())
		if err != nil {
			panic(err)
		}
		_, err = core.
			setUserKeyPhoneNumber(
				callCtx, newUserRef, phoneNumber)
		if err != nil {
			panic(err)
		}
		existingOwnerUserRef = newUserRef
	}

	ownerUserID := existingOwnerUserRef
	verificationID, codeExpiry, err = core.pnVerifier.
		StartVerification(callCtx, phoneNumber,
			0, userPreferredLanguages, verificationMethods)
	if err != nil {
		switch err.(type) {
		case pnv10n.InvalidPhoneNumberError:
			return iam.TerminalRefKeyZero(), 0, nil, errors.Arg("phoneNumber", err)
		}
		return iam.TerminalRefKeyZero(), 0, nil,
			errors.Wrap("pnVerifier.StartVerification", err)
	}

	termLangStrings := make([]string, 0, len(userPreferredLanguages))
	for _, tag := range userPreferredLanguages {
		termLangStrings = append(termLangStrings, tag.String())
	}

	displayName = strings.TrimSpace(displayName)

	terminalRef, _, err = core.RegisterTerminal(callCtx,
		TerminalRegistrationInput{
			ApplicationRef:   appRef,
			UserRef:          ownerUserID,
			DisplayName:      displayName,
			AcceptLanguage:   strings.Join(termLangStrings, ","),
			VerificationType: iam.TerminalVerificationResourceTypePhoneNumber,
			VerificationID:   verificationID,
		})
	if err != nil {
		panic(err)
	}

	return
}

func (core *Core) StartTerminalRegistrationByEmailAddress(
	callCtx iam.CallContext,
	appRef iam.ApplicationRefKey,
	emailAddress iam.EmailAddress,
	displayName string,
	userAgentString string,
	userPreferredLanguages []language.Tag,
	verificationMethods []eav10n.VerificationMethod,
) (terminalRef iam.TerminalRefKey, verificationID int64, codeExpiry *time.Time, err error) {
	authCtx := callCtx.Authorization()
	if authCtx.IsValid() && !authCtx.IsUserContext() {
		return iam.TerminalRefKeyZero(), 0, nil, iam.ErrAuthorizationInvalid
	}

	if !emailAddress.IsValid() && !core.isTestEmailAddress(emailAddress) {
		return iam.TerminalRefKeyZero(), 0, nil, errors.Arg("emailAddress", nil)
	}

	// Get the existing owner, whether already verified or not.
	existingOwnerUserRef, _, err := core.
		getUserIDByKeyEmailAddressAllowUnverified(emailAddress)
	if err != nil {
		panic(err)
	}

	if existingOwnerUserRef.IsValid() {
		// As the request is authenticated, check if the phone number
		// is associated to the authenticated user.
		if authCtx.IsUserContext() && !existingOwnerUserRef.EqualsUserRefKey(authCtx.UserRef()) {
			return iam.TerminalRefKeyZero(), 0, nil, errors.ArgMsg("emailAddress", "conflict")
		}
	} else {
		newUserRef, err := core.
			ensureOrNewUserRef(callCtx, authCtx.UserRef())
		if err != nil {
			panic(err)
		}
		_, err = core.
			setUserKeyEmailAddress(
				callCtx, newUserRef, emailAddress)
		if err != nil {
			panic(err)
		}
		existingOwnerUserRef = newUserRef
	}

	ownerUserID := existingOwnerUserRef
	verificationID, codeExpiry, err = core.eaVerifier.
		StartVerification(callCtx, emailAddress,
			0, userPreferredLanguages, verificationMethods)
	if err != nil {
		switch err.(type) {
		case eav10n.InvalidEmailAddressError:
			return iam.TerminalRefKeyZero(), 0, nil,
				errors.Arg("emailAddress", err)
		}
		return iam.TerminalRefKeyZero(), 0, nil,
			errors.Wrap("eaVerifier.StartVerification", err)
	}

	termLangStrings := make([]string, 0, len(userPreferredLanguages))
	for _, tag := range userPreferredLanguages {
		termLangStrings = append(termLangStrings, tag.String())
	}

	displayName = strings.TrimSpace(displayName)

	terminalRef, _, err = core.RegisterTerminal(callCtx,
		TerminalRegistrationInput{
			ApplicationRef:   appRef,
			UserRef:          ownerUserID,
			DisplayName:      displayName,
			AcceptLanguage:   strings.Join(termLangStrings, ","),
			VerificationType: iam.TerminalVerificationResourceTypeEmailAddress,
			VerificationID:   verificationID,
		})
	if err != nil {
		panic(err)
	}

	return
}

// ConfirmTerminalRegistrationVerification confirms authorization of a
// terminal by providing the verificationCode which was delivered through
// selected channel when the authorization was created.
func (core *Core) ConfirmTerminalRegistrationVerification(
	callCtx iam.CallContext,
	terminalRef iam.TerminalRefKey,
	verificationCode string,
) (secret string, userRef iam.UserRefKey, err error) {
	// The code is verified based on the identifier used when the verification
	// was requested. Each of the implementation required to implement
	// limit the number of failed attempts.

	ctxTime := callCtx.RequestReceiveTime()

	userTermModel, err := core.getTerminal(terminalRef.ID())
	if err != nil {
		panic(err)
	}
	if userTermModel == nil {
		return "", iam.UserRefKeyZero(), errors.ArgMsg("terminalID", "reference invalid")
	}
	disallowReplay := false

	if userTermModel.UserID.IsValid() {
		terminalUserID := userTermModel.UserID
		switch userTermModel.VerificationType {
		case iam.TerminalVerificationResourceTypeEmailAddress:
			err = core.eaVerifier.
				ConfirmVerification(
					callCtx,
					userTermModel.VerificationID,
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
					userTermModel.VerificationID)
			if err != nil {
				panic(err)
			}

			updated, err := core.
				ensureUserEmailAddressVerifiedFlag(
					terminalUserID,
					*emailAddress,
					&ctxTime,
					userTermModel.VerificationID)
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
				if existingOwnerUserID.IsValid() && existingOwnerUserID != terminalUserID {
					// The email address has been claimed by another user after
					// the current user requested the link.
					return "", iam.UserRefKeyZero(), iam.ErrTerminalVerificationResourceConflict
				}
			}

		case iam.TerminalVerificationResourceTypePhoneNumber:
			err = core.pnVerifier.
				ConfirmVerification(
					callCtx,
					userTermModel.VerificationID,
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
					userTermModel.VerificationID)
			if err != nil {
				panic(err)
			}

			updated, err := core.
				ensureUserPhoneNumberVerifiedFlag(
					terminalUserID,
					*phoneNumber,
					&ctxTime,
					userTermModel.VerificationID)
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
				if existingOwnerUserID.IsValid() && existingOwnerUserID != terminalUserID {
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
		setUserTerminalVerified(callCtx, userTermModel.ID, disallowReplay)
	if err != nil {
		if err == errTerminalVerificationConfirmationReplayed {
			return "", iam.UserRefKeyZero(), iam.ErrAuthorizationCodeAlreadyClaimed
		}
		panic(err)
	}

	return termSecret, iam.NewUserRefKey(userTermModel.UserID), nil
}

func (core *Core) getTerminal(id iam.TerminalID) (*terminalDBModel, error) {
	var err error
	var ut terminalDBModel
	err = core.db.QueryRowx(
		`SELECT `+
			`id, application_id, user_id, secret, `+
			`c_ts, c_uid, c_tid, c_origin_address, `+
			`display_name, accept_language, `+
			`verification_type, verification_id, verification_time `+
			`FROM `+terminalTableName+` `+
			`WHERE id = $1`,
		id).StructScan(&ut)
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
	authCtx := callCtx.Authorization()
	if !authCtx.IsUserContext() {
		return nil, nil
	}

	var ownerUserID iam.UserID
	var displayName string
	var acceptLanguage string
	err := core.db.QueryRow(
		`SELECT user_id, display_name, accept_language `+
			`FROM `+terminalTableName+` `+
			`WHERE id = $1`,
		terminalID).
		Scan(&ownerUserID, &displayName, &acceptLanguage)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if ownerUserID != authCtx.UserID() {
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
	authCtx := callCtx.Authorization()

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

	ctxTime := callCtx.RequestReceiveTime()

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

	//TODO: if id conflict, generate another id and retry
	termID, err := core.generateTerminalID()
	_, err = core.db.Exec(
		`INSERT INTO `+terminalTableName+` (`+
			`id, application_id, user_id, secret, `+
			`c_ts, c_uid, c_tid, `+
			`c_origin_address, c_origin_env, `+
			`display_name, accept_language, `+
			`verification_type, verification_id, verification_time `+
			`) VALUES (`+
			`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14`+
			`)`,
		termID.PrimitiveValue(),
		input.ApplicationRef.ID().PrimitiveValue(),
		input.UserRef.ID().PrimitiveValue(),
		termSecret,
		ctxTime,
		authCtx.UserIDPtr(),
		authCtx.TerminalIDPtr(),
		callCtx.RemoteAddress(),
		callCtx.RemoteEnvironmentString(),
		input.DisplayName,
		input.AcceptLanguage, //TODO: get from context
		input.VerificationType,
		input.VerificationID,
		input.VerificationTime)
	if err != nil {
		return iam.TerminalRefKeyZero(), "", err
	}

	terminalRef = iam.NewTerminalRefKey(input.ApplicationRef, input.UserRef, termID)
	if generateSecret {
		return terminalRef, termSecret, nil
	}
	return terminalRef, "", nil
}

func (core *Core) setUserTerminalVerified(
	callCtx iam.CallContext,
	terminalID iam.TerminalID,
	disallowReplay bool,
) (secret string, err error) {
	// A secret is used to obtain an access token. if an access token is
	// expired, the terminal could request another access token by
	// providing this secret. the secret is only provided after the
	// authorization has been verified.
	termSecret := core.generateTerminalSecret() //TODO(exa): JWT (or something similar)
	xres, err := core.db.
		Exec(
			`UPDATE `+terminalTableName+` `+
				`SET (secret, verification_time) = ($1, $2) `+
				`WHERE id = $3 AND verification_time IS NULL`,
			termSecret, callCtx.RequestReceiveTime(), terminalID)
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
		err = core.db.
			QueryRow(
				`SELECT secret FROM `+terminalTableName+` `+
					`WHERE id = $1`,
				terminalID).
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

func (core *Core) getTerminalAcceptLanguages(
	id iam.TerminalID,
) ([]language.Tag, error) {
	var acceptLanguage string
	err := core.db.QueryRow(
		`SELECT accept_language `+
			`FROM `+terminalTableName+` `+
			`WHERE id = $1`,
		id).Scan(&acceptLanguage)
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
