package builder

import (
	"fmt"

	"github.com/hari134/comet/builder/container"
	"github.com/hari134/comet/builder/modules"
	reactvitenode20 "github.com/hari134/comet/builder/modules/react_vite_node20"
	"github.com/hari134/comet/builder/relay"
	"github.com/hari134/comet/core/storage"
)

type Builder struct {
	Store            storage.Store
	ContainerManager container.ContainerManager
	PipelineFactory  modules.PipelineFactory
}

func NewBuilder(store storage.Store,containerManager container.ContainerManager,pipelineFactory modules.PipelineFactory) *Builder{
	return &Builder{
		Store:store,
		ContainerManager: containerManager,
		PipelineFactory: pipelineFactory,
	}
}
func (bSvc *Builder) DeployProject(projectStorageKey, projectStorageBucket, buildEnvType string) error {
	buildPipeline, err := bSvc.PipelineFactory.Get(buildEnvType)
	if err != nil {
		return err
	}
	switch buildEnvType {
	case "reactvitenode20":
		buildContainer, err := bSvc.ContainerManager.NewBuildContainer(buildEnvType)
		if err != nil {
			return err
		}

		cfg := reactvitenode20.PipelineConfig{
			BuildContainer: buildContainer,
			ProjectStorageConfig: reactvitenode20.ProjectStorageConfig{
				ProjectStorageKey:    projectStorageKey,
				ProjectStorageBucket: projectStorageBucket,
			},
			Store: bSvc.Store,
			StreamConfig: reactvitenode20.StreamConfig{
				StreamingEnabled: true,
				Output: make(chan relay.StreamData,1),
			},
		}
		if err := buildPipeline.Run(&cfg); err != nil {
			return err
		}
	default:
		return fmt.Errorf("no such build environment type %v",buildEnvType)
	}
	return nil
}
