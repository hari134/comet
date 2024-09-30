package pipeline

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/hari134/comet/builder/container"
	"github.com/hari134/comet/core/storage"
)

// PipelineContext holds shared data for stages
// The buildContainer is a default parameter and other optional parameters are store in data map
type PipelineContext struct {
	container      container.BuildContainer
	projectTarFile *bytes.Buffer
	store 					storage.Store
	data           map[string]interface{}
}

func NewPipelineContext() *PipelineContext {
	return &PipelineContext{
		container: nil,
		data:      make(map[string]interface{}),
	}
}

func (ctx *PipelineContext) WithContainer(buildContainer container.BuildContainer) *PipelineContext {
	ctx.container = buildContainer
	return ctx
}


func (ctx *PipelineContext) WithStore(store storage.Store) *PipelineContext {
	ctx.store= store
	return ctx
}

func (ctx *PipelineContext) GetStore() (storage.Store,error) {
	if ctx.store == nil {
		return nil, errors.New("store not set in pipeline context")
	}
	return ctx.store, nil
}

func (pipeline *PipelineContext) WithProjectTarFile(tarFile *bytes.Buffer) *PipelineContext {
	pipeline.projectTarFile = tarFile
	return pipeline
}

func (ctx *PipelineContext) GetContainer() (container.BuildContainer, error) {
	if ctx.container == nil {
		return nil, errors.New("container not set in pipeline context")
	}
	return ctx.container, nil
}

func (ctx *PipelineContext) GetProjectTarFile() (*bytes.Buffer, error) {
	if ctx.projectTarFile== nil {
		return nil, errors.New("project tar file not set in pipeline context")
	}
	return ctx.projectTarFile, nil
}

func (ctx *PipelineContext) Set(key string, value interface{}) {
	ctx.data[key] = value
}

func (ctx *PipelineContext) Get(key string) (interface{}, error) {
	val, ok := ctx.data[key]
	if !ok {
		return nil, fmt.Errorf("key %s not found in context", key)
	}
	return val, nil
}
