package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Invalid number of arguments\n")
		os.Exit(-1)
	}
	firstParty := os.Args[1] == "true"
	clientID := GenerateApplicationRefKey(firstParty, os.Args[2])
	clientSecret := genSecret(16)
	fmt.Fprintf(os.Stdout, "%s\n%s\n", clientID.AZIDText(), clientSecret)
}

// GenerateApplicationRefKey generates a new ApplicationRefKey. Note that this function is
// not consulting any database. To ensure that the generated ApplicationRefKey is
// unique, check the client database.
func GenerateApplicationRefKey(firstParty bool, clientTyp string) iam.ApplicationRefKey {
	var typeInfo uint32
	if firstParty {
		typeInfo = iam.ApplicationIDNumFirstPartyBits
	}
	switch clientTyp {
	case "service":
		typeInfo |= iam.ApplicationIDNumServiceBits
	case "ua-public":
		typeInfo |= iam.ApplicationIDNumUserAgentAuthorizationPublicBits
	case "ua-confidential":
		typeInfo |= iam.ApplicationIDNumUserAgentAuthorizationConfidentialBits
	default:
		panic("Unsupported client app type")
	}
	//TODO: reserve some ranges (?)
	appIDNum, err := iamserver.GenerateApplicationIDNum(typeInfo)
	if err != nil {
		panic(err)
	}
	return iam.NewApplicationRefKey(appIDNum)
}

func genSecret(len int) string {
	if len == 0 {
		len = 16
	}
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
