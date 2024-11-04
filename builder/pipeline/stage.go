package pipeline

// TODO add timeout or deadline for each stage and break with error

type Stage struct {
    Name     string
    Execute  func(PipelineConfig) error
}

func NewStage(name string, executeFunc func(config PipelineConfig) error) Stage {
    return Stage{
        Name:    name,
        Execute: executeFunc,
    }
}
