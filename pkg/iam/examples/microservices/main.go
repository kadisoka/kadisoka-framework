package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	restfulopenapi "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	_ "github.com/lib/pq"

	"github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/api/rest"
	"github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/app"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam/rest/logging"
)

var (
	appName        = "Kadisoka IAM Test Client Application"
	revisionID     = "unknown"
	buildTimestamp = "unknown"
)

var (
	log    = logging.NewPkgLogger()
	logCtx = log.WithContext
)

func main() {
	appInfo := app.Info{
		Name: appName,
		BuildInfo: app.BuildInfo{
			RevisionID: revisionID,
			Timestamp:  buildTimestamp,
		},
	}

	fmt.Fprintf(os.Stderr, "%s\n", appInfo.HeaderString())

	appBase, err := app.Init(appInfo)
	if err != nil {
		log.Fatal().Err(err).Msg("App initialization")
	}

	log.Info().Msg("Initializing app...")
	svcApp, err := iam.NewConsumerServerAppSimple(appBase, "IAM_")
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	log.Info().Msgf("App instance %v terminal %v", svcApp.InstanceID(), svcApp.TerminalID().AZIDText())

	httpListenPort := "8080"

	restV1BasePath := "/microservices/restv1"
	restV1DocsPath := restV1BasePath + "/apidocs.json"

	restV1Mux := http.NewServeMux()
	restV1Container := restful.NewContainer()
	restV1Container.ServeMux = restV1Mux
	restV1Container.EnableContentEncoding(true)
	// We need CORS for our web clients
	rest.SetUpCORSFilterByEnv(restV1Container, "CORS_", nil)

	restV1Svc := NewRESTService(svcApp, restV1BasePath)
	restV1Container.Add(restV1Svc.RestfulWebService())

	// Setup API specification handler
	restV1Container.Add(restfulopenapi.NewOpenAPIService(restfulopenapi.Config{
		WebServices:                   restV1Container.RegisteredWebServices(),
		APIPath:                       restV1DocsPath,
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}))

	restV1Mux.Handle(restV1BasePath+"/apidocs/",
		http.StripPrefix(restV1BasePath+"/apidocs/",
			http.FileServer(http.Dir("resources/swagger-ui"))))

	httpRootMux := http.NewServeMux()
	// Health check is used by load balancer and/or orchestrator
	httpRootMux.HandleFunc("/healthz",
		func(w http.ResponseWriter, _ *http.Request) {
			log.Debug().Msg("TODO: /healthz actual health check e.g., db and other services we depended on")
			w.Write([]byte("OK"))
		})

	statsFilter := rest.NewStatsFilter()
	restV1Container.Filter(statsFilter.Filter)
	httpRootMux.HandleFunc(restV1BasePath+"/statz", statsFilter.StatsHandler)
	httpRootMux.Handle("/", restV1Mux)

	//httpRootMux.Handle(uiBaseURL+"/", NewUIService(uiBaseURL))

	httpServer := &http.Server{
		Addr:    ":" + httpListenPort,
		Handler: httpRootMux}

	// For this point forward is about running the server.
	// First, we set up the signal handler (interrupt and terminate).
	// We are using the signal to gracefully and forcefully stop the server.
	var shuttingDown bool
	shutdownSignal := make(chan os.Signal)
	// Listen to interrupt and terminal signals
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	// Start a go-routine to run the server.
	go func() {
		log.Info().Msgf("API specification: %s", restV1DocsPath)
		log.Info().Msg("Service is ready")
		err = httpServer.ListenAndServe()
		if err != nil && (err != http.ErrServerClosed || !shuttingDown) {
			log.Error().Err(err).
				Msg("API server error")
		}
	}()

	// Wait for any signal
	<-shutdownSignal
	shuttingDown = true
	log.Info().Msg("Shutting down the server...")

	// Start another go-routine to catch another signal so the shutdown
	// could be forced. If we get another signal, we'll exit immediately.
	go func() {
		<-shutdownSignal
		os.Exit(0)
	}()

	// Gracefully shutdown the server.
	shutdownCtx, shutdownCtxCancel := context.WithTimeout(
		context.Background(), 15*time.Second)
	defer shutdownCtxCancel()

	httpServer.Shutdown(shutdownCtx)
	log.Info().Msg("Done.")
}

func enrichSwaggerObject(swo *spec.Swagger) {
	rev := revisionID
	if rev != "unknown" && len(rev) > 7 {
		rev = rev[:7]
	}
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Kadisoka IAM Test Client Service",
			Description: "A service which demonstrates 3-legged authorization",
			Version:     fmt.Sprintf("0.1.0-%s built at %s", rev, buildTimestamp),
		},
	}
}
