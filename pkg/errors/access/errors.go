package access

import (
	"github.com/citadelium/foundation/pkg/errors"
)

// Error is an abstract error type for all API access-related errors.
type Error interface {
	errors.CallError
	AccessError() Error
}

// New creates an API access-error which wraps another error.
func New(innerErr error) Error {
	return &errorWrap{innerErr}
}

// Msg creates an API access-error which wraps simple message error.
func Msg(errMsg string) Error {
	return &errorWrap{errors.New(errMsg)}
}

// Wrap creates an API access-error which provides additional context
// to another error.
func Wrap(contextMsg string, causeErr error) Error {
	return &errorWrap{errors.Wrap(contextMsg, causeErr)}
}

var (
	_ Error              = &errorWrap{}
	_ errors.Unwrappable = &errorWrap{}
	_ errors.CallError   = &errorWrap{}
)

type errorWrap struct {
	innerErr error
}

func (e *errorWrap) Error() string {
	if e != nil && e.innerErr != nil {
		return e.innerErr.Error()
	}
	return "access error"
}

func (e *errorWrap) Unwrap() error {
	if e != nil {
		return e.innerErr
	}
	return nil
}

func (e errorWrap) AccessError() Error { return &e }
func (e errorWrap) CallError()         {}
