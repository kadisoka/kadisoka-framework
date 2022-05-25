package main

import (
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/rez-go/stev/docgen"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/app"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/realm"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/webui"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam/logging"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver"
	srvapp "github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/app"
	srvgrpc "github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/grpc"
	srvrest "github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/rest"
	srvwebui "github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/webui"
)

var log = logging.NewPkgLogger()

var (
	appName        = "Kadisoka IAM Server"
	revisionID     = "unknown"
	buildTimestamp = "unknown"
)

const envVarsPrefix = "IAM_"

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

	srvApp, err := initApp(appBase)
	if err != nil {
		log.Fatal().Err(err).Msg("Servers initialization")
	}

	// to detect that all services are ready
	go func() {
		for {
			time.Sleep(200 * time.Millisecond)
			if srvApp.IsAllServersAcceptingClients() {
				log.Info().Msg("Services are ready")
				break
			}
		}
	}()

	srvApp.Run()
}

func initApp(appBase app.App) (app.App, error) {
	realmInfo, err := realm.InfoFromEnvOrDefault("")
	if err != nil {
		log.Fatal().Err(err).Msg("RealmInfo loading")
	}

	cfg := configSkeleton()
	cfg.RealmInfo = &realmInfo

	srvApp, err := srvapp.NewByEnv(appBase, envVarsPrefix, &cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("App initialization")
	}

	return srvApp, nil
}

func configSkeleton() srvapp.Config {
	return srvapp.Config{
		Core: iamserver.CoreConfigSkeleton(),
		// Web UI
		WebUIEnabled: true,
		WebUI: &srvwebui.ServerConfig{
			Server: webui.ServerConfig{
				ServePort: 8080,
			},
		},
		// REST API
		RESTEnabled: true,
		REST: &srvrest.ServerConfig{
			ServePort:          9080,
			SwaggerUIAssetsDir: "resources/swagger-ui",
		},
		// gRPC API
		GRPCEnabled: false,
		GRPC: &srvgrpc.ServerConfig{
			ServePort: 50051,
		},
	}
}
