package pipeline

import (
	"fmt"
)

/* Pipeline is a interface which defines functions:
 Run() - Runs the functions in pipeline stages
 AddDependency(constructor) - Takes a contructor function that returns the dependency
 AddStage(function) - Takes a function with the necessary dependencies as arguments in the function which is injected by any
 DI framework like dig. The function can return any error which is returned as is, this is useful to return at any point in the stage.

 Note : For simplicity this assumes that only a single instance exists for a dependency of a given type. If multiple dependencies exists
			  of same type exist then wrap then in structs of different types.
*/
type Pipeline interface {
    AddStage(stage Stage) Pipeline
    Run(depManager DependencyManager) error
}


type pipelineImpl struct {
    stages []Stage
}

// NewPipeline creates a new pipeline instance.
func NewPipeline() Pipeline {
    return &pipelineImpl{}
}

// AddStage adds a stage to the pipeline.
func (p *pipelineImpl) AddStage(st Stage) Pipeline {
    p.stages = append(p.stages, st)
    return p
}

// Run executes all stages in the pipeline using the dependency manager.
func (p *pipelineImpl) Run(depManager DependencyManager) error {
    for _, st := range p.stages {
        fmt.Printf("Running stage: %s\n", st.Name)
        err := st.Execute(depManager)
        if err != nil {
            fmt.Printf("Error in stage %s: %v\n", st.Name, err)
            return err
        }
    }
    return nil
}