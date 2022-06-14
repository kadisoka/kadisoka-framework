package telesign

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iamserver/pnv10n"
	"github.com/kadisoka/kadisoka-framework/pkg/volib/pkg/telephony"
)

type SMSDeliveryServiceConfig struct {
	CustomerID string `env:"CUSTOMER_ID,required"`
	APIKey     string `env:"API_KEY,required"`
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

	if len(conf.APIKey) <= 0 {
		panic("Telesign API Key not provided")
	}
	if len(conf.CustomerID) <= 0 {
		panic("Telesign Customer ID not provided")
	}

	return &SMSDeliveryService{
		config:      conf,
		endpointURL: "https://rest-api.telesign.com/v1/messaging",
	}
}

// SendTextMessage is use for send text message using sms delivery service
func (sms *SMSDeliveryService) SendTextMessage(
	recipient telephony.PhoneNumber,
	text string,
	opts pnv10n.SMSDeliveryOptions,
) error {
	endPoint := sms.endpointURL
	bodyReq := url.Values{}
	bodyReq.Set("phone_number", strings.Trim(recipient.String(), "+"))
	bodyReq.Set("message", text)
	bodyReq.Set("message_type", "ARN")
	payload := strings.NewReader(bodyReq.Encode())

	req, err := http.NewRequest("POST", endPoint, payload)

	if err != nil {
		return errors.New("Unable to build new request -> " + err.Error())
	}
	req.SetBasicAuth(sms.config.CustomerID, sms.config.APIKey)
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

	var errorData telesignResponse
	errBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(errBody, &errorData)
	if err != nil {
		return err
	}

	switch errorData.Status.Code {
	case 11000, 11001, 10033:
		return pnv10n.InvalidPhoneNumberError{Err: errors.New(string(errBody))}
	}

	return pnv10n.ConfigurationError{Err: errors.New(string(errBody))}
}

type telesignResponse struct {
	ReferenceID string                 `json:"reference_id"`
	ExternalID  string                 `json:"external_id"`
	Status      telesignStatusResponse `json:"status"`
}

type telesignStatusResponse struct {
	Code        int64  `json:"code"`
	Description string `json:"description"`
	UpdatedOn   string `json:"updated_on"`
}
