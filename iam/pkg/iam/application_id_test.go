package iam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplicationIDLimits(t *testing.T) {
	assert.Equal(t, ApplicationID(0), ApplicationIDZero)
	assert.Equal(t, false, ApplicationID(0).IsValid())
	assert.Equal(t, false, ApplicationID(-1).IsValid())
	assert.Equal(t, true, ApplicationID(1).IsValid())
	assert.Equal(t, true, ApplicationID(0xffff).IsValid())
	assert.Equal(t, true, ApplicationID(0xffffff).IsValid())
	assert.Equal(t, true, ApplicationID(0x7fffffff).IsValid())
	assert.Equal(t, false, ApplicationID(1<<28).IsValid())
	assert.Equal(t, true, ApplicationID(0x01000000).IsValid())
	assert.Equal(t, true, ApplicationID(0x01000001).IsValid())
	assert.Equal(t, true, ApplicationID(0x01ffffff).IsValid())
}

func TestApplicationRefKeyAZISEncoding(t *testing.T) {
	assert.Equal(t, "", _ApplicationRefKeyZero.AZIS())
	assert.Equal(t, "KAp00801", ApplicationRefKey(1).AZIS())
	assert.Equal(t, "KAp008e807", ApplicationRefKey(1000).AZIS())
}

func TestApplicationRefKeyAZISDecodingEmpty(t *testing.T) {
	refKey, err := ApplicationRefKeyFromAZIS("")
	assert.Nil(t, err)
	assert.Equal(t, _ApplicationRefKeyZero, refKey)
	assert.Equal(t, false, refKey.IsValid())
	assert.Equal(t, true, refKey.IsZero())
	assert.Equal(t, true, refKey.EqualsApplicationRefKey(_ApplicationRefKeyZero))
}

func TestApplicationRefKeyAZISDecodingValid(t *testing.T) {
	refKey, err := ApplicationRefKeyFromAZIS("KAp008e807")
	assert.Nil(t, err)
	assert.Equal(t, ApplicationRefKey(1000), refKey)
	assert.Equal(t, true, refKey.IsValid())
	assert.Equal(t, false, refKey.IsZero())
	assert.Equal(t, true, refKey.EqualsApplicationRefKey(ApplicationRefKey(1000)))
}
