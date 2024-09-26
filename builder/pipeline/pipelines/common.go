package pipelines

import (
	"github.com/hari134/comet/builder/pipeline"
)


func copyTarToContainer(ctx *pipeline.PipelineContext) error {
	buildContainer, err := ctx.GetContainer()
	if err != nil {
		return err
	}
	tarFile, err := ctx.GetProjectTarFile()
	if err != nil {
		return err
	}
	return buildContainer.CopyToContainer(tarFile, "/app")
}

func copyDistFromContainer(ctx *pipeline.PipelineContext) error {
	buildContainer, err := ctx.GetContainer()
	if err != nil {
		return err
	}
	_, err = buildContainer.CopyFromContainer("/app/dist")
	return err
}