package pnv10n

import "github.com/kadisoka/kadisoka-framework/volib/pkg/telephony"

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

func init() {
	type smsDeliveryServiceNULLConfig struct{}

	RegisterModule(
		"null",
		Module{
			ConfigSkeleton: func() interface{} {
				return &smsDeliveryServiceNULLConfig{}
			},
			NewSMSDeliveryService: func(config interface{}) SMSDeliveryService {
				return &smsDeliveryServiceNULL{}
			},
		})
}
