package pipeline

import (
	"github.com/hari134/comet/builder/container"
)
/*
	The pipeline interface defines function signatures for a build pipeline.
	The build pipeline will utilize a single container instance for all the stages.
	In the future the pipeline should be extended to allow multiple container across stages.
*/
type Pipeline interface{
	Run() error
	AddStage(stage Stage)
}


type SerialPipeline struct{
	container container.BuildContainer
	stages []Stage
}

func NewSerialPipeline(buildContainer container.BuildContainer) *SerialPipeline{
	return &SerialPipeline{
		buildContainer,
		[]Stage{},
	}
}

func (pipeline *SerialPipeline) AddStage(stage Stage){
	pipeline.stages = append(pipeline.stages,stage)
}

func (pipeline *SerialPipeline) Run() error{
	for _, stage := range pipeline.stages{
		_,err := stage.Execute(pipeline.container)
		if err != nil{
			return err
		}
	}
	pipeline.container.Stop()
	if err := pipeline.container.Remove(); err != nil{
		return err
	}
	return nil
}