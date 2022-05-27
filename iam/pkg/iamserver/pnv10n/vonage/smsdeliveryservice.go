package vonage

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/pnv10n"
	"github.com/kadisoka/kadisoka-framework/volib/pkg/telephony"
)

type SMSDeliveryServiceConfig struct {
	APIKey    string `env:"API_KEY,required"`
	APISecret string `env:"API_SECRET,required"`
	Sender    string `env:"SENDER,required"`
}

func SMSDeliveryServiceConfigSkeleton() SMSDeliveryServiceConfig { return SMSDeliveryServiceConfig{} }

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

	if conf.APIKey == "" {
		panic("NEXMO API Key not provided")
	}
	if conf.APISecret == "" {
		panic("NEXMO API Secret not provided")
	}

	return &SMSDeliveryService{
		config:      conf,
		endpointURL: "https://rest.nexmo.com/sms/json",
	}
}

func (sms *SMSDeliveryService) SendTextMessage(
	recipient telephony.PhoneNumber,
	text string,
	opts pnv10n.SMSDeliveryOptions,
) error {
	sender := sms.config.Sender
	if sender == "" {
		sender = "Nexmo"
	}
	endPoint := sms.endpointURL
	bodyReq := url.Values{}
	bodyReq.Set("to", strings.Trim(recipient.String(), "+"))
	bodyReq.Set("text", text)
	bodyReq.Set("type", "text")
	bodyReq.Set("from", sender)
	payload := strings.NewReader(bodyReq.Encode())

	req, err := http.NewRequest("POST", endPoint, payload)
	if err != nil {
		return errors.New("Unable to build new request > " + err.Error())
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)

	if err != nil {
		return errors.New("Unable to send request -> " + err.Error())
	}
	defer resp.Body.Close()

	// resp.StatusCode is between 200 and 300.
	// This is because an HTTP status code with the form 2XX signifies a successful HTTP POST request
	// https://standard.telesign.com/api-reference/apis/sms-api/send-an-sms/reference
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	errBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New(err.Error())
	}
	return errors.New(string(errBody))
}
