package pipeline

import "errors"


type PipelineFactory interface {
	Register(name string, factory func() Pipeline)
	Get(name string) (Pipeline, error)
}

type DefaultPipelineFactory struct {
	registry map[string]func() Pipeline
}

func (pf *DefaultPipelineFactory) Register(name string ,factory func() Pipeline){
	pf.registry[name] = factory
}

func (pf *DefaultPipelineFactory) Get(name string ) (Pipeline, error){
	factory, ok := pf.registry[name]
	if !ok{
		return nil, errors.New("no such pipeline exists")
	}
	return factory(), nil
}
