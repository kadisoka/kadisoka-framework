package iamserver

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/kadisoka/foundation/pkg/errors"
	dataerrs "github.com/kadisoka/foundation/pkg/errors/data"

	"github.com/kadisoka/iam/pkg/iam"
)

type RESTServiceServerBase struct {
	*Core
}

func RESTServiceServerWith(iamServerCore *Core) *RESTServiceServerBase {
	if iamServerCore == nil {
		panic("provided iamServerCore is nil")
	}
	return &RESTServiceServerBase{iamServerCore}
}

// RequestClient returns a Client info which identified by Basic authorization
// header.
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
func (svcBase *RESTServiceServerBase) RequestClient(
	req *http.Request,
) (client *iam.Client, err error) {
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
	clientID, err := iam.ClientIDFromString(creds[0])
	if err != nil {
		return nil, iam.ReqFieldErr("Authorization", errors.Ent("username", dataerrs.Malformed(err)))
	}
	if clientID.IsNotValid() {
		return nil, iam.ReqFieldErr("Authorization", errors.Ent("username", nil))
	}

	client, err = svcBase.ClientByID(clientID)
	if err != nil {
		return nil, errors.Wrap("client look up", err)
	}
	if client == nil {
		return nil, iam.ReqFieldErr("Authorization", errors.EntMsg("username", "reference invalid"))
	}
	if len(creds) == 0 || creds[1] != client.Secret {
		return nil, iam.ReqFieldErr("Authorization", errors.EntMsg("password", "mismatch"))
	}

	return client, nil
}

// RequestHasNoAuthorization checks if the request header with key
// Authorization has non empty value. This DOES NOT check if it's valid
// or not.
func (svcBase *RESTServiceServerBase) RequestHasNoAuthorization(req *http.Request) bool {
	if req == nil {
		return true
	}
	authHeaderVal := req.Header.Get("Authorization")
	return authHeaderVal == ""
}
