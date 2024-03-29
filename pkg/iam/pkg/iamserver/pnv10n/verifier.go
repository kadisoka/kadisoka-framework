//
package pnv10n

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/alloyzeus/go-azfl/errors"
	"github.com/jmoiron/sqlx"
	"golang.org/x/text/language"

	"github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/realm"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/pkg/volib/pkg/telephony"
)

const verificationDBTableName = "phone_number_verification_dt"

func NewVerifier(
	realmInfo realm.Info,
	db *sqlx.DB,
	config Config,
) *Verifier {
	if realmInfo.Name == "" {
		panic("Invalid config")
	}

	loadTemplates()

	smsDeliveryServices := map[int32]SMSDeliveryService{}
	serviceName := config.SMSDeliveryService
	if serviceName == "" || !strings.Contains(serviceName, ",") {
		if serviceName == "" {
			panic("SMS delivery service must be specified in the config")
		}
		moduleCfg := config.Modules[serviceName]
		if moduleCfg != nil {
			deliverySvc, err := NewSMSDeliveryService(serviceName, moduleCfg.SMSDeliveryServiceConfig())
			if err != nil || deliverySvc == nil {
				panic("SMS delivery service not configured")
			}
			smsDeliveryServices[0] = deliverySvc
		}
	} else {
		codeNames := strings.Split(serviceName, ",")
		instantiatedServices := map[string]SMSDeliveryService{}
		for _, codeName := range codeNames {
			parts := strings.Split(codeName, ":")
			svcName := strings.TrimSpace(parts[1])
			svcInst := instantiatedServices[svcName]
			if svcInst == nil {
				moduleCfg := config.Modules[serviceName]
				if moduleCfg != nil {
					deliverySvc, err := NewSMSDeliveryService(serviceName, moduleCfg.SMSDeliveryServiceConfig())
					if err != nil || deliverySvc == nil {
						panic("SMS delivery service not configured")
					}
					instantiatedServices[svcName] = deliverySvc
					svcInst = deliverySvc
				}
			}
			ccStr := parts[0]
			if ccStr == "*" {
				smsDeliveryServices[0] = svcInst
			} else {
				countryCode, err := strconv.ParseInt(ccStr, 10, 16)
				if err != nil {
					panic(err)
				}
				if countryCode <= 0 {
					panic("Invalid country code " + ccStr)
				}
				if _, dup := smsDeliveryServices[int32(countryCode)]; dup {
					panic("Duplicated country code " + ccStr)
				}
				smsDeliveryServices[int32(countryCode)] = svcInst
			}
		}
	}
	if _, ok := smsDeliveryServices[0]; !ok {
		panic("Requires at least one SMS delivery service")
	}

	var codeTTLDefault time.Duration
	if config.CodeTTLDefault > 0 {
		codeTTLDefault = config.CodeTTLDefault
	} else {
		codeTTLDefault = 5 * time.Minute //TODO: should be based on the length of the code
	}

	confirmationAttemptsMax := config.ConfirmationAttemptsMax
	if confirmationAttemptsMax == 0 {
		confirmationAttemptsMax = 5
	}

	return &Verifier{
		realmInfo:                    realmInfo,
		db:                           db,
		config:                       config,
		codeTTLDefaultValue:          codeTTLDefault,
		smsDeliveryServicesByCountry: smsDeliveryServices,
		confirmationAttemptsMax:      confirmationAttemptsMax,
	}
}

type Verifier struct {
	realmInfo                    realm.Info
	db                           *sqlx.DB
	codeTTLDefaultValue          time.Duration
	confirmationAttemptsMax      int16
	config                       Config
	smsDeliveryServicesByCountry map[int32]SMSDeliveryService
}

