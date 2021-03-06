package errors

import (
	"net/http"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	accesserrs "github.com/alloyzeus/go-azfl/azfl/errors/access"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api/rest"
)

const HTTPStatusUnknown = 0

func Response(err error) (statusCode int, respData *rest.ErrorResponse) {
	return responseStatusCode(err), responseBody(err)
}

func responseStatusCode(err error) (httpStatusCode int) {
	if err == nil {
		return http.StatusOK
	}
	if x, ok := err.(interface{ RESTStatusCode() int }); ok && x != nil {
		code := x.RESTStatusCode()
		return code
	}

	if err == errors.ErrUnimplemented {
		return http.StatusNotImplemented
	}

	switch err.(type) {
	case accesserrs.Error:
		return http.StatusNotFound
	case errors.CallError:
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

func responseBody(err error) *rest.ErrorResponse {
	if err == nil {
		return nil
	}
	if d, ok := err.(interface{ RESTErrorResponseBody() *rest.ErrorResponse }); ok && d != nil {
		return d.RESTErrorResponseBody()
	}
	return &rest.ErrorResponse{}
}
