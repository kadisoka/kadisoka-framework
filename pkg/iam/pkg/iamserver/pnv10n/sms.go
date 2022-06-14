package pnv10n

import (
	"github.com/kadisoka/kadisoka-framework/pkg/volib/pkg/telephony"
)

type SMSDeliveryService interface {
	SendTextMessage(
		recipient telephony.PhoneNumber,
		text string,
		opts SMSDeliveryOptions) error
}

type SMSDeliveryOptions struct{}
