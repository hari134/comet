package pipeline

import "errors"


type PipelineFactory interface {
	Register(name string, pipeline Pipeline)
	Get(name string) (Pipeline, error)
}

type DefaultPipelineFactory struct {
	registry map[string]Pipeline
}

func (pf *DefaultPipelineFactory) Register(name string ,pipeline Pipeline){
	pf.registry[name] = pipeline
}

func (pf *DefaultPipelineFactory) Get(name string ) (Pipeline, error){
	pipeline, ok := pf.registry[name]
	if !ok{
		return nil, errors.New("no such pipeline exists")
	}
	return pipeline, nil
}

func NewFactory() PipelineFactory{
	return &DefaultPipelineFactory{}
}
