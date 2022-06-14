package eav10n

import (
	"sync"
)

type ModuleConfig interface {
	EmailDeliveryServiceConfig() interface{}
}

// Module provides details about a module of pn10n
type Module struct {
	// ConfigSkeleton returns a configuration skeleton for the module.
	ConfigSkeleton func() ModuleConfig

	// NewEmailDeliveryService returns an instance of module's email delivery
	// service. If the module does not provide such functionality, it
	// must return nil.
	NewEmailDeliveryService func(config interface{}) EmailDeliveryService
}

var (
	modules   = map[string]Module{}
	modulesMu sync.RWMutex
)

func ModuleNames() []string {
	modulesMu.Lock()
	defer modulesMu.Unlock()

	var names []string
	for name := range modules {
		names = append(names, name)
	}

	return names
}

func RegisterModule(
	serviceName string,
	module Module,
) {
	modulesMu.Lock()
	defer modulesMu.Unlock()

	if _, dup := modules[serviceName]; dup {
		panic("called twice for service " + serviceName)
	}

	modules[serviceName] = module
}

func NewEmailDeliveryService(
	serviceName string, config interface{},
) (EmailDeliveryService, error) {
	if serviceName == "" {
		return nil, nil
	}

	var module Module
	modulesMu.RLock()
	module = modules[serviceName]
	modulesMu.RUnlock()

	return module.NewEmailDeliveryService(config), nil
}

func ModuleConfigSkeletons() map[string]ModuleConfig {
	modulesMu.RLock()
	defer modulesMu.RUnlock()

	configs := map[string]ModuleConfig{}
	for serviceName, module := range modules {
		if module.ConfigSkeleton != nil {
			configs[serviceName] = module.ConfigSkeleton()
		}
	}

	return configs
}
