package iamserver

import (
	"encoding/csv"
	"errors"
	"os"
	"strings"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

type clientStaticDataProvider struct {
	clients map[iam.ApplicationRefKey]*iam.ApplicationData
}

func newClientStaticDataProviderFromCSVFileByName(
	filename string, skipRows int,
) (*clientStaticDataProvider, error) {
	csvFile, err := os.Open(filename)
	if err != nil {
		//TODO: translate errors
		return nil, err
	}
	defer csvFile.Close()

	rows, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		//TODO: translate errors
		return nil, err
	}

	if len(rows) < (skipRows) {
		return nil, errors.New("header row required")
	}

	displayNameIdx := -1
	secretIdx := -1
	platformTypeIdx := -1
	requiredScopesIdx := -1
	oauth2RedirectURIIdx := -1

	for idx, key := range rows[0] {
		switch key {
		case "display_name":
			displayNameIdx = idx
		case "secret":
			secretIdx = idx
		case "platform_type":
			platformTypeIdx = idx
		case "required_scopes":
			requiredScopesIdx = idx
		case "oauth2_redirect_uri":
			oauth2RedirectURIIdx = idx
		}
	}

	indexexdValue := func(ls []string, idx int) string {
		if idx < 0 {
			return ""
		}
		if idx >= len(ls) {
			return ""
		}
		return ls[idx]
	}

	clList := map[iam.ApplicationRefKey]*iam.ApplicationData{}
	for _, r := range rows[skipRows:] {
		var clID iam.ApplicationRefKey
		clID, err = iam.ApplicationRefKeyFromAZERText(r[0])
		if err != nil {
			return nil, err
		}

		var requiredScopes []string
		if requiredScopeStr := indexexdValue(r, requiredScopesIdx); requiredScopeStr != "" {
			parts := strings.Split(requiredScopeStr, " ")
			if len(parts) == 1 {
				parts = strings.Split(requiredScopeStr, ",")
			}
			if len(parts) > 1 {
				for _, v := range parts {
					scopeStr := strings.TrimSpace(v)
					if scopeStr != "" {
						requiredScopes = append(requiredScopes, scopeStr)
					}
				}
			} else {
				requiredScopes = append(requiredScopes, requiredScopeStr)
			}
		}

		var redirectURIs []string
		if redirectURIStr := indexexdValue(r, oauth2RedirectURIIdx); redirectURIStr != "" {
			parts := strings.Split(redirectURIStr, ",")
			if len(parts) > 1 {
				for _, v := range parts {
					uriStr := strings.TrimSpace(v)
					if uriStr != "" {
						redirectURIs = append(redirectURIs, uriStr)
					}
				}
			} else {
				redirectURIs = append(redirectURIs, redirectURIStr)
			}
		}

		clList[clID] = &iam.ApplicationData{
			DisplayName:       indexexdValue(r, displayNameIdx),
			Secret:            indexexdValue(r, secretIdx),
			PlatformType:      indexexdValue(r, platformTypeIdx),
			RequiredScopes:    requiredScopes,
			OAuth2RedirectURI: redirectURIs,
		}
	}

	return &clientStaticDataProvider{clList}, nil
}

func (clientStaticDataStore *clientStaticDataProvider) GetApplication(
	appID iam.ApplicationRefKey,
) (*iam.Application, error) {
	cl := clientStaticDataStore.clients[appID]
	if cl == nil {
		return nil, nil
	}
	app := &iam.Application{ID: appID, Data: *cl}
	return app, nil
}
