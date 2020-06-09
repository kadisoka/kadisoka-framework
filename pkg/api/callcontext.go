package api

import (
	"context"

	"github.com/google/uuid"
)

type RequestID = uuid.UUID

type CallContext interface {
	context.Context

	// MethodName returns the name of the method this call is directed to.
	//
	// For HTTP, this method returns the value as "<HTTP_METHOD> <URL>", e.g.,
	// GET /users/me
	//
	MethodName() string

	// RequestID returns the idempotency token, if provided.
	//
	// https://www.youtube.com/watch?v=IP-rGJKSZ3s
	RequestID() *RequestID
}
