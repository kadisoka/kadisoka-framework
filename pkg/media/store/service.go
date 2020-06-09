package store

import (
	"io"
)

type Service interface {
	PutObject(objectKey string, content io.Reader) (publicURL string, err error)
}
