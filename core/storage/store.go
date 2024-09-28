package storage

import "context"


type Store interface{
	Get(ctx context.Context,key string) ([]byte,error)
	Put(ctx context.Context,fileData []byte, key string) error
}