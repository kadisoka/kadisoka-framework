package ses

import "github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/eav10n"

const ServiceName = "ses"

func init() {
	eav10n.RegisterModule(
		ServiceName,
		eav10n.Module{
			ConfigSkeleton: func() eav10n.ModuleConfig {
				cfg := ModuleConfigSkeleton()
				return &cfg
			},
			NewEmailDeliveryService: NewEmailDeliveryService,
		})
}

func ModuleConfigSkeleton() ModuleConfig {
	emailDeliveryCfg := EmailDeliveryServiceConfigSkeleton()
	return ModuleConfig{
		Email: &emailDeliveryCfg,
	}
}

type ModuleConfig struct {
	Email *EmailDeliveryServiceConfig `env:",squash"`
}

func (moduleCfg ModuleConfig) EmailDeliveryServiceConfig() interface{} {
	return moduleCfg.Email
}
