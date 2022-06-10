// Package eav10n provides utilities for verifying email addresses.
package eav10n

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	htmltpl "html/template"
	"strings"
	texttpl "text/template"
	"time"

	"github.com/alloyzeus/go-azfl/errors"
	"github.com/jmoiron/sqlx"
	"golang.org/x/text/language"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/realm"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/volib/pkg/email"
)

const verificationDBTableName = "email_address_verification_dt"

func NewVerifier(
	realmInfo realm.Info,
	db *sqlx.DB,
	config Config,
) *Verifier {
	emailSenderAddress := realmInfo.ServiceNotificationEmailsSenderAddress
	if emailSenderAddress == "" {
		emailSenderAddress = realmInfo.Contact.EmailAddress
	}

	if config.SenderAddress == "" && emailSenderAddress == "" {
		panic("sender address not configured")
	}

	resDir := config.ResourcesDir
	if resDir == "" {
		resDir = ResourcesDirDefault
	}
	loadTemplates(resDir)

	if config.SenderAddress != "" {
		emailSenderAddress = config.SenderAddress
	}

	emailDeliveryServices := map[string]EmailDeliveryService{}
	serviceName := config.EmailDeliveryService
	if serviceName == "" || !strings.Contains(serviceName, ",") {
		if serviceName == "" {
			panic("Email delivery service must be specified in the config")
		}
		moduleCfg := config.Modules[serviceName]
		if moduleCfg != nil {
			deliverySvc, err := NewEmailDeliveryService(serviceName, moduleCfg.EmailDeliveryServiceConfig())
			if err != nil || deliverySvc == nil {
				panic("Email delivery service not configured")
			}
			emailDeliveryServices[""] = deliverySvc
		}
	} else {
		codeNames := strings.Split(serviceName, ",")
		instantiatedServices := map[string]EmailDeliveryService{}
		for _, codeName := range codeNames {
			parts := strings.Split(codeName, ":")
			svcName := strings.TrimSpace(parts[1])
			svcInst := instantiatedServices[svcName]
			if svcInst == nil {
				moduleCfg := config.Modules[serviceName]
				if moduleCfg != nil {
					deliverySvc, err := NewEmailDeliveryService(serviceName, moduleCfg.EmailDeliveryServiceConfig())
					if err != nil || deliverySvc == nil {
						panic("Email delivery service not configured")
					}
					instantiatedServices[svcName] = deliverySvc
					svcInst = deliverySvc
				}
			}
			ccStr := parts[0]
			if ccStr == "*" {
				emailDeliveryServices[""] = svcInst
			} else {
				targetDomain := ccStr
				if targetDomain == "" {
					panic("Invalid domain " + ccStr)
				}
				if _, dup := emailDeliveryServices[targetDomain]; dup {
					panic("Duplicated domain " + ccStr)
				}
				emailDeliveryServices[targetDomain] = svcInst
			}
		}
	}
	if _, ok := emailDeliveryServices[""]; !ok {
		panic("Requires at least one email delivery service")
	}

	var codeTTLDefault time.Duration
	if config.CodeTTLDefault > 0 {
		codeTTLDefault = config.CodeTTLDefault
	} else {
		codeTTLDefault = 15 * time.Minute //TODO: should be based on the length of the code
	}

	confirmationAttemptsMax := config.ConfirmationAttemptsMax
	if confirmationAttemptsMax == 0 {
		confirmationAttemptsMax = 5
	}

	return &Verifier{
		realmInfo:                     realmInfo,
		db:                            db,
		senderAddress:                 emailSenderAddress,
		codeTTLDefault:                codeTTLDefault,
		emailDeliveryServicesByDomain: emailDeliveryServices,
		confirmationAttemptsMax:       confirmationAttemptsMax,
	}
}

