package rest

import (
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam/rest/logging"
)

var (
	log    = logging.NewPkgLogger()
	logReq = log.WithRequest
)
