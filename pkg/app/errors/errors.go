package errors

import (
	"github.com/kadisoka/foundation/pkg/errors"
)

type Error interface {
	//TODO: details: which part of the application. Configuration?

	error
	ApplicationError() Error
}

func New() Error {
	return &applicationBase{}
}

type applicationBase struct{}

var _ Error = &applicationBase{}

func (e *applicationBase) Error() string           { return "application error" }
func (e *applicationBase) ApplicationError() Error { return e }

type Configuration interface {
	//TODO: details: which configuration fields, etc.

	Error
	ConfigurationError() Configuration
}

func NewConfiguration(innerErr error) Configuration {
	return &configurationWrap{innerErr}
}

func NewConfigurationMsg(errMsg string) Configuration {
	return &configurationWrap{errors.New(errMsg)}
}

type configurationWrap struct {
	innerErr error
}

var _ Configuration = &configurationWrap{}

func (e *configurationWrap) Error() string {
	if e != nil && e.innerErr != nil {
		return e.innerErr.Error()
	}
	return "configuration error"
}

func (e *configurationWrap) Unwrap() error {
	if e != nil {
		return e.innerErr
	}
	return nil
}

func (e *configurationWrap) ApplicationError() Error           { return e }
func (e *configurationWrap) ConfigurationError() Configuration { return e }
