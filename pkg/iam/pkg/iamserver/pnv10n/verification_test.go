package pnv10n

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerificationMethodsValidity(t *testing.T) {
	assert.False(t, VerificationMethodUnspecified.IsValid())
	assert.False(t, VerificationMethodUnknown.IsValid())
	assert.True(t, VerificationMethodNone.IsValid())
	assert.True(t, VerificationMethodSMS.IsValid())
}

func TestVerificationMethodsFromStrings(t *testing.T) {
	assert.Equal(t, VerificationMethodFromString(""), VerificationMethodUnspecified)
	assert.Equal(t, VerificationMethodFromString("none"), VerificationMethodNone)
	assert.Equal(t, VerificationMethodFromString("sms"), VerificationMethodSMS)
	assert.Equal(t, VerificationMethodFromString("call"), VerificationMethodUnknown)
	assert.Equal(t, VerificationMethodFromString("email"), VerificationMethodUnknown)
}
