package local

import (
	"io"
	"os"
	"path/filepath"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/errors"
	mediastore "github.com/kadisoka/kadisoka-framework/foundation/pkg/media/store"
)

type Config struct {
	FolderPath string `env:"FOLDER_PATH"`
}

const ServiceName = "local"

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

func NewService(config mediastore.ServiceConfig) (mediastore.Service, error) {
	if config == nil {
		return nil, errors.ArgMsg("config", "missing")
	}

	conf, ok := config.(*Config)
	if !ok {
		return nil, errors.ArgMsg("config", "type invalid")
	}

	return &Service{
		folderPath: conf.FolderPath,
	}, nil
}

type Service struct {
	folderPath string
}

var _ mediastore.Service = &Service{}

// PutObject is required by mediastore.Service interface.
func (objStoreCl *Service) PutObject(
	targetKey string, contentSource io.Reader,
) (finalURL string, err error) {
	targetName := filepath.Join(objStoreCl.folderPath, targetKey)
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
