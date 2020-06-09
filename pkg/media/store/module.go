package store

import (
	"sync"
)

type Module struct {
	ConfigSkeleton   func() interface{} // returns pointer
	NewServiceClient func(config interface{}) (Service, error)
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

func NewServiceClient(
	serviceName string,
	config interface{},
) (Service, error) {
	if serviceName == "" {
		return nil, nil
	}

	var module Module
	modulesMu.RLock()
	module, _ = modules[serviceName]
	modulesMu.RUnlock()

	return module.NewServiceClient(config)
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

func ModuleConfigSkeletons() map[string]interface{} {
	modulesMu.RLock()
	defer modulesMu.RUnlock()

	configs := map[string]interface{}{}
	for serviceName, mod := range modules {
		if mod.ConfigSkeleton != nil {
			configs[serviceName] = mod.ConfigSkeleton()
		}
	}

	return configs
}
