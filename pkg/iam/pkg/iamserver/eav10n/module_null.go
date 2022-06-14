package eav10n

import (
	"github.com/kadisoka/kadisoka-framework/pkg/volib/pkg/email"
	"github.com/rez-go/stev"
)

func init() {
	RegisterModule(
		"null",
		Module{
			ConfigSkeleton: func() ModuleConfig {
				return &moduleNULLConfig{}
			},
			NewEmailDeliveryService: func(config interface{}) EmailDeliveryService {
				return &emailDeliveryServiceNULL{}
			},
		})
}

type moduleNULLConfig struct {
}

func (moduleNULLConfig) EmailDeliveryServiceConfig() interface{} {
	return &emailDeliveryServiceNULLConfig{}
}

type emailDeliveryServiceNULL struct {
}

func (smsDS emailDeliveryServiceNULL) SendHTMLMessage(
	recipient email.Address,
	subjectText string,
	htmlContent string,
	opts EmailDeliveryOptions,
) error {
	return nil
}

type emailDeliveryServiceNULLConfig struct{}

func (emailDeliveryServiceNULLConfig) SelfDocsDescriptor() stev.SelfDocsDescriptor {
	return stev.SelfDocsDescriptor{
		ShortDesc: "Emails will not be delivered",
	}
}
