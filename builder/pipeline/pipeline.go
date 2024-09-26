package pipeline

// Pipeline interface defines function signatures for a build pipeline.
// The build pipeline will take a container instance to execute for all stages.
type Pipeline interface {
	Run(ctx *PipelineContext) error
	AddStage(stage Stage) Pipeline
}

type SerialPipeline struct {
	stages []Stage
}

// NewSerialPipeline creates a new SerialPipeline instance.
func NewSerialPipeline() Pipeline {
	return &SerialPipeline{
		stages: []Stage{},
	}
}

// AddStage adds a stage to the serial pipeline and returns the pipeline for chaining.
func (pipeline *SerialPipeline) AddStage(stage Stage) Pipeline{
	pipeline.stages = append(pipeline.stages, stage)
	return pipeline
}

// Run executes all stages in sequence. If a stage fails, the execution stops.
func (pipeline *SerialPipeline) Run(ctx *PipelineContext) error {
	for _, stage := range pipeline.stages {
		if err := stage.Execute(ctx); err != nil {
			return err
		}
	}

	// Cleanup the container after all stages are executed
	container, err := ctx.GetContainer()
	if err != nil{
		return err
	}

	if err := container.Stop(); err != nil {
		return err
	}
	if err := container.Remove(); err != nil {
		return err
	}

	return nil
}
