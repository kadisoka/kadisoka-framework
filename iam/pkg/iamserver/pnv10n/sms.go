package pnv10n

type SMSDeliveryService interface {
	SendTextMessage(
		recipientPhoneNumber string, //TODO: iam.PhoneNumber or telephony.Number
		text string,
		opts SMSDeliveryOptions) error
}

type SMSDeliveryOptions struct{}

type smsDeliveryServiceNULL struct {
}

func (smsDS smsDeliveryServiceNULL) SendTextMessage(
	recipient string,
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
