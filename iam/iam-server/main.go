package main

import (
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"

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

func main() {
	fmt.Fprintf(os.Stderr,
		"%s revision %s built at %s\n",
		appName, revisionID, buildTimestamp)
	app.Init(app.Info{
		Name: appName,
		BuildInfo: app.BuildInfo{
			RevisionID: revisionID,
			Timestamp:  buildTimestamp,
		},
	})

	srvApp, err := initApp()
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

func initApp() (app.App, error) {
	envVarsPrefix := "IAM_"

	realmInfo, err := realm.InfoFromEnvOrDefault("")
	if err != nil {
		log.Fatal().Err(err).Msg("RealmInfo loading")
	}

	cfg := srvapp.Config{
		RealmInfo: &realmInfo,
		Core:      iamserver.CoreConfigSkeleton(),
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

	srvApp, err := srvapp.NewByEnv(envVarsPrefix, &cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("App initialization")
	}

	return srvApp, nil
}
