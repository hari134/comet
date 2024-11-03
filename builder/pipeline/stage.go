package pipeline


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
