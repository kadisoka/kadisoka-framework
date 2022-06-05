package pnv10n

import (
	"github.com/kadisoka/kadisoka-framework/volib/pkg/telephony"
	"github.com/rez-go/stev"
)

func init() {
	RegisterModule(
		"null",
		Module{
			ConfigSkeleton: func() ModuleConfig {
				return &moduleNULLConfig{}
			},
			NewSMSDeliveryService: func(config interface{}) SMSDeliveryService {
				return &smsDeliveryServiceNULL{}
			},
		})
}

type moduleNULLConfig struct {
}

func (moduleNULLConfig) SMSDeliveryServiceConfig() interface{} {
	return &smsDeliveryServiceNULLConfig{}
}

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
