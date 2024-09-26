package storage


type Store interface{
	Get(key string) ([]byte,error)
	Put(fileData []byte, key string) error
}