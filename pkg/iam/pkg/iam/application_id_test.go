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

func TestApplicationIDAZIDTextEncoding(t *testing.T) {
	assert.Equal(t, "", _ApplicationIDZero.AZIDText())
	assert.Equal(t, "KAp02c000001", ApplicationID(1).AZIDText())
	assert.Equal(t, "KAp02c000002", ApplicationID(2).AZIDText())
	assert.Equal(t, "KAp02c000003", ApplicationID(3).AZIDText())
	assert.Equal(t, "KAp02c000004", ApplicationID(4).AZIDText())
	assert.Equal(t, "KAp02c000005", ApplicationID(5).AZIDText())
	assert.Equal(t, "KAp02c0000z8", ApplicationID(1000).AZIDText())
	assert.Equal(t, "KAp02c0000z9", ApplicationID(1001).AZIDText())
	assert.Equal(t, "KAp02c0000za", ApplicationID(1002).AZIDText())
	assert.Equal(t, "KAp02c0000zb", ApplicationID(1003).AZIDText())
	assert.Equal(t, "KAp02c0000zc", ApplicationID(1004).AZIDText())
	assert.Equal(t, "KAp02c0001ym", ApplicationID(2004).AZIDText())
	assert.Equal(t, "KAp02c000172", ApplicationID(1250).AZIDText())
	assert.Equal(t, "KAp02dr00001", ApplicationID(0x70000001).AZIDText())
	assert.Equal(t, "KAp02dr004jg", ApplicationID(0x70001250).AZIDText())
	assert.Equal(t, "KAp02dr028t5", ApplicationID(0x70012345).AZIDText())
	assert.Equal(t, "KAp02dr08nkr", ApplicationID(0x70045678).AZIDText())
}

func TestApplicationIDAZIDTextDecodingEmpty(t *testing.T) {
	id, err := ApplicationIDFromAZIDText("")
	assert.Nil(t, err)
	assert.Equal(t, _ApplicationIDZero, id)
	assert.Equal(t, true, id.Equal(_ApplicationIDZero))
	assert.Equal(t, false, id.IsStaticallyValid())
	assert.Equal(t, true, id.IsZero())
	assert.Equal(t, true, id.EqualsApplicationID(_ApplicationIDZero))
}

func TestApplicationIDAZIDTextDecodingValid(t *testing.T) {
	id, err := ApplicationIDFromAZIDText("KAp02c0000z8")
	assert.Nil(t, err)
	assert.Equal(t, ApplicationID(1000), id)
	assert.Equal(t, false, id.Equal(_ApplicationIDZero))
	assert.Equal(t, true, id.IsStaticallyValid())
	assert.Equal(t, false, id.IsZero())
	assert.Equal(t, true, id.EqualsApplicationID(ApplicationID(1000)))
}

func TestApplicationIDJSONEncodingZero(t *testing.T) {
	id := ApplicationIDZero()
	b, err := json.Marshal(&id)
	assert.Nil(t, err)
	assert.Equal(t, `""`, string(b))
}

func TestApplicationIDJSONEncodingValid(t *testing.T) {
	id := ApplicationID(0x70012345)
	b, err := json.Marshal(&id)
	assert.Nil(t, err)
	assert.Equal(t, `"KAp02dr028t5"`, string(b))
}

func TestApplicationIDJSONDecodingValid(t *testing.T) {
	var id ApplicationID
	err := json.Unmarshal([]byte(`"KAp02dr028t5"`), &id)
	assert.Nil(t, err)
	assert.Equal(t, true, id.Equals(ApplicationID(0x70012345)))
}
