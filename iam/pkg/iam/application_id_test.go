package iam

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplicationIDLimits(t *testing.T) {
	assert.Equal(t, ApplicationID(0), ApplicationIDZero)
	assert.Equal(t, true, ApplicationIDZero.IsZero())
	assert.Equal(t, ApplicationIDFromPrimitiveValue(0), ApplicationIDZero)
	assert.Equal(t, int32(1), ApplicationID(1).PrimitiveValue())
	assert.Equal(t, false, ApplicationID(0).IsValid())
	assert.Equal(t, false, ApplicationID(-1).IsValid())
	assert.Equal(t, true, ApplicationID(1).IsValid())
	assert.Equal(t, true, ApplicationID(0xffff).IsValid())
	assert.Equal(t, true, ApplicationID(0xffffff).IsValid())
	assert.Equal(t, true, ApplicationID(0x7fffffff).IsValid())
	assert.Equal(t, false, ApplicationID(1<<28).IsValid())
	assert.Equal(t, true, ApplicationID((1<<30)|0x1).IsFirstParty())
	assert.Equal(t, true, ApplicationID(0x01000000).IsValid())
	assert.Equal(t, true, ApplicationID(0x01000001).IsValid())
	assert.Equal(t, true, ApplicationID(0x01ffffff).IsValid())
}

func TestApplicationRefKeyAZERTextEncoding(t *testing.T) {
	assert.Equal(t, "", _ApplicationRefKeyZero.AZERText())
	assert.Equal(t, "KAp02c000001", ApplicationRefKey(1).AZERText())
	assert.Equal(t, "KAp02c000002", ApplicationRefKey(2).AZERText())
	assert.Equal(t, "KAp02c000003", ApplicationRefKey(3).AZERText())
	assert.Equal(t, "KAp02c000004", ApplicationRefKey(4).AZERText())
	assert.Equal(t, "KAp02c000005", ApplicationRefKey(5).AZERText())
	assert.Equal(t, "KAp02c0000z8", ApplicationRefKey(1000).AZERText())
	assert.Equal(t, "KAp02c0000z9", ApplicationRefKey(1001).AZERText())
	assert.Equal(t, "KAp02c0000za", ApplicationRefKey(1002).AZERText())
	assert.Equal(t, "KAp02c0000zb", ApplicationRefKey(1003).AZERText())
	assert.Equal(t, "KAp02c0000zc", ApplicationRefKey(1004).AZERText())
	assert.Equal(t, "KAp02c0001ym", ApplicationRefKey(2004).AZERText())
	assert.Equal(t, "KAp02c000172", ApplicationRefKey(1250).AZERText())
	assert.Equal(t, "KAp02dr00001", ApplicationRefKey(0x70000001).AZERText())
	assert.Equal(t, "KAp02dr004jg", ApplicationRefKey(0x70001250).AZERText())
	assert.Equal(t, "KAp02dr028t5", ApplicationRefKey(0x70012345).AZERText())
	assert.Equal(t, "KAp02dr08nkr", ApplicationRefKey(0x70045678).AZERText())
}

func TestApplicationRefKeyAZERTextDecodingEmpty(t *testing.T) {
	refKey, err := ApplicationRefKeyFromAZERText("")
	assert.Nil(t, err)
	assert.Equal(t, _ApplicationRefKeyZero, refKey)
	assert.Equal(t, true, refKey.Equal(_ApplicationRefKeyZero))
	assert.Equal(t, false, refKey.IsValid())
	assert.Equal(t, true, refKey.IsZero())
	assert.Equal(t, true, refKey.EqualsApplicationRefKey(_ApplicationRefKeyZero))
}

func TestApplicationRefKeyAZERTextDecodingValid(t *testing.T) {
	refKey, err := ApplicationRefKeyFromAZERText("KAp02c0000z8")
	assert.Nil(t, err)
	assert.Equal(t, ApplicationRefKey(1000), refKey)
	assert.Equal(t, false, refKey.Equal(_ApplicationRefKeyZero))
	assert.Equal(t, true, refKey.IsValid())
	assert.Equal(t, false, refKey.IsZero())
	assert.Equal(t, true, refKey.EqualsApplicationRefKey(ApplicationRefKey(1000)))
}

func TestApplicationRefKeyJSONEncodingZero(t *testing.T) {
	refKey := ApplicationRefKeyZero()
	b, err := json.Marshal(&refKey)
	assert.Nil(t, err)
	assert.Equal(t, `""`, string(b))
}

func TestApplicationRefKeyJSONEncodingValid(t *testing.T) {
	refKey := ApplicationRefKey(0x70012345)
	b, err := json.Marshal(&refKey)
	assert.Nil(t, err)
	assert.Equal(t, `"KAp02dr028t5"`, string(b))
}

func TestApplicationRefKeyJSONDecodingValid(t *testing.T) {
	var refKey ApplicationRefKey
	err := json.Unmarshal([]byte(`"KAp02dr028t5"`), &refKey)
	assert.Nil(t, err)
	assert.Equal(t, true, refKey.Equals(ApplicationRefKey(0x70012345)))
}
