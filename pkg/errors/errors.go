// Package errors provides extra functionalities to that of Go's stdlib
// errors package.
//
//TODO: fields.
// e.g., errors.Msg("error message", errors.Str("name", name), errors.Err(err))
//    or errors.With().Str("name", name).Err(err).Msg("error message")
//    or errors.With().StrErr("name", nameErr).Msg("error message")
// (sounds like structured logging? exactly!)
package errors

import (
	"errors"
)

// Wraps Go's errors
var (
	As     = errors.As
	Is     = errors.Is
	New    = errors.New // Prefer Msg instead of New as it has better semantic
	Msg    = errors.New
	Unwrap = errors.Unwrap
)

type Unwrappable interface {
	error
	Unwrap() error
}

// Wrap creates a new error by providing context message to another error.
// It's recommended for the message to describe what the program did which
// caused the error.
//
//     err := fetchData(...)
//     if err != nil { return errors.Wrap("fetching data", err) }
//
func Wrap(contextMessage string, causeErr error) error {
	return &errorWrap{contextMessage, causeErr}
}

var _ Unwrappable = &errorWrap{}

type errorWrap struct {
	msg string
	err error
}

func (e errorWrap) Error() string {
	if e.msg != "" {
		if e.err != nil {
			return e.msg + ": " + e.err.Error()
		}
	}
	return e.msg
}

func (e errorWrap) Unwrap() error {
	return e.err
}

// ErrUnimplemented is used to declare that a functionality, or part of it,
// has not been implemented. This could be well mapped to some protocols'
// status code, e.g., HTTP's 501 and gRPC's 12 .
var ErrUnimplemented = Msg("unimplemented")
