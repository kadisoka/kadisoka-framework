package webui

import (
	"path/filepath"
	"reflect"

	"github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/webui"

	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
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
