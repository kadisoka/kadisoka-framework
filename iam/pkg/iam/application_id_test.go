package iam

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplicationIDNumLimits(t *testing.T) {
	assert.Equal(t, ApplicationIDNum(0), ApplicationIDNumZero)
	assert.Equal(t, true, ApplicationIDNumZero.IsZero())
	assert.Equal(t, ApplicationIDNumFromPrimitiveValue(0), ApplicationIDNumZero)
	assert.Equal(t, int32(1), ApplicationIDNum(1).PrimitiveValue())
	assert.Equal(t, false, ApplicationIDNum(0).IsStaticallyValid())
	assert.Equal(t, false, ApplicationIDNum(-1).IsStaticallyValid())
	assert.Equal(t, true, ApplicationIDNum(1).IsStaticallyValid())
	assert.Equal(t, true, ApplicationIDNum(0xffff).IsStaticallyValid())
	assert.Equal(t, true, ApplicationIDNum(0xffffff).IsStaticallyValid())
	assert.Equal(t, true, ApplicationIDNum(0x7fffffff).IsStaticallyValid())
	assert.Equal(t, false, ApplicationIDNum(1<<28).IsStaticallyValid())
	assert.Equal(t, true, ApplicationIDNum((1<<30)|0x1).IsFirstParty())
	assert.Equal(t, true, ApplicationIDNum(0x01000000).IsStaticallyValid())
	assert.Equal(t, true, ApplicationIDNum(0x01000001).IsStaticallyValid())
	assert.Equal(t, true, ApplicationIDNum(0x01ffffff).IsStaticallyValid())
}

func TestApplicationRefKeyAZIDTextEncoding(t *testing.T) {
	assert.Equal(t, "", _ApplicationRefKeyZero.AZIDText())
	assert.Equal(t, "KAp02c000001", ApplicationRefKey(1).AZIDText())
	assert.Equal(t, "KAp02c000002", ApplicationRefKey(2).AZIDText())
	assert.Equal(t, "KAp02c000003", ApplicationRefKey(3).AZIDText())
	assert.Equal(t, "KAp02c000004", ApplicationRefKey(4).AZIDText())
	assert.Equal(t, "KAp02c000005", ApplicationRefKey(5).AZIDText())
	assert.Equal(t, "KAp02c0000z8", ApplicationRefKey(1000).AZIDText())
	assert.Equal(t, "KAp02c0000z9", ApplicationRefKey(1001).AZIDText())
	assert.Equal(t, "KAp02c0000za", ApplicationRefKey(1002).AZIDText())
	assert.Equal(t, "KAp02c0000zb", ApplicationRefKey(1003).AZIDText())
	assert.Equal(t, "KAp02c0000zc", ApplicationRefKey(1004).AZIDText())
	assert.Equal(t, "KAp02c0001ym", ApplicationRefKey(2004).AZIDText())
	assert.Equal(t, "KAp02c000172", ApplicationRefKey(1250).AZIDText())
	assert.Equal(t, "KAp02dr00001", ApplicationRefKey(0x70000001).AZIDText())
	assert.Equal(t, "KAp02dr004jg", ApplicationRefKey(0x70001250).AZIDText())
	assert.Equal(t, "KAp02dr028t5", ApplicationRefKey(0x70012345).AZIDText())
	assert.Equal(t, "KAp02dr08nkr", ApplicationRefKey(0x70045678).AZIDText())
}

func TestApplicationRefKeyAZIDTextDecodingEmpty(t *testing.T) {
	refKey, err := ApplicationRefKeyFromAZIDText("")
	assert.Nil(t, err)
	assert.Equal(t, _ApplicationRefKeyZero, refKey)
	assert.Equal(t, true, refKey.Equal(_ApplicationRefKeyZero))
	assert.Equal(t, false, refKey.IsStaticallyValid())
	assert.Equal(t, true, refKey.IsZero())
	assert.Equal(t, true, refKey.EqualsApplicationRefKey(_ApplicationRefKeyZero))
}

func TestApplicationRefKeyAZIDTextDecodingValid(t *testing.T) {
	refKey, err := ApplicationRefKeyFromAZIDText("KAp02c0000z8")
	assert.Nil(t, err)
	assert.Equal(t, ApplicationRefKey(1000), refKey)
	assert.Equal(t, false, refKey.Equal(_ApplicationRefKeyZero))
	assert.Equal(t, true, refKey.IsStaticallyValid())
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
