package iamserver

import (
	"strings"
	"time"

	"github.com/alloyzeus/go-azfl/errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/lib/pq"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/volib/pkg/email"
	"github.com/kadisoka/kadisoka-framework/volib/pkg/telephony"
)

//TODO:SEC: harden
func (core *Core) AuthorizeTerminalByUserIdentifierAndPassword(
	callCtx iam.CallInputContext,
	reqApp *iam.Application,
	terminalDisplayName string,
	identifier string,
	password string,
) (terminalRef iam.TerminalRefKey, terminalSecret string, userRef iam.UserRefKey, err error) {
	//TODO: check context

	identifier = strings.TrimSpace(identifier)

	// Username with scheme. The format is '<scheme>:<scheme-specific-identifier>'
	if names := strings.SplitN(identifier, ":", 2); len(names) == 2 {
		switch names[0] {
		case "terminal":
			panic("TODO")
		default:
		}
	}

	var userIDNum iam.UserIDNum

	//TODO: create a method `isAuthenticationByEmailAddressAllowed`
	if emailAddress, err := email.AddressFromString(identifier); err == nil {
		//TODO: by email
		ownerUserIDNum, err := core.getUserIDNumByKeyEmailAddressInsecure(emailAddress)
		if err != nil {
			logCtx(callCtx).Error().Err(err).
				Msg("getUserIDNumByKeyEmailAddressInsecure")
		} else {
			userIDNum = ownerUserIDNum
		}
	}

	if userIDNum.IsNotStaticallyValid() {
		if phoneNumber, err := telephony.PhoneNumberFromString(identifier); err == nil {
			//TODO: by phone number
			ownerUserIDNum, err := core.getUserIDNumByKeyPhoneNumberInsecure(phoneNumber)
			if err != nil {
				logCtx(callCtx).Error().Err(err).
					Msg("getUserIDNumByKeyPhoneNumberInsecure")
			} else {
				userIDNum = ownerUserIDNum
			}
		}
	}

	//TODO: last, check if it matches ourr specification of usernames

	if userIDNum.IsNotStaticallyValid() {
		// No errors
		return iam.TerminalRefKeyZero(), "", iam.UserRefKeyZero(), nil
	}

	userRef = iam.NewUserRefKey(userIDNum)

	passwordMatch, err := core.MatchUserPassword(userRef, password)
	if err != nil {
		return iam.TerminalRefKeyZero(), "", iam.UserRefKeyZero(),
			errors.Wrap("matching user password", err)
	}

	if !passwordMatch {
		return iam.TerminalRefKeyZero(), "", iam.UserRefKeyZero(), nil
	}

	var appRef iam.ApplicationRefKey
	if reqApp != nil {
		appRef = reqApp.RefKey
	}
	regOutput := core.RegisterTerminal(TerminalRegistrationInput{
		Context:        callCtx,
		ApplicationRef: appRef,
		Data: TerminalRegistrationInputData{
			UserRef:          userRef,
			DisplayName:      terminalDisplayName,
			VerificationType: iam.TerminalVerificationResourceTypeOAuthPassword,
			VerificationID:   0, //TODO: request ID or such
		}})
	if err = regOutput.Context.Err; err != nil {
		return iam.TerminalRefKeyZero(), "", iam.UserRefKeyZero(),
			errors.Wrap("RegisterTerminal", err)
	}

	return regOutput.Data.TerminalRef, regOutput.Data.TerminalSecret, userRef, nil
}

func (core *Core) issueSession(
	callCtx iam.CallInputContext,
	terminalRef iam.TerminalRefKey,
	userRef iam.UserRefKey,
) (
	sessionRef iam.SessionRefKey,
	issueTime time.Time,
	expiry time.Time,
	err error,
) {
	ctxAuth := callCtx.Authorization()

	const attemptNumMax = 5

	timeZero := time.Time{}
	sessionStartTime := timeZero
	sessionExpiry := timeZero
	var sessionIDNum iam.SessionIDNum

	for attemptNum := 0; ; attemptNum++ {
		sessionStartTime = time.Now().UTC()
		sessionExpiry = sessionStartTime.Add(iam.AccessTokenTTLDefault)
		sessionIDNum, err = GenerateSessionIDNum(0)
		if err != nil {
			return iam.SessionRefKeyZero(), timeZero, timeZero, err
		}
		sqlString, _, _ := goqu.
			Insert(sessionDBTableName).
			Rows(
				goqu.Record{
					"terminal_id": terminalRef.IDNum().PrimitiveValue(),
					"id_num":      sessionIDNum.PrimitiveValue(),
					"expiry":      sessionExpiry,
					"_mc_ts":      sessionStartTime,
					"_mc_tid":     ctxAuth.TerminalIDNumPtr(),
					"_mc_uid":     ctxAuth.UserIDNumPtr(),
				},
			).
			ToSQL()
		_, err = core.db.
			Exec(sqlString)
		if err == nil {
			break
		}

		pqErr, _ := err.(*pq.Error)
		if pqErr != nil &&
			pqErr.Code == "23505" &&
			pqErr.Constraint == sessionDBTableName+"_pkey" {
			if attemptNum >= attemptNumMax {
				return iam.SessionRefKeyZero(), timeZero, timeZero,
					errors.Wrap("insert max attempts", err)
			}
			continue
		}

		return iam.SessionRefKeyZero(), timeZero, timeZero,
			errors.Wrap("insert", err)
	}

	return iam.NewSessionRefKey(terminalRef, sessionIDNum),
		sessionStartTime, sessionExpiry, nil
}
