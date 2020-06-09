package data

// Error is an abstract error type for all data-related errors.
type Error interface {
	error
	DataError() Error
}

func Err(err error) error { return &wrappingError{err} }

type wrappingError struct {
	err error
}

var _ Error = &wrappingError{}

func (e wrappingError) Error() string    { return e.err.Error() }
func (e wrappingError) DataError() Error { return &e }

type msgError struct {
	msg string
}

var _ Error = &msgError{}

func (e msgError) Error() string    { return e.msg }
func (e msgError) DataError() Error { return &e }

func Malformed(err error) error {
	return &malformedError{err}
}

type malformedError struct {
	err error
}

var _ Error = &malformedError{}

func (e malformedError) Error() string {
	if e.err != nil {
		return "malformed: " + e.err.Error()
	}
	return "malformed"
}

func (e malformedError) DataError() Error { return &e }

var (
	ErrEmpty           = &msgError{"empty"}
	ErrMalformed       = Malformed(nil)
	ErrTypeUnsupported = &msgError{"type unsupported"}
)
