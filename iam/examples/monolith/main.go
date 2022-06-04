package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/rez-go/stev"
	"github.com/rez-go/stev/docgen"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/app"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/realm"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/webui"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam/logging"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver"
	iamapp "github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/app"
)

var (
	appName        = "Kadisoka Monolith Example Server"
	revisionID     = "unknown"
	buildTimestamp = "unknown"
)

const envVarsPrefix = ""

var log = logging.NewPkgLogger()

func main() {
	appInfo := app.Info{
		Name: appName,
		BuildInfo: app.BuildInfo{
			RevisionID: revisionID,
			Timestamp:  buildTimestamp,
		},
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "env_file_template":
			fmt.Fprintf(os.Stderr, "Generating env file template...\n")
			cfg := configSkeleton()
			fmt.Fprintf(os.Stdout,
				"# Env file template generated\n# by %s\n",
				appInfo.HeaderString())
			err := docgen.WriteEnvTemplate(os.Stdout, &cfg,
				docgen.EnvTemplateWriteOptions{FieldPrefix: envVarsPrefix})
			if err != nil {
				log.Fatal().Err(err).Msg("Unable to generate env file template")
			} else {
				fmt.Fprintln(os.Stderr, "Done.")
			}
			return
		}
	}

	fmt.Fprintf(os.Stderr, "%s\n", appInfo.HeaderString())
	appBase, err := app.Init(appInfo)
	if err != nil {
		log.Fatal().Err(err).Msg("App initialization")
	}

	err = initApp(appBase)
	if err != nil {
		log.Fatal().Err(err).Msg("Servers initialization")
	}

	http.ListenAndServe(":8080", nil)
}

// Config is the configuration of our app. This config includes config for
// all subsystems in our application.
type Config struct {
	// All of IAM components configurations will be under namespace 'IAM' (i.e., prefixed with 'IAM_')
	IAM   iamapp.Config      `env:"IAM"`
	WebUI webui.ServerConfig `env:"WEBUI"`
}

func initApp(appBase app.App) error {
	var err error

	realmInfo := realm.Info{
		Name:    "Kadisoka Monolith Example",
		Contact: realm.ContactInfo{EmailAddress: "info@example.com"},
	}
	realmInfo, err = realm.InfoFromEnv("REALM_", &realmInfo)
	if err != nil {
		log.Fatal().Err(err).Msg("RealmInfo loading")
	}

	cfg := configSkeleton()
	cfg.IAM.RealmInfo = &realmInfo

	err = stev.LoadEnv("", &cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Config loading")
	}

	mux := http.DefaultServeMux

	// Init IAM core but don't let it init the services. We'll init
	// the services in our application.
	_, err = iamapp.NewWithCombinedHTTPServers(appBase, cfg.IAM, mux)
	if err != nil {
		log.Fatal().Err(err).Msg("IAM initialization")
	}

	webUIServer, err := webui.NewServer(
		cfg.WebUI,
		map[string]interface{}{})
	if err != nil {
		log.Fatal().Err(err).Msg("Web UI initialization")
	}
	mux.Handle("/", webUIServer)

	return nil
}

func configSkeleton() Config {
	curDir, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	return Config{
		WebUI: webui.ServerConfig{
			ServePath: "/",
			FilesDir:  filepath.Join(curDir, "resources", "monolith-webui"),
		},
		IAM: iamapp.Config{
			Core: iamserver.CoreConfigSkeleton(),
			// Serve HTTP services under /accounts
			HTTPBasePath: "/accounts",
			// Web UI
			WebUIEnabled: true,
			// REST API
			RESTEnabled: true,
		},
	}
}
