package store

import (
	"context"
	"io"
)

type ServiceConfig interface{}

type Service interface {
	PutObject(ctx context.Context, objectKey string, content io.Reader) (publicURL string, err error)
}
