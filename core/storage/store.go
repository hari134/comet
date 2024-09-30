package storage

import (
	"bytes"
	"context"
)

type Store interface {
	Get(ctx context.Context, bucket string, key string) (*bytes.Buffer, error)
	Put(ctx context.Context, fileData *bytes.Buffer, bucket string, key string) error
}
