package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"

	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iamserver"
)

func main() {
	var params Params
	err := survey.Ask(qs, &params)
	if err != nil {
		panic(err)
	}
	clientID := GenerateApplicationID(params.FirstParty, params.AppType)
	clientSecret := genSecret(16)
	fmt.Fprintf(os.Stdout, "%s\n%s\n", clientID.AZIDText(), clientSecret)
}

// GenerateApplicationID generates a new ApplicationID. Note that this function is
// not consulting any database. To ensure that the generated ApplicationID is
// unique, check the client database.
func GenerateApplicationID(firstParty bool, clientTyp string) iam.ApplicationID {
	var typeInfo uint32
	if firstParty {
		typeInfo = iam.ApplicationIDNumFirstPartyBits
	}
	switch clientTyp {
	case "service":
		typeInfo |= iam.ApplicationIDNumServiceBits
	case "ua-public", "user-agent-public", "user-agent-direct-auth":
		typeInfo |= iam.ApplicationIDNumUserAgentAuthorizationPublicBits
	case "ua-confidential", "user-agent-confidential", "user-agent-3-legged-auth":
		typeInfo |= iam.ApplicationIDNumUserAgentAuthorizationConfidentialBits
	default:
		panic("Unsupported client app type")
	}
	//TODO: reserve some ranges (?)
	appIDNum, err := iamserver.GenerateApplicationIDNum(typeInfo)
	if err != nil {
		panic(err)
	}
	return iam.NewApplicationID(appIDNum)
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

// the questions to ask
var qs = []*survey.Question{
	{
		Name:     "firstParty",
		Prompt:   &survey.Confirm{Message: "Is this first-party application?"},
		Validate: survey.Required,
	},
	{
		Name: "appType",
		Prompt: &survey.Select{
			Message: "The type op application:",
			Options: []string{
				"service", "user-agent-direct-auth", "user-agent-3-legged-auth",
			},
			Default: "user-agent-direct-auth",
		},
	},
}

type Params struct {
	FirstParty bool   `survey:"firstParty"`
	AppType    string `survey:"appType"`
}
