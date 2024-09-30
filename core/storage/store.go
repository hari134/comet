package storage

import "context"

type Store interface {
	Get(ctx context.Context, bucket string, key string) ([]byte, error)
	Put(ctx context.Context, fileData []byte, bucket string, key string) error
}
