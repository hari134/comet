package modules

import (
	"errors"

	"github.com/hari134/comet/builder/pipeline"
	"github.com/hari134/comet/builder/modules/react_vite_node20"
)


type PipelineFactory interface {
	Get(name string) (pipeline.Pipeline, error)
}

type DefaultPipelineFactory struct {
	registry map[string]pipeline.Pipeline
}

func (pf *DefaultPipelineFactory) Register(name string ,pipeline pipeline.Pipeline){
	pf.registry[name] = pipeline
}

func (pf *DefaultPipelineFactory) Get(name string ) (pipeline.Pipeline, error){
	pipeline, ok := pf.registry[name]
	if !ok{
		return nil, errors.New("no such pipeline exists")
	}
	return pipeline, nil
}

func NewFactory() PipelineFactory{
	factory := &DefaultPipelineFactory{}
	factory.registry["reactvitenode20"] = reactvitenode20.ReactViteNode20Pipeline
	return &DefaultPipelineFactory{}
}
