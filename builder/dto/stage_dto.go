package dto

import "bytes"

type ProjectStorageConfig struct{
	ProjectStorageKey string
	ProjectStorageBucket string
}


type ProjectFileData struct{
	DirName string
	ProjectTarFile *bytes.Buffer
}