package rest

import (
	"encoding/json"
	"net/http"
)

func IsCallErrorStatusCode(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
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

func RespondErrorEmpty(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte("{}"))
}

func RespondError(w http.ResponseWriter, statusCode int, errorData ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(errorData)
	if err != nil {
		panic(err)
	}
}
