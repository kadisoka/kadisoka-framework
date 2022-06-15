package iamserver

import (
	"strings"
	"time"

	"github.com/alloyzeus/go-azfl/errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/lib/pq"

	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/pkg/volib/pkg/email"
	"github.com/kadisoka/kadisoka-framework/pkg/volib/pkg/telephony"
)

//TODO:SEC: harden
func (core *Core) AuthorizeTerminalByUserIdentifierAndPassword(
	inputCtx iam.CallInputContext,
	reqApp *iam.Application,
	terminalDisplayName string,
	identifier string,
	password string,
) (terminalID iam.TerminalID, terminalSecret string, userID iam.UserID, err error) {
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
			logCtx(inputCtx).Error().Err(err).
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
				logCtx(inputCtx).Error().Err(err).
					Msg("getUserIDNumByKeyPhoneNumberInsecure")
			} else {
				userIDNum = ownerUserIDNum
			}
		}
	}

	//TODO: last, check if it matches ourr specification of usernames

	if userIDNum.IsNotStaticallyValid() {
		// No errors
		return iam.TerminalIDZero(), "", iam.UserIDZero(), nil
	}

	userID = iam.NewUserID(userIDNum)

	passwordMatch, err := core.MatchUserPassword(userID, password)
	if err != nil {
		return iam.TerminalIDZero(), "", iam.UserIDZero(),
			errors.Wrap("matching user password", err)
	}

	if !passwordMatch {
		return iam.TerminalIDZero(), "", iam.UserIDZero(), nil
	}

	var appID iam.ApplicationID
	if reqApp != nil {
		appID = reqApp.ID
	}

	regOutCtx, regOutData := core.RegisterTerminal(inputCtx, TerminalRegistrationInputData{
		ApplicationID:    appID,
		UserID:           userID,
		DisplayName:      terminalDisplayName,
		VerificationType: iam.TerminalVerificationResourceTypeOAuthPassword,
		VerificationID:   0, //TODO: request ID or such
	})
	if err = regOutCtx.Err; err != nil {
		return iam.TerminalIDZero(), "", iam.UserIDZero(),
			errors.Wrap("RegisterTerminal", err)
	}

	return regOutData.TerminalID, regOutData.TerminalSecret, userID, nil
}

func (core *Core) issueSession(
	inputCtx iam.CallInputContext,
	terminalID iam.TerminalID,
	userID iam.UserID,
) (
	sessionID iam.SessionID,
	issueTime time.Time,
	expiry time.Time,
	err error,
) {
	ctxAuth := inputCtx.Authorization()

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
			return iam.SessionIDZero(), timeZero, timeZero, err
		}
		sqlString, _, _ := goqu.
			Insert(sessionDBTableName).
			Rows(
				goqu.Record{
					sessionDBColTerminalID:           terminalID.IDNum().PrimitiveValue(),
					sessionDBColIDNum:                sessionIDNum.PrimitiveValue(),
					"expiry":                         sessionExpiry,
					sessionDBColMDCreationTimestamp:  sessionStartTime,
					sessionDBColMDCreationTerminalID: ctxAuth.TerminalIDNumPtr(),
					sessionDBColMDCreationUserID:     ctxAuth.UserIDNumPtr(),
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
			pqErr.Constraint == sessionDBTablePrimaryKeyName {
			if attemptNum >= attemptNumMax {
				return iam.SessionIDZero(), timeZero, timeZero,
					errors.Wrap("insert max attempts", err)
			}
			continue
		}

		return iam.SessionIDZero(), timeZero, timeZero,
			errors.Wrap("insert", err)
	}

	return iam.NewSessionID(terminalID, sessionIDNum),
		sessionStartTime, sessionExpiry, nil
}
