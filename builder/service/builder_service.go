package service

import (
	"fmt"

	"github.com/hari134/comet/builder/container"
	"github.com/hari134/comet/builder/modules"
	reactvitenode20 "github.com/hari134/comet/builder/modules/react_vite_node20"
	"github.com/hari134/comet/core/storage"
)

type BuilderService struct {
	Store            storage.Store
	ContainerManager container.ContainerManager
	PipelineFactory  modules.PipelineFactory
}

func NewBuilderService(store storage.Store,containerManager container.ContainerManager,pipelineFactory modules.PipelineFactory) *BuilderService{
	return &BuilderService{
		Store:store,
		ContainerManager: containerManager,
		PipelineFactory: pipelineFactory,
	}
}
func (bSvc *BuilderService) DeployProject(projectStorageKey, projectStorageBucket, buildEnvType string) error {
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
		cfg := reactvitenode20.Config{
			BuildContainer: buildContainer,
			ProjectStorageConfig: reactvitenode20.ProjectStorageConfig{
				ProjectStorageKey:    projectStorageKey,
				ProjectStorageBucket: projectStorageBucket,
			},
			Store: bSvc.Store,
		}
		if err := buildPipeline.Run(cfg); err != nil {
			return err
		}
	default:
		return fmt.Errorf("no such buildEnvType %v",buildEnvType)
	}
	return nil
}
