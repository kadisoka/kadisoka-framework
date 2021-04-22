package email

import (
	"regexp"
	"strings"

	azfl "github.com/alloyzeus/go-azfl/azfl"
	"github.com/alloyzeus/go-azfl/azfl/errors"
	dataerrs "github.com/alloyzeus/go-azfl/azfl/errors/data"
)

type Address struct {
	localPart  string
	domainPart string
	rawInput   string
}

var _ azfl.ValueObject = Address{}

func AddressFromString(str string) (Address, error) {
	parts := strings.SplitN(str, "@", 2)
	if len(parts) < 2 {
		return Address{}, dataerrs.ErrMalformed
	}
	//TODO(exa): normalize localPart and domainPart
	if parts[0] == "" {
		return Address{}, errors.EntMsg("local part", "empty")
	}
	if parts[1] == "" || !addressDomainRE.MatchString(parts[1]) {
		return Address{}, errors.Ent("domain part", nil)
	}
	//TODO(exa): perform more extensive checking

	return Address{
		localPart:  parts[0],
		domainPart: strings.ToLower(parts[1]),
		rawInput:   str,
	}, nil
}

//TODO: at least common address convention
func (addr Address) IsSound() bool {
	return addr.localPart != "" && addr.domainPart != ""
}

func (addr Address) Equal(other interface{}) bool {
	return addr.Equals(other)
}

func (addr Address) Equals(other interface{}) bool {
	//TODO: compare the normalized representations
	if o, ok := other.(Address); ok {
		return strings.EqualFold(o.domainPart, addr.domainPart) &&
			o.localPart == addr.localPart
	}
	if o, _ := other.(*Address); o != nil {
		return strings.EqualFold(o.domainPart, addr.domainPart) &&
			o.localPart == addr.localPart
	}
	return false
}

func (addr Address) String() string {
	return addr.localPart + "@" + addr.domainPart
}

func (addr Address) LocalPart() string {
	return addr.localPart
}

func (addr Address) DomainPart() string {
	return addr.domainPart
}

func (addr Address) RawInput() string {
	return addr.rawInput
}

//NOTE: actually, it's not recommended to use regex to
// identify if a string is an email address:
// https://www.regular-expressions.info/email.html
var addressRE = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var addressDomainRE = regexp.MustCompile("^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsValidAddress(str string) bool {
	return addressRE.MatchString(str)
}
