package webui

import (
	"path/filepath"
	"reflect"

	"github.com/kadisoka/foundation/pkg/webui"

	"github.com/kadisoka/iam/pkg/iam"
)

var ResourcesDirDefault string

func init() {
	type t int
	pkgPath := reflect.TypeOf(t(0)).PkgPath()
	ResourcesDirDefault = filepath.Join(pkgPath, "resources")
}

type ServerConfig struct {
	Server webui.ServerConfig `env:",squash"`
	URLs   iam.WebUIURLs      `env:"URLS,squash"`
}
