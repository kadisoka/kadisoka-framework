package rest

import (
	"encoding/json"
	"net/http"
)

func IsCallErrorStatusCode(httpStatusCode int) bool {
	return httpStatusCode >= 400 && httpStatusCode < 500
}

type ErrorResponse struct {
	Code string `json:"code,omitempty"`

	// We use the term description because it describes the error
	// to the developer rather than a message for the end user.
	Description string `json:"description,omitempty"`

	Fields []ErrorResponseField `json:"fields,omitempty"`
	DocURL string               `json:"doc_url,omitempty"`
}

type ErrorResponseField struct {
	Field       string `json:"field"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	DocURL      string `json:"doc_url,omitempty"`
}

type EmptyRequest struct{}

type EmptyResponse struct{}

func RespondTo(w http.ResponseWriter) Responder { return Responder{w} }

type Responder struct {
	w http.ResponseWriter
}

// EmptyError responses with empty ErrorResponse with status code set to httpStatusCode.
func (r Responder) EmptyError(httpStatusCode int) {
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(httpStatusCode)
	r.w.Write([]byte("{}"))
}

// Error responses with payload provided as errorResp
func (r Responder) Error(errorData ErrorResponse, httpStatusCode int) {
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(httpStatusCode)
	err := json.NewEncoder(r.w).Encode(errorData)
	if err != nil {
		panic(err)
	}
}

func (r Responder) Success(successData interface{}) {
	//TODO: ensure that it's trully nil. It's possible that the data is not
	// actually nil.
	if successData == nil {
		r.w.WriteHeader(http.StatusNoContent)
		return
	}
	r.SuccessWithHTTPStatusCode(successData, http.StatusOK)
}

func (r Responder) SuccessWithHTTPStatusCode(successData interface{}, httpStatusCode int) {
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(httpStatusCode)
	err := json.NewEncoder(r.w).Encode(successData)
	if err != nil {
		panic(err)
	}
}

func RespondErrorEmpty(w http.ResponseWriter, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	w.Write([]byte("{}"))
}

func RespondError(w http.ResponseWriter, httpStatusCode int, errorData ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	err := json.NewEncoder(w).Encode(errorData)
	if err != nil {
		panic(err)
	}
}
