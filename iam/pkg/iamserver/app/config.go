package app

import (
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/app"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/realm"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/grpc"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/rest"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/webui"
)

type Config struct {
	RealmInfo *realm.Info `env:"REALM"`

	AppInfo *app.Info            `env:"APP"`
	Core    iamserver.CoreConfig `env:",squash"`

	HTTPBasePath         string `env:"HTTP_BASE_PATH"`
	RESTCanonicalBaseURL string `env:"REST_CANONICAL_BASE_URL"`

	WebUIEnabled bool                `env:"WEBUI_ENABLED"`
	WebUI        *webui.ServerConfig `env:"WEBUI"`
	RESTEnabled  bool                `env:"REST_ENABLED"`
	REST         *rest.ServerConfig  `env:"REST"`
	GRPCEnabled  bool                `env:"GRPC_ENABLED"`
	GRPC         *grpc.ServerConfig  `env:"GRPC"`
}

func (Config) EnvVarsDocs() map[string]string {
	return map[string]string{
		"HTTPBasePath": "The default serving path for all HTTP services, e.g., " +
			"web UI and REST. Individual HTTP services might accept specific " +
			"serve path in their configurations.",
		"WebUIEnabled": "The master flag for enabling/disabling serving the " +
			"web UI service.",
	}
}
