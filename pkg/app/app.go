//

package app

import (
	"sync"

	"github.com/rez-go/stev"

	"github.com/kadisoka/foundation/pkg/errors"
)

const EnvPrefixDefault = "APP_"

const (
	NameDefault                    = "Kadisoka"
	URLDefault                     = "https://github.com/kadisoka"
	EmailDefault                   = "nop@example.com"
	NotificationEmailSenderDefault = "no-reply@example.com"
	TeamNameDefault                = "Team Kadisoka"
)

func DefaultInfo() Info {
	return Info{
		Name:                    NameDefault,
		URL:                     URLDefault,
		Email:                   EmailDefault,
		NotificationEmailSender: NotificationEmailSenderDefault,
		TeamName:                TeamNameDefault,
	}
}

type Info struct {
	// Name of the app
	Name string
	// Canonical URL of the app
	URL                     string
	TermsOfServiceURL       string
	PrivacyPolicyURL        string
	Email                   string
	NotificationEmailSender string
	TeamName                string
}

type App interface {
	AppInfo() Info
	InstanceID() string

	AddServer(ServiceServer)
	Run()
}

type AppBase struct {
	appInfo    Info
	instanceID string

	servers   []ServiceServer
	serversMu sync.RWMutex
}

func (appBase AppBase) AppInfo() Info      { return appBase.appInfo }
func (appBase AppBase) InstanceID() string { return appBase.instanceID }

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
	RunServers(appBase.Servers())
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
	info := DefaultInfo()
	err := stev.LoadEnv(EnvPrefixDefault, &info)
	if err != nil {
		return nil, errors.Wrap("info loading from environment variables", err)
	}
	return Init(&info)
}

func Init(info *Info) (App, error) {
	var err error
	defAppOnce.Do(func() {
		var appInfo Info
		if info != nil {
			appInfo = *info
		} else {
			appInfo = DefaultInfo()
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

		defApp = &AppBase{appInfo: appInfo, instanceID: instanceID}
	})

	if err != nil {
		return nil, err
	}

	return defApp, nil
}
