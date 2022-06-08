package app

import (
	"net/http"
	"strings"

	"github.com/alloyzeus/go-azfl/errors"
	"github.com/rez-go/stev"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/app"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/realm"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/webui"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam/logging"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/grpc"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/rest"
	iamwebui "github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/webui"
)

var log = logging.NewPkgLogger()

func NewByEnv(
	appBase app.App,
	envVarsPrefix string,
	defaultConfig *Config,
) (*App, error) {
	cfg := defaultConfig
	if cfg == nil {
		cfg = &Config{}
	}
	err := stev.LoadEnv(envVarsPrefix, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("App config loading")
	}

	resolveConfig(cfg)

	log.Info().Msg("Initializing server app...")
	srvApp, err := newWithoutServices(appBase, *cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("App initialization")
	}

	err = srvApp.initServers(appBase, *cfg)
	if err != nil {
		return nil, err
	}

	return srvApp, nil
}

func NewWithCombinedHTTPServers(
	appBase app.App,
	cfg Config,
	mux *http.ServeMux,
) (*App, error) {
	resolveConfig(&cfg)

	srvApp, err := newWithoutServices(appBase, cfg)
	if err != nil {
		return nil, errors.Wrap("app initialization", err)
	}

	iamServerCore := srvApp.core

	if cfg.RESTEnabled {
		log.Info().Msg("Initializing REST server...")
		restServer, err := rest.NewServer(
			appBase,
			*cfg.REST,
			iamServerCore,
			&cfg.WebUI.URLs)
		if err != nil {
			return nil, errors.Wrap("REST server initialization", err)
		}

		mux.Handle(cfg.REST.ServePath, restServer)
	}

	if cfg.WebUIEnabled {
		log.Info().Msg("Initializing web UI server...")
		webUIServer, err := setUpWebUIServer(srvApp, cfg)
		if err != nil {
			return nil, errors.Wrap("web UI server initialization", err)
		}

		mux.Handle(cfg.WebUI.Server.ServePath, webUIServer)
	}

	return srvApp, nil
}

func setUpWebUIServer(srvApp *App, cfg Config) (*webui.Server, error) {
	webUICfg := cfg.WebUI.Server

	templateData := map[string]interface{}{
		"AppInfo":   srvApp.RealmInfo(),
		"RealmName": srvApp.RealmInfo().Name,
	}
	restAPIURLReplacer := &webui.StringReplacer{
		Old: "http://localhost:11121/rest/v1",
		New: strings.TrimRight(cfg.RESTCanonicalBaseURL, "/"),
	}
	webUIServeURLReplacer := &webui.StringReplacer{
		Old: "/kadisoka-iam-webui-base-path/",
		New: webUICfg.ServePath,
	}
	homeURLReplacer := &webui.StringReplacer{
		Old: "http://localhost:3000/",
		New: "/",
	}

	webUICfg.FileProcessors = map[string][]webui.FileProcessor{
		"*.html": {&webui.HTMLRenderer{
			Config: webui.HTMLRendererConfig{
				TemplateDelimBegin: "{:[",
				TemplateDelimEnd:   "]:}",
			},
			TemplateData: templateData,
		}, restAPIURLReplacer, webUIServeURLReplacer, homeURLReplacer},
		"*.js": {&webui.JSRenderer{
			Config: webui.JSRendererConfig{
				TemplateDelimBegin: "{:[",
				TemplateDelimEnd:   "]:}",
			},
			TemplateData: templateData,
		}, restAPIURLReplacer, webUIServeURLReplacer, homeURLReplacer},
	}

	webUIServer, err := webui.NewServer(
		webUICfg,
		templateData)
	if err != nil {
		return nil, errors.Wrap("web UI instantiation", err)
	}

	return webUIServer, nil
}

func newWithoutServices(appBase app.App, appCfg Config) (*App, error) {
	var realmInfo realm.Info
	if appCfg.RealmInfo != nil {
		realmInfo = *appCfg.RealmInfo
	}

	log.Info().Msg("Instantiating IAM Server Core...")
	srvCore, err := iamserver.NewCoreByConfig(appCfg.Core, appBase, realmInfo)
	if err != nil {
		return nil, errors.Wrap("core initialization", err)
	}

	return &App{
		App:  appBase,
		core: srvCore,
	}, nil
}

func resolveConfig(cfg *Config) {
	if cfg.WebUI == nil {
		cfg.WebUI = &iamwebui.ServerConfig{
			Server: webui.ServerConfig{},
		}
	}
	if cfg.WebUI.Server.ServePath == "" {
		cfg.WebUI.Server.ServePath = cfg.HTTPBasePath
	}
	cfg.WebUI.Server.ServePath = strings.TrimRight(cfg.WebUI.Server.ServePath, "/") + "/"
	if cfg.WebUI.Server.FilesDir == "" {
		cfg.WebUI.Server.FilesDir = "resources/iam-webui"
	}
	if cfg.WebUI.URLs.SignIn == "" {
		cfg.WebUI.URLs.SignIn = cfg.WebUI.Server.ServePath + "signin"
	}

	if cfg.REST == nil {
		cfg.REST = &rest.ServerConfig{}
	}
	if cfg.REST.ServePath == "" {
		cfg.REST.ServePath = strings.TrimRight(cfg.HTTPBasePath, "/") + "/rest/"
	}
	if cfg.RESTCanonicalBaseURL == "" {
		cfg.RESTCanonicalBaseURL = cfg.REST.ServePath + rest.ServerLatestVersionString + "/"
	}

	if cfg.Core.EAV.ResourcesDir == "" {
		cfg.Core.EAV.ResourcesDir = "resources/iam-pnv10n-resources"
	}
}

type App struct {
	app.App
	core *iamserver.Core
}

// RealmInfo returns information about the realm of the app.
func (srvApp App) RealmInfo() realm.Info { return srvApp.core.RealmInfo() }

func (srvApp *App) initServers(appBase app.App, cfg Config) error {
	iamServerCore := srvApp.core

	if cfg.RESTEnabled {
		log.Info().Msg("Initializing REST API server...")
		restServer, err := rest.NewServer(
			appBase,
			*cfg.REST,
			iamServerCore,
			&cfg.WebUI.URLs)
		if err != nil {
			return errors.Wrap("REST API server initialization", err)
		}

		srvApp.AddServiceServer(restServer)
	}

	if cfg.WebUIEnabled {
		log.Info().Msg("Initializing Web UI server...")
		webUIServer, err := setUpWebUIServer(srvApp, cfg)
		if err != nil {
			return errors.Wrap("Web UI server initialization", err)
		}

		srvApp.AddServiceServer(webUIServer)
	}

	if cfg.GRPCEnabled {
		log.Info().Msg("Initializing gRPC API server...")
		grpcServer, err := grpc.NewServer(
			*cfg.GRPC,
			iamServerCore)
		if err != nil {
			return errors.Wrap("gRPC API server initialization", err)
		}

		srvApp.AddServiceServer(grpcServer)
	}

	return nil

}

func ConfigSkeleton() Config {
	return Config{
		Core: iamserver.CoreConfigSkeleton(),
	}
}

func ConfigSkeletonPtr() *Config {
	cfg := ConfigSkeleton()
	return &cfg
}
