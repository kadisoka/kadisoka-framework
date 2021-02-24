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

func TestApplicationRefKeyAZRSEncoding(t *testing.T) {
	assert.Equal(t, "", _ApplicationRefKeyZero.AZRS())
	assert.Equal(t, "KAp0201", ApplicationRefKey(1).AZRS())
	assert.Equal(t, "KAp0HT07", ApplicationRefKey(1000).AZRS())
	assert.Equal(t, "KAp0HRG9", ApplicationRefKey(1250).AZRS())
	assert.Equal(t, "KAp08G6081007", ApplicationRefKey(0x70000001).AZRS())
	assert.Equal(t, "KAp08T2J81007", ApplicationRefKey(0x70001250).AZRS())
	assert.Equal(t, "KAp08RQ389007", ApplicationRefKey(0x70012345).AZRS())
	assert.Equal(t, "KAp08Z2P93007", ApplicationRefKey(0x70045678).AZRS())
}

func TestApplicationRefKeyAZRSDecodingEmpty(t *testing.T) {
	refKey, err := ApplicationRefKeyFromAZRS("")
	assert.Nil(t, err)
	assert.Equal(t, _ApplicationRefKeyZero, refKey)
	assert.Equal(t, true, refKey.Equal(_ApplicationRefKeyZero))
	assert.Equal(t, false, refKey.IsValid())
	assert.Equal(t, true, refKey.IsZero())
	assert.Equal(t, true, refKey.EqualsApplicationRefKey(_ApplicationRefKeyZero))
}

func TestApplicationRefKeyAZRSDecodingValid(t *testing.T) {
	refKey, err := ApplicationRefKeyFromAZRS("KAp0ht07")
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
	assert.Equal(t, `"KAp08RQ389007"`, string(b))
}

func TestApplicationRefKeyJSONDecodingValid(t *testing.T) {
	var refKey ApplicationRefKey
	err := json.Unmarshal([]byte(`"KAp08RQ389007"`), &refKey)
	assert.Nil(t, err)
	assert.Equal(t, true, refKey.Equals(ApplicationRefKey(0x70012345)))
}