type Verifier struct {
	realmInfo                     realm.Info
	db                            *sqlx.DB
	senderAddress                 string
	codeTTLDefault                time.Duration
	confirmationAttemptsMax       int16
	emailDeliveryServicesByDomain map[string]EmailDeliveryService
}

//TODO(exa): make the operations atomic
func (verifier *Verifier) StartVerification(
	inputCtx iam.CallInputContext,
	emailAddress email.Address,
	codeTTL time.Duration,
	userPreferredLanguages []language.Tag,
	preferredVerificationMethods []VerificationMethod,
) (idNum int64, codeExpiry *time.Time, err error) {
	if inputCtx == nil {
		return 0, nil, errors.ArgMsg("inputCtx", "missing")
	}

	ctxAuth := inputCtx.Authorization()
	ctxTime := inputCtx.CallInputMetadata().ReceiveTime

	var prevAttempts int16
	var prevVerificationID int64
	var prevCodeExpiry time.Time
	err = verifier.db.
		QueryRow(
			"SELECT id_num, code_expiry, attempts_remaining "+
				`FROM `+verificationDBTableName+` `+
				"WHERE domain_part = $1 AND local_part = $2 AND confirmation_ts IS NULL "+
				"ORDER BY id_num DESC "+
				"LIMIT 1",
			emailAddress.DomainPart(),
			emailAddress.LocalPart()).
		Scan(&prevVerificationID, &prevCodeExpiry, &prevAttempts)
	if err == nil {
		// Return previous verification code
		if prevAttempts > 0 && prevCodeExpiry.After(ctxTime.Add(-10*time.Second)) {
			return prevVerificationID, &prevCodeExpiry, nil
		}
	}

	if codeTTL <= 0 {
		codeTTL = verifier.codeTTLDefault
	}

	code := verifier.generateVerificationCode()
	// Truncate because sub-ms value might be problematic
	// for some parsers. To minute because it's more humane.
	codeExp := ctxTime.Add(codeTTL).Truncate(time.Minute)
	codeExpiry = &codeExp

	err = verifier.db.
		QueryRow(
			`INSERT INTO `+verificationDBTableName+` (`+
				`domain_part, local_part, `+
				"_mc_ts, _mc_uid, _mc_tid, "+
				"code, code_expiry, attempts_remaining"+
				") VALUES ($1, $2, $3, $4, $5, $6, $7, $8) "+
				"RETURNING id_num",
			emailAddress.DomainPart(),
			emailAddress.LocalPart(),
			ctxTime,
			ctxAuth.UserIDNumPtr(),
			ctxAuth.TerminalIDNumPtr(),
			code,
			codeExp,
			verifier.confirmationAttemptsMax,
		).Scan(&idNum)
	if err != nil {
		return 0, nil, err
	}

	noDelivery := len(preferredVerificationMethods) == 1 &&
		preferredVerificationMethods[0] == VerificationMethodNone
	err = verifier.sendVerificationEmail(
		emailAddress, code, userPreferredLanguages, noDelivery)
	if err != nil {
		return 0, nil, err
	}

	return
}

func (verifier *Verifier) generateVerificationCode() string {
	b := make([]byte, 8)
	_, err := rand.Read(b[5:])
	if err != nil {
		panic(err)
	}
	h := binary.BigEndian.Uint64(b)
	return fmt.Sprintf("%06d", h%1000000)
}

