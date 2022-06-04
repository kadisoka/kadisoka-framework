package local

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/alloyzeus/go-azfl/azfl/errors"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/app"
	imagesrv "github.com/kadisoka/kadisoka-framework/foundation/pkg/media/image/server"
	mediastore "github.com/kadisoka/kadisoka-framework/foundation/pkg/media/store"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/webui"
)

type Config struct {
	FolderPath string `env:"FOLDER_PATH"`

	ImagesServePath    string                  `env:"IMAGES_SERVE_PATH"`
	ImageServerHandler *imagesrv.HandlerConfig `env:"IMAGE_SERVER"`

	ServerServePath string `env:"SERVER_SERVE_PATH"`
	ServerServePort int16  `env:"SERVER_SERVE_PORT"`
}

const ServiceName = "local"

var serviceInfo = app.ServiceInfo{
	Name:        "Local Object Store Service",
	Description: "A generic service for serving objects",
}

func init() {
	mediastore.RegisterModule(
		ServiceName,
		mediastore.Module{
			NewService: NewService,
			ServiceConfigSkeleton: func() mediastore.ServiceConfig {
				cfg := ConfigSkeleton()
				return &cfg
			},
		})
}

func ConfigSkeleton() Config { return Config{} }

func NewService(
	config mediastore.ServiceConfig,
	appApp app.App,
) (mediastore.Service, error) {
	if config == nil {
		return nil, errors.ArgMsg("config", "missing")
	}

	conf, ok := config.(*Config)
	if !ok || conf == nil {
		return nil, errors.ArgMsg("config", "type invalid")
	}

	if conf.ServerServePath == "" {
		return nil, errors.ArgMsg("config.ServerServePath", "empty")
	}
	if conf.ServerServePort < 1024 {
		return nil, errors.ArgMsg("config.ServerServePort", "invalid")
	}

	filesBasePath, err := filepath.Abs(conf.FolderPath)
	if err != nil {
		return nil, errors.ArgWrap("config.FolderPath", "absolute resolution", err)
	}

	filesDirNoSlash := filesBasePath

	fileServer := webui.ETagHandler(
		http.StripPrefix(conf.ServerServePath,
			http.FileServer(
				http.Dir(filesDirNoSlash))))

	cfg := *conf
	if cfg.ImagesServePath == "" {
		cfg.ImagesServePath = conf.ServerServePath + "/_imgs"
	}

	var imageServerHandlerCfg imagesrv.HandlerConfig
	if cfg.ImageServerHandler != nil {
		imageServerHandlerCfg = *cfg.ImageServerHandler
	}
	imageServerHandlerCfg.RawFilesDir = filesDirNoSlash
	imageServer, err := imagesrv.NewHandler(imageServerHandlerCfg)
	if err != nil {
		return nil, errors.Wrap("image server handler initialization", err)
	}

	httpServeMux := http.NewServeMux()
	httpServeMux.Handle("/", fileServer)
	httpServeMux.Handle(cfg.ImagesServePath+"/",
		http.StripPrefix(cfg.ImagesServePath+"/", imageServer))

	svc := &Service{
		config:       cfg,
		httpServeMux: httpServeMux,
	}

	appApp.AddServiceServer(svc)

	return svc, nil
}

type Service struct {
	config Config

	shutdownOnce sync.Once
	shuttingDown bool

	httpServer   *http.Server
	httpServeMux *http.ServeMux
}

var _ mediastore.Service = &Service{}
var _ app.ServiceServer = &Service{}

// PutObject is required by mediastore.Service interface.
func (svc *Service) PutObject(
	targetKey string, contentSource io.Reader,
) (finalURL string, err error) {
	folderPath := svc.config.FolderPath
	targetKeyParts := strings.Split(targetKey, "/")
	if len(targetKeyParts) > 1 {
		targetPath := filepath.Join(targetKeyParts[:len(targetKeyParts)-1]...)
		targetPath = filepath.Join(folderPath, targetPath)
		os.MkdirAll(targetPath, 0700)
	}

	targetName := filepath.Join(folderPath, targetKey)
	targetFile, err := os.Create(targetName)
	if err != nil {
		return "", errors.Wrap("create file", err)
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, contentSource)
	if err != nil {
		return "", errors.Wrap("write content", err)
	}

	//TODO: final URL! ask the HTTP server to provide this.
	return targetKey, nil
}

// ServiceInfo conforms app.ServiceServer interface.
func (svc *Service) ServiceInfo() app.ServiceInfo { return serviceInfo }

// IsAcceptingClients is required by app.ServiceServer
func (svc *Service) IsAcceptingClients() bool {
	return !svc.shuttingDown && svc.IsHealthy()
}

// IsHealthy is required by app.ServiceServer
func (svc *Service) IsHealthy() bool { return true }

// Serve is required by app.ServiceServer
func (svc *Service) Serve() error {
	servePort := svc.config.ServerServePort

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", servePort),
		Handler: svc}
	svc.httpServer = httpServer
	err := svc.httpServer.ListenAndServe()
	if err == nil {
		if !svc.shuttingDown {
			return errors.Msg("server stopped unexpectedly")
		}
		return nil
	}
	if err == http.ErrServerClosed && svc.shuttingDown {
		return nil
	}
	return err
}

// Shutdown conforms app.ServiceServer interface.
func (svc *Service) Shutdown(ctx context.Context) error {
	svc.shutdownOnce.Do(func() {
		svc.shuttingDown = true
		svc.httpServer.Shutdown(ctx)
	})

	return nil
}

// ServeHTTP conforms Go's HTTP Handler interface.
func (svc *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	svc.httpServeMux.ServeHTTP(w, r)
}
