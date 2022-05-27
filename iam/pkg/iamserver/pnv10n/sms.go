package pnv10n

import (
	"github.com/rez-go/stev"

	"github.com/kadisoka/kadisoka-framework/volib/pkg/telephony"
)

type SMSDeliveryService interface {
	SendTextMessage(
		recipient telephony.PhoneNumber,
		text string,
		opts SMSDeliveryOptions) error
}

type SMSDeliveryOptions struct{}

type smsDeliveryServiceNULL struct {
}

func (smsDS smsDeliveryServiceNULL) SendTextMessage(
	recipient telephony.PhoneNumber,
	text string,
	opts SMSDeliveryOptions,
) error {
	return nil
}

type smsDeliveryServiceNULLConfig struct{}

func (smsDeliveryServiceNULLConfig) SelfDocsDescriptor() stev.SelfDocsDescriptor {
	return stev.SelfDocsDescriptor{
		ShortDesc: "SMSes will not be delivered",
	}
}

type moduleNULLConfig struct {
}

func (moduleNULLConfig) SMSDeliveryServiceConfig() interface{} {
	return &smsDeliveryServiceNULLConfig{}
}
