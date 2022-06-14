package twilio

import (
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iamserver/pnv10n"
)

const ServiceName = "twilio"

func init() {
	pnv10n.RegisterModule(
		ServiceName,
		pnv10n.Module{
			ConfigSkeleton: func() pnv10n.ModuleConfig {
				cfg := ModuleConfigSkeleton()
				return &cfg
			},
			NewSMSDeliveryService: NewSMSDeliveryService,
		})
}

func ModuleConfigSkeleton() ModuleConfig {
	smsCfg := SMSDeliveryServiceConfigSkeleton()
	return ModuleConfig{
		SMS: &smsCfg,
	}
}

type ModuleConfig struct {
	SMS *SMSDeliveryServiceConfig `env:"SMS"`
}

func (moduleCfg ModuleConfig) SMSDeliveryServiceConfig() interface{} {
	return moduleCfg.SMS
}
