package iam

import (
	"github.com/alloyzeus/go-azfl/azfl/errors"
)

func ReqFieldErr(fieldName string, err error) error {
	return &reqFieldError{errors.Ent(fieldName, err)}
}

func ReqFieldErrMsg(fieldName, errMsg string) error {
	return &reqFieldError{errors.EntMsg(fieldName, errMsg)}
}

type reqFieldError struct {
	errors.EntityError
}

var (
	_ errors.CallError = &reqFieldError{}
)

func (e reqFieldError) CallError()        {}
func (e reqFieldError) FieldName() string { return e.EntityError.EntityIdentifier() }
