package modules

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/hari134/comet/builder/modules/react_vite_node20"
	"github.com/hari134/comet/builder/pipeline"
)


type PipelineFactory interface {
	Get(name string) (pipeline.Pipeline, error)
}

type DefaultPipelineFactory struct {
	registry map[string]pipeline.Pipeline
}

func (pf *DefaultPipelineFactory) Get(name string ) (pipeline.Pipeline, error){
	slog.Debug(fmt.Sprintf("got key : %s in factory get",name))
	pipeline, ok := pf.registry[name]
	if !ok{
		return nil, errors.New("no such pipeline exists")
	}
	return pipeline, nil
}

func NewFactory() PipelineFactory{
	factory := &DefaultPipelineFactory{
		registry : make(map[string]pipeline.Pipeline),
	 }
	factory.registry["reactvitenode20"] = reactvitenode20.ReactViteNode20Pipeline
	return factory
}
