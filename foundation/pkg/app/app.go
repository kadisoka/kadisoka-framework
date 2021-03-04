//

package app

import (
	"sync"

	"github.com/alloyzeus/go-azcore/azcore/errors"
	"github.com/rez-go/stev"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/realm"
)

const EnvPrefixDefault = "APP_"

type App interface {
	RealmInfo() realm.Info

	AppInfo() Info
	InstanceID() string

	AddServer(ServiceServer)
	Run()
	IsAllServersAcceptingClients() bool
}

type Info struct {
	Name string
}

const (
	NameDefault = "Kadisoka-based App"
)

func DefaultInfo() Info {
	return Info{
		Name: NameDefault,
	}
}

func InfoFromEnvOrDefault() (Info, error) {
	info := DefaultInfo()
	err := stev.LoadEnv(EnvPrefixDefault, &info)
	if err != nil {
		return DefaultInfo(), errors.Wrap("info loading from environment variables", err)
	}
	return info, nil
}

type AppBase struct {
	realmInfo realm.Info

	appInfo    Info
	instanceID string

	servers   []ServiceServer
	serversMu sync.RWMutex
}

func (appBase *AppBase) RealmInfo() realm.Info { return appBase.realmInfo }

func (appBase *AppBase) AppInfo() Info { return appBase.appInfo }

func (appBase *AppBase) InstanceID() string { return appBase.instanceID }

// AddServer adds a server to be run simultaneously. Do NOT call this
// method after the app has been started.
func (appBase *AppBase) AddServer(srv ServiceServer) {
	appBase.serversMu.Lock()
	appBase.servers = append(appBase.servers, srv)
	appBase.serversMu.Unlock()
}

// Run runs all the servers. Do NOT add any new server after this method
// was called.
func (appBase *AppBase) Run() {
	RunServers(appBase.Servers(), nil)
}

// IsAllServersAcceptingClients checks if every server is accepting clients.
func (appBase *AppBase) IsAllServersAcceptingClients() bool {
	servers := appBase.Servers()
	for _, srv := range servers {
		if !srv.IsAcceptingClients() {
			return false
		}
	}
	return true
}

// Servers returns an array of servers added to this app.
func (appBase *AppBase) Servers() []ServiceServer {
	out := make([]ServiceServer, len(appBase.servers))
	appBase.serversMu.RLock()
	copy(out, appBase.servers)
	appBase.serversMu.RUnlock()
	return out
}

var (
	defApp     App
	defAppOnce sync.Once
)

func InitByEnvDefault() (App, error) {
	appInfo, err := InfoFromEnvOrDefault()
	if err != nil {
		return nil, errors.Wrap("app info loading", err)
	}
	realmInfo, err := realm.InfoFromEnvOrDefault()
	if err != nil {
		return nil, errors.Wrap("realm info loading", err)
	}
	return Init(&realmInfo, &appInfo)
}

func Init(realmInfo *realm.Info, appInfo *Info) (App, error) {
	var err error
	defAppOnce.Do(func() {
		if realmInfo == nil {
			i := realm.DefaultInfo()
			realmInfo = &i
		}

		if appInfo == nil {
			i := DefaultInfo()
			appInfo = &i
		}

		var unameStr string
		unameStr, err = unameString()
		if err != nil {
			return
		}

		var taskID string
		taskID, _, err = getECSTaskID()
		if err != nil {
			return
		}

		var instanceID string
		if taskID != "" {
			if unameStr != "" {
				instanceID = taskID + " (" + unameStr + ")"
			} else {
				instanceID = taskID
			}
		} else {
			instanceID = unameStr
		}

		defApp = &AppBase{
			realmInfo:  *realmInfo,
			appInfo:    *appInfo,
			instanceID: instanceID,
		}
	})

	if err != nil {
		return nil, err
	}

	return defApp, nil
}
