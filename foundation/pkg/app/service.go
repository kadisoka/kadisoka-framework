package app

import (
	"context"
)

// ServiceServer abstracts all service servers
//TODO: ServiceInfo
type ServiceServer interface {
	// ServiceInfo returns basic information about the service.
	ServiceInfo() ServiceInfo

	// Serve starts the server. This method is blocking and won't return
	// until the server is stopped (e.g., through Shutdown).
	Serve() error

	// Shutdown gracefully stops the server.
	Shutdown(ctx context.Context) error

	// IsAcceptingClients returns true if the service is ready to serve clients.
	IsAcceptingClients() bool

	// IsHealthy returns true if the service is considerably healthy.
	IsHealthy() bool
}

type ServiceInfo struct {
	Name        string
	Description string
}
