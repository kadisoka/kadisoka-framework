package telephony

import (
	"strconv"
	"strings"

	azfl "github.com/alloyzeus/go-azfl/azfl"
	"github.com/nyaruka/phonenumbers"
)

// PhoneNumber represents a phone number as we need.
type PhoneNumber struct {
	countryCode    int32
	nationalNumber int64
	rawInput       string
	isValid        bool
}

var _ azfl.ValueObject = PhoneNumber{}

func NewPhoneNumber(countryCode int32, nationalNumber int64) PhoneNumber {
	return PhoneNumber{countryCode: countryCode, nationalNumber: nationalNumber}
}

func PhoneNumberFromString(phoneNumberStr string) (PhoneNumber, error) {
	// Check if the country code is doubled
	if parts := strings.Split(phoneNumberStr, "+"); len(parts) == 3 {
		// We assume that the first part was automatically added at the client
		phoneNumberStr = "+" + parts[2]
	}

	parsedPhoneNumber, err := phonenumbers.Parse(phoneNumberStr, "")
	if err != nil {
		return PhoneNumber{}, err
	}

	phoneNumber := PhoneNumber{
		countryCode:    *parsedPhoneNumber.CountryCode,
		nationalNumber: int64(*parsedPhoneNumber.NationalNumber),
		rawInput:       phoneNumberStr,
		isValid:        phonenumbers.IsValidNumber(parsedPhoneNumber),
	}

	return phoneNumber, nil
}

func (phoneNumber PhoneNumber) IsSound() bool { return phoneNumber.isValid }

func (phoneNumber PhoneNumber) Equal(other interface{}) bool {
	return phoneNumber.Equals(other)
}

func (phoneNumber PhoneNumber) Equals(other interface{}) bool {
	if o, ok := other.(PhoneNumber); ok {
		return o.countryCode == phoneNumber.countryCode &&
			o.nationalNumber == phoneNumber.nationalNumber
	}
	if o, _ := other.(*PhoneNumber); o != nil {
		return o.countryCode == phoneNumber.countryCode &&
			o.nationalNumber == phoneNumber.nationalNumber
	}
	return false
}

func (phoneNumber PhoneNumber) CountryCode() int32    { return phoneNumber.countryCode }
func (phoneNumber PhoneNumber) NationalNumber() int64 { return phoneNumber.nationalNumber }
func (phoneNumber PhoneNumber) RawInput() string      { return phoneNumber.rawInput }

//TODO: get E.164 string
//TODO: consult the standards
func (phoneNumber PhoneNumber) String() string {
	if phoneNumber.countryCode == 0 && phoneNumber.nationalNumber == 0 {
		return "+"
	}
	return "+" + strconv.FormatInt(int64(phoneNumber.countryCode), 10) +
		strconv.FormatInt(phoneNumber.nationalNumber, 10)
}

// RawOrFormatted returns a string which prefers raw input with formatted
// string as the default.
func (phoneNumber PhoneNumber) RawOrFormatted() string {
	if phoneNumber.rawInput != "" {
		return phoneNumber.rawInput
	}
	return phoneNumber.String()
}
