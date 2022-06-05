package pnv10n

import (
	"github.com/kadisoka/kadisoka-framework/volib/pkg/telephony"
)

type SMSDeliveryService interface {
	SendTextMessage(
		recipient telephony.PhoneNumber,
		text string,
		opts SMSDeliveryOptions) error
}

type SMSDeliveryOptions struct{}
