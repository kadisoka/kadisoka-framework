package iam

import (
	"testing"

	"github.com/alloyzeus/go-azfl/azid"
	"github.com/stretchr/testify/assert"
)

func TestUserIDNumLimits(t *testing.T) {
	assert.Equal(t, UserIDNum(0), UserIDNumZero, "zero equal")
	assert.Equal(t, false, UserIDNum(0).IsStaticallyValid(), "zero")
	assert.Equal(t, false, UserIDNum(-1).IsStaticallyValid(), "neg is invalid")
	assert.Equal(t, true, UserIDNum(1).IsStaticallyValid(), "")
	assert.Equal(t, true, UserIDNum(0xffff).IsStaticallyValid(), "")
	assert.Equal(t, false, UserIDNum(0x0001000000000000).IsStaticallyValid(), "over limit")
	assert.Equal(t, true, UserIDNum(4294967296).IsStaticallyValid(), "lowest normal")
	assert.Equal(t, true, UserIDNum(4294967296).IsNormalAccount(), "lowest normal")
	assert.Equal(t, false, UserIDNum(4294967296).IsBot(), "lowest normal")
}

func TestUserIDNumEncode(t *testing.T) {
	testCases := []struct {
		input         UserIDNum
		expectedBytes []byte
		expectedType  azid.BinDataType
		label         string
	}{
		{UserIDNumZero, []byte{0, 0, 0, 0, 0, 0, 0, 0}, azid.BinDataTypeInt64, "zero is empty"},
		{UserIDNum(-1), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, azid.BinDataTypeInt64, "neg"},
	}

	for _, testData := range testCases {
		d, typ := testData.input.AZIDBinField()
		assert.Equal(t, testData.expectedBytes, d, testData.label)
		assert.Equal(t, testData.expectedType, typ, testData.label)
	}
	// assert.Equal(t, []byte(""), UserIDNumZero.AZIDBinField(), "zero is empty")
	// assert.Equal(t, "", UserIDNum(0).String(), "zero is empty")
	// assert.Equal(t, "", UserIDNum(-1).String(), "neg is empty")
	// assert.Equal(t, "", UserIDNum(1).String(), "reserved is empty")
	// assert.Equal(t, "ISv0T2000", UserIDNum(0x10000).String(), "service account")
	// assert.Equal(t, "INo0T4000000", UserIDNum(4294967296).String(), "normal account")
	// assert.Equal(t, "INo0T7zz6ya1v0x", UserIDNum(281448076602397).String(), "normal account")
	//TODO: more cases
}

func TestUserIDNumDecode(t *testing.T) {
	// var cases = []struct {
	// 	encoded  string
	// 	expected UserIDNum
	// 	err      error
	// }{
	// 	{"", UserIDNumZero, nil},
	// 	{"ISv0T2000", UserIDNum(0x10000), nil},
	// 	{"INo0T4000000", UserIDNum(4294967296), nil},
	// 	{"INo0T7zz6ya1v0x", UserIDNum(281448076602397), nil},
	// }

	// for _, c := range cases {
	// 	uid, err := UserIDNumFromString(c.encoded)
	// 	assert.Equal(t, c.err, err, "error")
	// 	assert.Equal(t, c.expected, uid, "uid")
	// }
}

//TODO: more tests
