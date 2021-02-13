package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/logging"
)

var log = logging.NewPkgLogger()

// RunServers run all provided servers. It won't return until it received
// a signal, SIGTERM or SIGINT by default, and all servers are stopped.
func RunServers(servers []ServiceServer, shutdownSignal <-chan os.Signal) {
	if len(servers) == 0 {
		return
	}

	// Used to determine if all servers have stopped.
	var serversStopWaiter sync.WaitGroup

	// Start the servers
	for _, srv := range servers {
		serversStopWaiter.Add(1)
		go func(innerSrv ServiceServer) {
			srvName := innerSrv.ServiceInfo().Name
			log.Info().Msgf("Starting server %s...", srvName)
			err := innerSrv.Serve()
			if err != nil {
				log.Fatal().Err(err).Msgf("%s serve", srvName)
			} else {
				log.Info().Msgf("%s stopped", srvName)
			}
			serversStopWaiter.Done()
		}(srv)
	}

	// We set up the signal handler (interrupt and terminate).
	// We are using the signal to gracefully and forcefully stop the server.
	if shutdownSignal == nil {
		sigChan := make(chan os.Signal)
		// Listen to interrupt and terminal signals
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		shutdownSignal = sigChan
	}

	// Wait for the shutdown signal
	<-shutdownSignal

	log.Info().Msg("Shutting down servers...")

	// Start another go-routine to catch another signal so the shutdown
	// could be forced. If we get another signal, we'll exit immediately.
	go func() {
		<-shutdownSignal
		log.Info().Msg("Forced shutdown.")
		os.Exit(0)
	}()

	// Gracefully shutdown the servers
	shutdownCtx, shutdownCtxCancel := context.WithTimeout(
		context.Background(), 15*time.Second)
	defer shutdownCtxCancel()

	for _, srv := range servers {
		go func(innerSrv ServiceServer) {
			srvName := innerSrv.ServiceInfo().Name
			log.Info().Msgf("Shutting down server %s...", srvName)
			err := innerSrv.Shutdown(shutdownCtx)
			if err != nil {
				log.Err(err).Msgf("Server %s shutdown with error", srvName)
			}
		}(srv)
	}

	// Wait for all servers to stop.
	serversStopWaiter.Wait()

	log.Info().Msg("Done.")
}
