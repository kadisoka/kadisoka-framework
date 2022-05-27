package pnv10n

import (
	"sync"
)

type ModuleConfig interface {
	SMSDeliveryServiceConfig() interface{}
}

// Module provides details about a module of pn10n
type Module struct {
	// ConfigSkeleton returns a configuration skeleton for the module.
	ConfigSkeleton func() ModuleConfig

	// NewSMSDeliveryService returns an instance of module's SMS delivery
	// service. If the module does not provide such functionality, it
	// must return nil.
	NewSMSDeliveryService func(config interface{}) SMSDeliveryService
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

func NewSMSDeliveryService(
	serviceName string, config interface{},
) (SMSDeliveryService, error) {
	if serviceName == "" {
		return nil, nil
	}

	var module Module
	modulesMu.RLock()
	module = modules[serviceName]
	modulesMu.RUnlock()

	return module.NewSMSDeliveryService(config), nil
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
