package util

import (
	"context"

	"github.com/hari134/comet/builder/pipeline"
	"github.com/hari134/comet/builder/util"
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

func pullProjectFromStore(ctx *pipeline.PipelineContext) error{
	store ,err := ctx.GetStore()
	if err != nil{
		return err
	}
	projectStorageKeyRaw ,err := ctx.Get("projectStorageKey")
	if err != nil{
		return err
	}
	projectStorageKey,err := util.TypeAssert[string](projectStorageKeyRaw,"string")
	if err != nil{
		return err
	}
	projectStorageBucketRaw ,err := ctx.Get("projectStorageBucket")
	if err != nil{
		return err
	}
	projectStorageBucket,err := util.TypeAssert[string](projectStorageBucketRaw,"string")
	if err != nil{
		return err
	}
	projectTarFile , err := store.Get(context.Background(),projectStorageBucket,projectStorageKey)
	if err != nil{
		return err
	}
	ctx.SetProjectTarFile(projectTarFile)
	return nil
}

func copyDistFromContainer(ctx *pipeline.PipelineContext) error {
	buildContainer, err := ctx.GetContainer()
	if err != nil {
		return err
	}
	_, err = buildContainer.CopyFromContainer("/app/dist")
	return err
}