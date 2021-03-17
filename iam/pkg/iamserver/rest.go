package iamserver

import (
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	dataerrs "github.com/alloyzeus/go-azfl/azfl/errors/data"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

type RESTServiceServerBase struct {
	*Core
}

var _ iam.RESTServiceClientServer = &RESTServiceServerBase{}

func RESTServiceServerWith(iamServerCore *Core) *RESTServiceServerBase {
	if iamServerCore == nil {
		panic("provided iamServerCore is nil")
	}
	return &RESTServiceServerBase{iamServerCore}
}

// RequestApplication returns a Client info which identified by Basic
// authorization header field.
//
// If the authorization is not provided, the returned client will be nil,
// and the err value will be nil.
//
// If the authorization is provided and it's invalid, the returned client
// will be nil and err value will contain the information about why it
// failed.
//
// If the authorization is provided and it's valid, the returned client
// will be a valid client and err will be nil.
func (svcBase *RESTServiceServerBase) RequestApplication(
	req *http.Request,
) (client *iam.Application, err error) {
	authorizationHeader := req.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, nil
	}

	authorizationParts := strings.SplitN(authorizationHeader, " ", 2)
	if len(authorizationParts) != 2 {
		return nil, iam.ErrReqFieldAuthorizationMalformed
	}
	if authorizationParts[0] != "Basic" {
		return nil, iam.ErrReqFieldAuthorizationTypeUnsupported
	}

	credsBytes, err := base64.StdEncoding.
		DecodeString(strings.TrimSpace(authorizationParts[1]))
	if err != nil {
		return nil, iam.ReqFieldErr("Authorization", dataerrs.Malformed(err))
	}

	creds := strings.SplitN(string(credsBytes), ":", 2)
	if creds[0] == "" {
		return nil, iam.ReqFieldErr("Authorization", errors.EntMsg("username", "empty"))
	}
	appRef, err := iam.ApplicationRefKeyFromAZERText(creds[0])
	if err != nil {
		return nil, iam.ReqFieldErr("Authorization", errors.Ent("username", dataerrs.Malformed(err)))
	}
	if appRef.IsNotValid() {
		return nil, iam.ReqFieldErr("Authorization", errors.Ent("username", nil))
	}

	client, err = svcBase.ApplicationByRefKey(appRef)
	if err != nil {
		return nil, errors.Wrap("client look up", err)
	}
	if client == nil {
		return nil, iam.ReqFieldErr("Authorization", errors.EntMsg("username", "reference invalid"))
	}
	if len(creds) == 0 || subtle.ConstantTimeCompare([]byte(creds[1]), []byte(client.Data.Secret)) != 1 {
		return nil, iam.ReqFieldErr("Authorization", errors.EntMsg("password", "mismatch"))
	}

	return client, nil
}