func (verifier *Verifier) ConfirmVerification(
	inputCtx iam.CallInputContext,
	verificationID int64, code string,
) error {
	if inputCtx == nil {
		return errors.ArgMsg("inputCtx", "missing")
	}
	ctxAuth := inputCtx.Authorization()

	ctxTime := inputCtx.CallInputMetadata().ReceiveTime
	var dbData verificationDBModel

	err := verifier.db.QueryRowx(
		`UPDATE `+verificationDBTableName+` `+
			`SET attempts_remaining = attempts_remaining - 1 `+
			`WHERE id_num = $1 `+
			`RETURNING *`,
		verificationID).
		StructScan(&dbData)
	if err != nil {
		return err
	}

	if dbData.AttemptsRemaining < 0 {
		return ErrVerificationCodeExpired
	}
	if dbData.Code != code {
		return ErrVerificationCodeMismatch
	}
	if dbData.CodeExpiry != nil && dbData.CodeExpiry.Before(ctxTime) {
		return ErrVerificationCodeExpired
	}

	if dbData.ConfirmationTime != nil {
		return nil
	}

	_, err = verifier.db.Exec(
		`UPDATE `+verificationDBTableName+` `+
			"SET confirmation_ts = $1, confirmation_uid = $2, confirmation_tid = $3 "+
			"WHERE id_num = $4 AND confirmation_ts IS NULL",
		ctxTime, ctxAuth.UserIDNumPtr(), ctxAuth.TerminalIDNumPtr(), verificationID)
	return err //TODO: determine if it's race-condition
}

func (verifier *Verifier) sendVerificationEmail(
	emailAddress email.Address,
	code string,
	userPreferredLanguages []language.Tag,
	noDelivery bool,
) error {
	var subjectTemplate *texttpl.Template
	var bodyTemplate *htmltpl.Template
	if len(userPreferredLanguages) != 0 {
		for _, locale := range userPreferredLanguages {
			bodyTemplate = localizedAccountActivationBodyHTMLTemplates[locale.String()]
			subjectTemplate = localizedAccountActivationSubjectTemplates[locale.String()]
			if bodyTemplate != nil {
				break
			}
		}
	}
	if bodyTemplate == nil {
		bodyTemplate = localizedAccountActivationBodyHTMLTemplates[messageLocaleDefault.String()]
	}
	if subjectTemplate == nil {
		subjectTemplate = localizedAccountActivationSubjectTemplates[messageLocaleDefault.String()]
	}

	var err error

	buf := new(bytes.Buffer)
	err = subjectTemplate.
		Execute(buf, map[string]interface{}{
			"RealmName": verifier.realmInfo.Name,
		})
	if err != nil {
		return err
	}
	subjectText := buf.String()

	buf = new(bytes.Buffer)
	if err = bodyTemplate.Execute(buf, map[string]interface{}{
		"RealmInfo": verifier.realmInfo,
		"Title":     subjectText, //TODO: title == subject?
		"Code":      code,
	}); err != nil {
		panic(err)
	}
	htmlBody := buf.String()

	if !noDelivery {
		targetDomain := emailAddress.DomainPart()
		doSend := true
		switch targetDomain {
		case "example.com", "example.org", "example.net":
			doSend = false
		}

		if doSend {
			var deliverySvc EmailDeliveryService
			if len(verifier.emailDeliveryServicesByDomain) > 1 {
				if svc := verifier.emailDeliveryServicesByDomain[targetDomain]; svc != nil {
					deliverySvc = svc
				}
			}
			if deliverySvc == nil {
				deliverySvc = verifier.emailDeliveryServicesByDomain[""]
			}
			err = deliverySvc.SendHTMLMessage(
				emailAddress,
				subjectText,
				htmlBody,
				EmailDeliveryOptions{
					MessageCharset: messageCharset,
					SenderAddress:  verifier.senderAddress,
				})
		}
	}

	return err
}

func (verifier *Verifier) GetEmailAddressByVerificationID(
	verificationID int64,
) (*email.Address, error) {
	var localPart, domainPart string
	err := verifier.db.QueryRow(
		`SELECT domain_part, local_part `+
			`FROM `+verificationDBTableName+` `+
			`WHERE id_num = $1 `,
		verificationID).
		Scan(&domainPart, &localPart)
	if err != nil {
		return nil, err
	}
	emailAddress, err := email.AddressFromString(localPart + "@" + domainPart)
	if err != nil {
		panic(err)
	}
	return &emailAddress, nil
}
