package pipeline


type Stage struct {
    Name     string
    Execute  func(depManager DependencyManager) error
}

func NewStage(name string, executeFunc func(depManager DependencyManager) error) Stage {
    return Stage{
        Name:    name,
        Execute: executeFunc,
    }
}
