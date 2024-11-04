package pipeline

import (
	"fmt"
)

type PipelineConfig interface{
    IsStreamingEnabled() bool
}

// TODO add support for streaming logs
type Pipeline interface {
    AddStage(stage Stage) Pipeline
    Run(config PipelineConfig) error
}


type pipeline struct {
    stages []Stage
}

// NewPipeline creates a new pipeline instance.
func NewPipeline() Pipeline {
    return &pipeline{}
}

// AddStage adds a stage to the pipeline.
func (p *pipeline) AddStage(st Stage) Pipeline {
    p.stages = append(p.stages, st)
    return p
}

// Run executes all stages in the pipeline using the dependency manager.
func (p *pipeline) Run(config PipelineConfig) error {
    for _, st := range p.stages {
        fmt.Printf("Running stage: %s\n", st.Name)
        err := st.Execute(config)
        if err != nil {
            fmt.Printf("Error in stage %s: %v\n", st.Name, err)
            return err
        }
    }
    return nil
}