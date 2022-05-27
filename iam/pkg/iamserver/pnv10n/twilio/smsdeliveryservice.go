package twilio

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/rez-go/stev"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/pnv10n"
	"github.com/kadisoka/kadisoka-framework/volib/pkg/telephony"
)

type SMSDeliveryServiceConfig struct {
	AccountSID string `env:"ACCOUNT_SID,required"`
	AuthToken  string `env:"AUTH_TOKEN,required"`
	Sender     string `env:"SENDER,required"`
}

func SMSDeliveryServiceConfigSkeleton() SMSDeliveryServiceConfig { return SMSDeliveryServiceConfig{} }

func (SMSDeliveryServiceConfig) SelfDocsDescriptor() stev.SelfDocsDescriptor {
	return stev.SelfDocsDescriptor{
		ShortDesc: "Use Twilio to deliver the SMS",
	}
}

type SMSDeliveryService struct {
	config      *SMSDeliveryServiceConfig
	endpointURL string
}

var _ pnv10n.SMSDeliveryService = &SMSDeliveryService{}

func NewSMSDeliveryService(config interface{}) pnv10n.SMSDeliveryService {
	if config == nil {
		panic(errors.New("configuration required"))
	}
	conf, ok := config.(*SMSDeliveryServiceConfig)
	if !ok {
		panic(errors.New("configuration of invalid type"))
	}

	if conf.AccountSID == "" {
		panic("Twilio Account SID not found")
	}
	if conf.AuthToken == "" {
		panic("Unable to find twilio auth token")
	}

	return &SMSDeliveryService{
		config:      conf,
		endpointURL: "https://api.twilio.com/2010-04-01/Accounts/" + conf.AccountSID + "/Messages.json",
	}
}

// SendTextMessage is use for send text message using sms delivery service
func (sms *SMSDeliveryService) SendTextMessage(
	recipient telephony.PhoneNumber,
	text string,
	opts pnv10n.SMSDeliveryOptions,
) error {
	// Docs: https://www.twilio.com/docs/sms/send-messages

	endPoint := sms.endpointURL

	bodyReq := url.Values{}
	bodyReq.Set("To", recipient.String())
	bodyReq.Set("From", sms.config.Sender)
	bodyReq.Set("Body", text)
	payload := strings.NewReader(bodyReq.Encode())

	req, err := http.NewRequest("POST", endPoint, payload)

	if err != nil {
		return errors.New("Unable to build new request -> " + err.Error())
	}

	req.SetBasicAuth(sms.config.AccountSID, sms.config.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	// resp.StatusCode is between 200 and 300.
	// This is because an HTTP status code with the form 2XX signifies a successful HTTP POST request
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	errBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var errorData twilioErrorResponse
	err = json.Unmarshal(errBody, &errorData)
	if err != nil {
		return err
	}

	switch errorData.Code {
	case 20003:
		return pnv10n.ConfigurationError{Err: errors.New(string(errBody))}
	case 20404:
		return pnv10n.GatewayError{Err: errors.New(string(errBody))}
	case 21211: // Invalid phone number
		return pnv10n.InvalidPhoneNumberError{Err: errors.New(string(errBody))}
	case 21408: // Permission to send an SMS has not been enabled for the region
		return pnv10n.PhoneNumberRegionNotSupportedError{Err: errors.New(string(errBody))}
	case 21614: // Not a mobile phone number
		return pnv10n.InvalidPhoneNumberError{Err: errors.New(string(errBody))}
	case 30008: // Unknown error
	default:
		return pnv10n.GatewayError{Err: errors.New(string(errBody))}
	}

	return nil
}

type twilioErrorResponse struct {
	Code     int64  `json:"code"`
	Message  string `json:"message"`
	MoreInfo string `json:"more_info"`
	Status   int64  `json:"status"`
}