//TODO(exa): make the operations atomic
func (verifier *Verifier) StartVerification(
	inputCtx iam.CallInputContext,
	phoneNumber telephony.PhoneNumber,
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
			"SELECT id_num, code_expiry, confirmation_attempts_remaining "+
				`FROM `+verificationDBTableName+` `+
				"WHERE country_code = $1 AND national_number = $2 AND confirmation_ts IS NULL "+
				"ORDER BY id_num DESC "+
				"LIMIT 1",
			phoneNumber.CountryCode(),
			phoneNumber.NationalNumber()).
		Scan(&prevVerificationID, &prevCodeExpiry, &prevAttempts)
	if err == nil {
		// Return previous verification code
		if prevAttempts > 0 && prevCodeExpiry.After(ctxTime.Add(-10*time.Second)) {
			return prevVerificationID, &prevCodeExpiry, nil
		}
	}

	if int64(codeTTL) <= 0 {
		codeTTL = verifier.codeTTLDefaultValue
	}

	code := verifier.generateVerificationCode()
	// Truncate because sub-ms value might be problematic
	// for some parsers. To minute because it's more humane.
	codeExp := ctxTime.Add(codeTTL).Truncate(time.Minute)
	codeExpiry = &codeExp

	err = verifier.db.
		QueryRow(
			`INSERT INTO `+verificationDBTableName+` (`+
				"country_code, national_number, "+
				"md_c_ts, md_c_uid, md_c_tid, "+
				"code, code_expiry, confirmation_attempts_remaining"+
				") VALUES ($1, $2, $3, $4, $5, $6, $7, $8) "+
				"RETURNING id_num",
			phoneNumber.CountryCode(),
			phoneNumber.NationalNumber(),
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
	err = verifier.sendTextMessage(
		phoneNumber, code, userPreferredLanguages, noDelivery)
	if err != nil {
		return 0, nil, err
	}

	return
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
			`SET confirmation_attempts_remaining = confirmation_attempts_remaining - 1 `+
			`WHERE id_num = $1 `+
			`RETURNING *`,
		verificationID).
		StructScan(&dbData)
	if err != nil {
		return err
	}

	if dbData.ConfirmationAttemptsRemaining < 0 {
		return ErrVerificationCodeExpired
	}
	if dbData.Code != code {
		return ErrVerificationCodeMismatch
	}
	if dbData.CodeExpiry != nil && dbData.CodeExpiry.Before(ctxTime) {
		// Delete?
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

func (verifier *Verifier) GetPhoneNumberByVerificationID(
	verificationID int64,
) (*telephony.PhoneNumber, error) {
	var countryCode int32
	var nationalNumber int64
	err := verifier.db.QueryRow(
		"SELECT country_code, national_number "+
			`FROM `+verificationDBTableName+` `+
			"WHERE id_num = $1 LIMIT 1",
		verificationID).
		Scan(&countryCode, &nationalNumber)
	if err != nil {
		return nil, err
	}
	result := telephony.NewPhoneNumber(countryCode, nationalNumber)
	return &result, nil
}

func (verifier *Verifier) GetVerificationCodeByPhoneNumber(
	phoneNumber telephony.PhoneNumber,
) (code string, err error) {
	err = verifier.db.QueryRow(
		"SELECT code "+
			`FROM `+verificationDBTableName+` `+
			"WHERE country_code = $1 AND national_number = $2 "+
			"AND confirmation_ts IS NULL "+
			"ORDER BY md_c_ts DESC LIMIT 1",
		phoneNumber.CountryCode(), phoneNumber.NationalNumber()).
		Scan(&code)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return code, nil
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

func (verifier *Verifier) sendTextMessage(
	phoneNumber telephony.PhoneNumber,
	code string,
	userPreferredLanguages []language.Tag,
	noDelivery bool,
) error {
	var messageTemplate *template.Template
	if len(userPreferredLanguages) != 0 {
		for _, locale := range userPreferredLanguages {
			messageTemplate = localizedMessageTemplates[locale.String()]
			if messageTemplate != nil {
				break
			}
		}
	}
	if messageTemplate == nil {
		messageTemplate = localizedMessageTemplates[messageLocaleDefault.String()]
	}

	var bodyBuilder strings.Builder
	err := messageTemplate.
		Execute(&bodyBuilder, map[string]interface{}{
			"RealmName": verifier.realmInfo.Name,
			"Code":      code,
		})
	if err != nil {
		return err
	}

	// https://developers.google.com/identity/sms-retriever/verify
	bodyString := "<#> " + bodyBuilder.String()
	if verifier.config.SMSRetrieverAppHash != "" {
		bodyString += "\n" + verifier.config.SMSRetrieverAppHash
	}

	if !noDelivery {
		//NOTE: special treatment for +1555xxxx numbers (for testing)
		if !(phoneNumber.CountryCode() == 1 && phoneNumber.NationalNumber() > 5550000 && phoneNumber.NationalNumber() <= 5559999) {
			var deliverySvc SMSDeliveryService
			if len(verifier.smsDeliveryServicesByCountry) > 1 {
				if svc := verifier.smsDeliveryServicesByCountry[phoneNumber.CountryCode()]; svc != nil {
					deliverySvc = svc
				}
			}
			if deliverySvc == nil {
				deliverySvc = verifier.smsDeliveryServicesByCountry[0]
			}
			err = deliverySvc.SendTextMessage(
				phoneNumber,
				bodyString,
				SMSDeliveryOptions{})
		}
	}

	verifier.notifyMaintainersChannels(phoneNumber, bodyBuilder.String(), err)

	return err
}

func (verifier *Verifier) notifyMaintainersChannels(
	phoneNumber telephony.PhoneNumber, messageBody string, sendError error,
) {
	textMessage := fmt.Sprintf(
		"Phone number verification for `%s`\n```\n%s```",
		phoneNumber.String(), messageBody)
	if sendError != nil {
		textMessage += "\nWith error: `" + sendError.Error() + "`"
	}
}

var localizedMessageTemplates map[string]*template.Template

//TODO: load these from somewhere (e.g., Firebase remote config)
var localizedMessageTemplateSources = map[string][]string{
	"{{ .RealmName }} - verification code: {{ .Code }}":    {"en", "en-US", "en-GB"},
	"{{ .RealmName }} - kode verifikasi Anda: {{ .Code }}": {"id", "id-ID"},
}

func loadTemplates() {
	localizedMessageTemplates = make(map[string]*template.Template)
	// Load all message templates. We also ensure that there's no
	// duplicates for the same language.
	for tplstr, locales := range localizedMessageTemplateSources {
		if len(locales) == 0 {
			continue
		}
		tpl := template.Must(template.New("verification-message").Parse(tplstr))
		for _, locale := range locales {
			if locale == "" {
				continue
			}
			langTag := language.MustParse(locale)
			if _, ok := localizedMessageTemplates[langTag.String()]; ok {
				panic("duplicate for locale " + locale + " (" + langTag.String() + ")")
			}
			localizedMessageTemplates[langTag.String()] = tpl
		}
	}
	// Ensure that we have a message template for the default locale.
	if v := localizedMessageTemplates[messageLocaleDefault.String()]; v == nil {
		panic("no template for default locale " + messageLocaleDefault.String())
	}
}
