package reactvitenode20

import (
	"bytes"

	"github.com/hari134/comet/builder/container"
	"github.com/hari134/comet/builder/relay"
	"github.com/hari134/comet/core/storage"
)


type PipelineConfig struct{
	BuildContainer container.BuildContainer
	ProjectStorageConfig ProjectStorageConfig
	ProjectFileData ProjectFileData
	Store storage.Store
	StreamConfig StreamConfig
}

type StreamConfig struct{
	StreamingEnabled bool
	Output chan relay.StreamData
}

func (config *PipelineConfig) IsStreamingEnabled() bool{
	return config.StreamConfig.StreamingEnabled
}


func (config *PipelineConfig) EnableStreaming() {
	config.StreamConfig.StreamingEnabled = true
}

type ProjectStorageConfig struct{
	ProjectStorageKey string
	ProjectStorageBucket string
}


type ProjectFileData struct{
	DirName string
	ProjectTarFile *bytes.Buffer
}