package iam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserIDNumLimits(t *testing.T) {
	assert.Equal(t, UserIDNum(0), UserIDNumZero, "zero equal")
	assert.Equal(t, false, UserIDNum(0).IsValid(), "zero")
	assert.Equal(t, false, UserIDNum(-1).IsValid(), "neg is invalid")
	assert.Equal(t, false, UserIDNum(1).IsValid(), "reserved")
	assert.Equal(t, false, UserIDNum(0xffff).IsValid(), "reserved")
	assert.Equal(t, false, UserIDNum(0x0001000000000000).IsValid(), "over limit")
	assert.Equal(t, true, UserIDNum(4294967296).IsValid(), "lowest normal")
	assert.Equal(t, true, UserIDNum(4294967296).IsNormalAccount(), "lowest normal")
	assert.Equal(t, false, UserIDNum(4294967296).IsServiceAccount(), "lowest normal")
}

func TestUserIDNumEncode(t *testing.T) {
	// assert.Equal(t, "", UserIDNumZero.String(), "zero is empty")
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
