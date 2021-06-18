//

package app

import (
	"sync"

	"github.com/alloyzeus/go-azfl/azfl/errors"
)

// App abstracts the application itself. There should be only one instance
// for a running instance of an app.
type App interface {
	AppInfo() Info
	InstanceID() string

	AddServer(ServiceServer)
	Run()
	IsAllServersAcceptingClients() bool
}

type Info struct {
	Name string

	BuildInfo BuildInfo
}

// HeaderString returns a line which contains build information about
// the app that can be used to provide such info in logs.
func (info Info) HeaderString() string {
	return info.Name +
		" revision " + info.BuildInfo.RevisionID +
		" built at " + info.BuildInfo.Timestamp
}

// AppBase is the base layer for an app.
type AppBase struct {
	appInfo    Info
	instanceID string

	servers   []ServiceServer
	serversMu sync.RWMutex
}

var _ App = &AppBase{}

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

func Instance() App {
	if defApp == nil {
		panic("App has not been initialized. Call app.Init to initialize App.")
	}
	return defApp
}

// Instantiate global instance of App with the default implementation.
func Init(appInfo Info) (App, error) {
	err := errors.Msg("already initialized")
	defAppOnce.Do(func() {
		err = nil

		if appInfo.BuildInfo.RevisionID == "" {
			err = errors.ArgMsg("appInfo.BuildInfo.RevisionID", "empty")
			return
		}
		if appInfo.BuildInfo.Timestamp == "" {
			err = errors.ArgMsg("appInfo.BuildInfo.RevisionID", "empty")
			return
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
			appInfo:    appInfo,
			instanceID: instanceID,
		}
	})

	if err != nil {
		return nil, err
	}

	return defApp, nil
}

// InitCustom
func InitCustom(customApp App) error {
	err := errors.Msg("already initialized")
	defAppOnce.Do(func() {
		err = nil
		defApp = customApp
	})
	return err
}
