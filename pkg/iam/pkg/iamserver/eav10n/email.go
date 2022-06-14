package eav10n

import (
	"github.com/kadisoka/kadisoka-framework/pkg/volib/pkg/email"
)

type EmailDeliveryService interface {
	SendHTMLMessage(
		recipient email.Address,
		subjectText string,
		htmlContent string,
		opts EmailDeliveryOptions) error
}

type EmailDeliveryOptions struct {
	MessageCharset string
	SenderAddress  string
}
