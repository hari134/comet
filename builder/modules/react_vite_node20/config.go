package reactvitenode20

import (
	"bytes"

	"github.com/hari134/comet/builder/container"
	"github.com/hari134/comet/core/storage"
)


type Config struct{
	BuildContainer container.BuildContainer
	ProjectStorageConfig ProjectStorageConfig
	ProjectFileData ProjectFileData
	Store storage.Store
}


type ProjectStorageConfig struct{
	ProjectStorageKey string
	ProjectStorageBucket string
}


type ProjectFileData struct{
	DirName string
	ProjectTarFile *bytes.Buffer
}