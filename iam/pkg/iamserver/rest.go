package iamserver

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/alloyzeus/go-azfl/errors"
	dataerrs "github.com/alloyzeus/go-azfl/errors/data"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

type RESTServiceServerBase struct {
	*Core
}

var _ iam.ConsumerRESTServer = &RESTServiceServerBase{}

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
) (app *iam.Application, err error) {
	authorizationHeader := req.Header.Get(iam.AuthorizationMetadataKey)
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
		return nil, iam.ReqFieldErr(iam.AuthorizationMetadataKey, dataerrs.Malformed(err))
	}

	creds := strings.SplitN(string(credsBytes), ":", 2)
	if creds[0] == "" {
		return nil, iam.ReqFieldErr(iam.AuthorizationMetadataKey, errors.EntMsg("username", "empty"))
	}

	appID, err := iam.ApplicationIDFromAZIDText(creds[0])
	if err != nil {
		return nil, iam.ReqFieldErr(iam.AuthorizationMetadataKey, errors.Ent("username", dataerrs.Malformed(err)))
	}
	if appID.IsNotStaticallyValid() {
		return nil, iam.ReqFieldErr(iam.AuthorizationMetadataKey, errors.Ent("username", nil))
	}

	app, err = svcBase.AuthenticatedApplication(appID, creds[1])
	if err != nil {
		return nil, errors.Wrap("client look up", err)
	}
	if app == nil {
		return nil, iam.ReqFieldErr(iam.AuthorizationMetadataKey, errors.EntMsg("username", "reference invalid"))
	}

	return app, nil
}
